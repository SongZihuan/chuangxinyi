package selfpay

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	"github.com/google/uuid"
	"github.com/wuntsong-org/wterrors"
	"time"
)

var BadName = errors.NewClass("bad name")

func NewWithdraw(ctx context.Context, cny int64, user *db.User, withdrawWay string, name string) (string, errors.WTError) {
	if withdrawWay == alipay.WithdrawAlipay || withdrawWay == wechatpay.WithdrawWechatpay {
		return "", BadPayWay.New()
	}

	if cny < 1000 { // 小于10
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
		_, err := withdrawModel.FindByWithdrawID(ctx, SelfpayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := SelfpayID(OutTradeNoUUID.String())

	withdraw := &db.Withdraw{
		UserId:      user.Id,
		WalletId:    user.WalletId,
		WithdrawId:  OutTradeNo,
		WithdrawWay: withdrawWay,
		Name:        name,
		Cny:         cny,
		Status:      db.WithdrawWait,
		WithdrawAt:  time.Now(),
	}

	_, err = balance.WithdrawWithInsert(ctx, user, withdraw) // 自带insert
	if errors.Is(err, balance.Insufficient) {
		return "", Insufficient.New()
	} else if err != nil {
		return "", errors.WarpQuick(err)
	}

	logger.Logger.WXInfo("收到用户人工提现申请，单号：%s", withdraw.WithdrawId)

	audit.NewUserAudit(user.Id, "人工提现（%s）发起成功（%.2f）", withdrawWay, float64(cny)/100.0)
	return OutTradeNo, nil
}

func NewAdminWithdraw(ctx context.Context, cny int64, user *db.User, withdrawWay string, name string) (string, errors.WTError) {
	if withdrawWay == alipay.WithdrawAlipay || withdrawWay == wechatpay.WithdrawWechatpay {
		return "", BadPayWay.New()
	}

	if cny < 0 { // 小于10
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
		_, err := withdrawModel.FindByWithdrawID(ctx, SelfpayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := SelfpayID(OutTradeNoUUID.String())

	withdraw := &db.Withdraw{
		UserId:      user.Id,
		WalletId:    user.WalletId,
		WithdrawId:  OutTradeNo,
		WithdrawWay: withdrawWay,
		Name:        name,
		Cny:         cny,
		Status:      db.WithdrawWait,
		WithdrawAt:  time.Now(),
	}

	_, err = balance.WithdrawWithInsert(ctx, user, withdraw) // 自带insert
	if errors.Is(err, balance.Insufficient) {
		return "", Insufficient.New()
	} else if err != nil {
		return "", errors.WarpQuick(err)
	}

	withdraw.Status = db.WithdrawOK
	withdraw.PayAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	err = withdrawModel.Update(ctx, withdraw)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	logger.Logger.WXInfo("收到管理员录入人工提现，单号：%s", withdraw.WithdrawId)

	sender.PhoneSendChange(withdraw.UserId, "提现额度（人工提现）")
	sender.EmailSendChange(withdraw.UserId, "提现额度（人工提现）")
	sender.MessageSendWithdraw(withdraw.UserId, withdraw.Cny, "人工")
	audit.NewUserAudit(withdraw.UserId, "人工提现已到账（%.2f）", float64(withdraw.Cny)/100.0)

	return OutTradeNo, nil
}
