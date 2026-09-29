package selfpay

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

func Withdraw(ctx context.Context, withdraw *db.Withdraw) errors.WTError { // 不限制支付方式
	if withdraw.WithdrawWay == alipay.WithdrawAlipay || withdraw.WithdrawWay == wechatpay.WithdrawWechatpay {
		return errors.Errorf("bad withdraw way")
	}

	if withdraw.Status != db.WithdrawWait {
		return errors.Errorf("double withdraw")
	}

	keyPay := fmt.Sprintf("withdraw:%s", withdraw.WithdrawId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return errors.Errorf("can not lock")
	}
	defer redis.ReleaseLock(keyPay)

	withdraw.Status = db.WithdrawOK
	withdraw.PayAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	err := withdrawModel.Update(ctx, withdraw)
	if err != nil {
		return errors.WarpQuick(err)
	}

	logger.Logger.WXInfo("人工提现成功：%.2f", float64(withdraw.Cny)/100.0)
	sender.PhoneSendChange(withdraw.UserId, "提现额度（人工提现）")
	sender.EmailSendChange(withdraw.UserId, "提现额度（人工提现）")
	sender.MessageSendWithdraw(withdraw.UserId, withdraw.Cny, "人工")
	audit.NewUserAudit(withdraw.UserId, "人工提现已到账（%.2f）", float64(withdraw.Cny)/100.0)

	return nil
}

func WithdrawFail(ctx context.Context, user *db.User, withdraw *db.Withdraw) errors.WTError { // 不限制支付方式
	if withdraw.WithdrawWay == alipay.WithdrawAlipay || withdraw.WithdrawWay == wechatpay.WithdrawWechatpay {
		return errors.Errorf("bad withdraw way")
	}

	if withdraw.Status != db.WithdrawWait {
		return errors.Errorf("double withdraw")
	}

	keyPay := fmt.Sprintf("withdraw:%s", withdraw.WithdrawId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return errors.Errorf("can not lock")
	}
	defer redis.ReleaseLock(keyPay)

	_, err := balance.WithdrawFail(ctx, user, withdraw)
	if err != nil {
		return errors.WarpQuick(err)
	}

	audit.NewUserAudit(withdraw.UserId, "人工提现失败（%.2f）", float64(withdraw.Cny)/100.00)

	return nil
}
