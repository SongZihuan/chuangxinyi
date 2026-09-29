package alipay

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
	WithdrawAlipay = "支付宝提现到余额"
)

var BadName = errors.NewClass("bad name")

func NewWithdraw(ctx context.Context, user *db.User, cny int64, identity string, name string) (string, errors.WTError) {
	if !config.BackendConfig.Alipay.UseWithdraw {
		return "", errors.Errorf("not ok")
	}

	if cny < 100 || cny > 20000 {
		return "", BadCNY.New()
	}

	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	companyModel := db.NewCompanyModel(mysql.MySQLConn)

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

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "withdraw", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
		_, err := withdrawModel.FindByWithdrawID(ctx, AlipayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := AlipayID(OutTradeNoUUID.String())

	bm := make(gopay.BodyMap)
	bm.Set("out_biz_no", OutTradeNo).
		Set("trans_amount", fmt.Sprintf("%.2f", float64(cny)/100.0)).
		Set("biz_scene", "DIRECT_TRANSFER").
		Set("product_code", "TRANS_ACCOUNT_NO_PWD").
		Set("order_title", fmt.Sprintf("%s提现", config.BackendConfig.User.ReadableName)).
		Set("remark", fmt.Sprintf("用户（%s）提现%.2f，提现ID：%s。", user.Uid, float64(cny)/100.0, OutTradeNo)).
		Set("business_params", `{"payer_show_name_use_alias":"true"}`)

	payee := make(gopay.BodyMap)
	payee.Set("identity", identity).
		Set("identity_type", "ALIPAY_LOGON_ID").
		Set("name", name)

	bm.Set("payee_info", payee)

	withdraw := &db.Withdraw{
		UserId:      user.Id,
		WalletId:    user.WalletId,
		WithdrawId:  OutTradeNo,
		WithdrawWay: WithdrawAlipay,
		Name:        name,
		AlipayLoginId: sql.NullString{
			Valid:  true,
			String: identity,
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

	resp, err := AlipayClient.FundTransUniTransfer(ctx, bm)
	if err != nil {
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", errors.WarpQuick(err)
	}

	if resp.Response.SubCode == "EXCEED_LIMIT_SM_AMOUNT" {
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", BadCNY.New()
	} else if resp.Response.SubCode == "PAYER_USER_INFO_ERROR" {
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", BadName.New()
	} else if resp.Response.SubCode == "EXCEED_LIMIT_MM_AMOUNT" {
		logger.Logger.WXInfo("支付宝月提现已达到上线")
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", errors.Errorf("system error")
	} else if resp.Response.SubCode == "EXCEED_LIMIT_DM_AMOUNT" {
		logger.Logger.WXInfo("支付宝日提现已达到上线")
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", errors.Errorf("system error")
	} else if resp.Response.SubCode == "PAYER_BALANCE_NOT_ENOUGH" {
		logger.Logger.Error("支付宝提现余额不足")
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", errors.Errorf("system error")
	} else if resp.Response.Status != "SUCCESS" {
		logger.Logger.Error("支付宝快捷提现失败：%s (%s)", resp.Response.Code, resp.Response.SubCode)
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
		return "", errors.Errorf("system error")
	}

	withdraw.OrderId = sql.NullString{
		Valid:  true,
		String: resp.Response.OrderId,
	}
	withdraw.PayFundOrderId = sql.NullString{
		Valid:  true,
		String: resp.Response.PayFundOrderId,
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	err = withdrawModel.Update(ctx, withdraw)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
	}

	audit.NewUserAudit(user.Id, "支付宝提现（%s）发起成功（%.2f）", WithdrawAlipay, float64(cny)/100.0)
	return OutTradeNo, nil
}
