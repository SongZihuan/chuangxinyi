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
	"gitee.com/wuntsong-auth/backend/src/sender"
	"github.com/SuperH-0630/gopay"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

func QueryWithdraw(ctx context.Context, user *db.User, withdraw *db.Withdraw) (int64, errors.WTError) {
	if !config.BackendConfig.Alipay.UseWithdraw {
		return 0, errors.Errorf("not ok")
	}

	if withdraw.WithdrawWay != WithdrawAlipay {
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
	bm.Set("product_code", "TRANS_ACCOUNT_NO_PWD")
	bm.Set("biz_scene", "DIRECT_TRANSFER")
	bm.Set("out_biz_no", withdraw.WithdrawId)

	res, err := AlipayClient.FundTransCommonQuery(ctx, bm)
	if err != nil {
		logger.Logger.Error("alipay withdraw error: %s", err.Error())
		return db.WithdrawWait, nil
	}

	if res.Response.Code != "10000" {
		logger.Logger.Error("alipay withdraw error: %s", res.Response.Msg)
		return db.WithdrawWait, nil
	}

	if res.Response.Status == "SUCCESS" {
		logger.Logger.WXInfo("支付宝提现成功（通过检查）：%.2f", float64(withdraw.Cny)/100.0)
		withdraw.Status = db.WithdrawOK

		payDate, err := time.Parse("2006-01-02 15:04:05", res.Response.PayDate)
		if err == nil {
			withdraw.PayAt = sql.NullTime{
				Valid: true,
				Time:  payDate,
			}
		}

		sender.PhoneSendChange(withdraw.UserId, "提现额度（支付宝快捷提现）")
		sender.EmailSendChange(withdraw.UserId, "提现额度（支付宝快捷提现）")
		sender.MessageSendWithdraw(withdraw.UserId, withdraw.Cny, "支付宝")
		audit.NewUserAudit(withdraw.UserId, "支付宝提现已到账（%.2f）", float64(withdraw.Cny)/100.0)
	} else if res.Response.Status == "FAIL" || res.Response.Status == "REFUND" {
		_, _ = balance.WithdrawFail(ctx, user, withdraw)
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	err = withdrawModel.Update(ctx, withdraw)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return 0, errors.WarpQuick(err)
	}

	return withdraw.Status, nil
}
