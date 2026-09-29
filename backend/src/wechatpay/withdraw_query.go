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
	"gitee.com/wuntsong-auth/backend/src/sender"
	"github.com/SuperH-0630/gopay"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

func QueryWithdraw(ctx context.Context, user *db.User, withdraw *db.Withdraw) (int64, errors.WTError) {
	if !config.BackendConfig.WeChatPay.UseWithdraw {
		return 0, errors.Errorf("not ok")
	}

	if withdraw.WithdrawWay != WithdrawWechatpay {
		return 0, errors.Errorf("bad withdraw way")
	}

	if withdraw.Status != db.WithdrawWait {
		return withdraw.Status, nil
	}

	keyPay := fmt.Sprintf("withdraw:%s", withdraw.WithdrawId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return db.WithdrawWait, nil
	}
	defer redis.ReleaseLock(keyPay)

	bm := make(gopay.BodyMap)
	bm.Set("need_query_detail", false)

	res, err := WeChatPayClient.V3TransferMerchantQuery(ctx, withdraw.WithdrawId, bm)
	if err != nil {
		logger.Logger.Error("wechat withdraw query error: %s", err.Error())
		return db.WithdrawWait, nil
	}

	if res.Code != 0 {
		logger.Logger.Error("wechat withdraw query error: %s", res.Error)
		return db.WithdrawWait, nil
	}

	switch res.Response.TransferBatch.BatchStatus {
	default:
		withdraw.Status = db.WithdrawWait
	case "FINISHED":
		logger.Logger.WXInfo("微信支付提现成功（通过检查）：%.2f", float64(withdraw.Cny)/100.0)
		withdraw.Status = db.WithdrawOK
		payDate, err := time.Parse(time.RFC3339, res.Response.TransferBatch.UpdateTime)
		if err == nil {
			withdraw.PayAt = sql.NullTime{
				Valid: true,
				Time:  payDate,
			}
		}

		sender.PhoneSendChange(withdraw.UserId, "提现额度（微信支付快捷提现）")
		sender.EmailSendChange(withdraw.UserId, "提现额度（微信支付快捷提现）")
		sender.MessageSendWithdraw(withdraw.UserId, withdraw.Cny, "微信支付")
		audit.NewUserAudit(withdraw.UserId, "微信支付提现已到账（%.2f）", float64(withdraw.Cny)/100.0)
	case "CLOSED":
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	err = withdrawModel.Update(ctx, withdraw)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	return withdraw.Status, nil
}
