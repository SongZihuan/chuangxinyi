package wechatpay

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"github.com/SuperH-0630/gopay"
	"github.com/google/uuid"
	"github.com/wuntsong-org/wterrors"
	"time"
)

const (
	WithdrawWechatpay = "微信支付提现到零钱"
)

var BadName = errors.NewClass("bad name")
var WithoutBindFuwuhao = errors.NewClass("without bind fuwuhao")

func NewWithdraw(ctx context.Context, user *db.User, cny int64, name string) (string, errors.WTError) {
	if !config.BackendConfig.WeChatPay.UseWithdraw {
		return "", errors.Errorf("not ok")
	}

	if cny < 100 || cny > 20000 {
		return "", BadCNY.New()
	}

	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	companyModel := db.NewCompanyModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)

	idcard, err := idcardModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return "", BadName.New()
	} else if err != nil {
		return "", errors.WarpQuick(err)
	}

	company, err := companyModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		company = &db.Company{}
	} else if err != nil {
		return "", errors.WarpQuick(err)
	}

	if idcard.UserName != name && company.CompanyName != name && company.LegalPersonName != name {
		return "", BadName.New()
	}

	wechat, err := wechatModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return "", WithoutBindFuwuhao.New()
	} else if err != nil {
		return "", errors.WarpQuick(err)
	} else if !wechat.OpenId.Valid {
		return "", WithoutBindFuwuhao.New()
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "withdraw", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
		_, err := withdrawModel.FindByWithdrawID(ctx, WeChatPayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := WeChatPayID(OutTradeNoUUID.String())

	remark := fmt.Sprintf("用户（%s）提现%.2f，提现ID：%s。", user.Uid, float64(cny)/100.0, OutTradeNo)

	bm := make(gopay.BodyMap)
	bm.Set("appid", config.BackendConfig.WeChatPay.AppID).
		Set("out_batch_no", OutTradeNoUUID).
		Set("batch_name", fmt.Sprintf("%s提现", config.BackendConfig.User.ReadableName)).
		Set("batch_remark", remark).
		Set("total_amount", cny).
		Set("total_num", 1)

	dt := make(gopay.BodyMap)
	dt.Set("out_detail_no", OutTradeNoUUID).
		Set("transfer_amount", cny).
		Set("transfer_remark", remark).
		Set("openid", wechat.Fuwuhao.String).
		Set("user_name", name)

	bm.Set("transfer_detail_list", []gopay.BodyMap{dt})

	withdraw := &db.Withdraw{
		UserId:      user.Id,
		WalletId:    user.WalletId,
		WithdrawId:  OutTradeNo,
		WithdrawWay: WithdrawWechatpay,
		Name:        name,
		WechatpayOpenId: sql.NullString{
			Valid:  true,
			String: wechat.Fuwuhao.String,
		},
		WechatpayUnionId: sql.NullString{
			Valid:  true,
			String: wechat.UnionId.String,
		},
		WechatpayNickname: sql.NullString{
			Valid:  true,
			String: wechat.Nickname.String,
		},
		Cny:        cny,
		Status:     db.WithdrawWait,
		WithdrawAt: time.Now(),
	}

	_, err = balance.WithdrawWithInsert(ctx, user, withdraw) // 自带insert
	if errors.Is(err, balance.Insufficient) {
		return "", Insufficient.New()
	} else if err != nil {
		return "", errors.WarpQuick(err)
	}

	res, err := WeChatPayClient.V3Transfer(ctx, bm)
	if err != nil {
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", errors.WarpQuick(err)
	} else if res.Code != 0 {
		if res.Code == 403 {
			logger.Logger.Error("微信支付提现余额不足，请注意。提现金额：%.2f", float64(withdraw.Cny)/100.0)
		} else {
			logger.Logger.Error("微信支付退款失败：%d", res.Code)
		}
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", errors.Errorf("system error")
	}

	withdraw.OrderId = sql.NullString{
		Valid:  true,
		String: res.Response.BatchId,
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	err = withdrawModel.Update(ctx, withdraw)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
	}

	audit.NewUserAudit(user.Id, "微信支付快捷提现（%s）发起成功（%.2f）", WithdrawWechatpay, float64(cny)/100.0)
	return OutTradeNo, nil
}
