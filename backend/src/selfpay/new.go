package selfpay

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/config"
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

var BadCNY = errors.NewClass("pay cny must in 0 - 10000000")
var BadPayWay = errors.NewClass("pay way error")
var Insufficient = errors.NewClass("insufficient") // 余额不足

func NewSelfPay(ctx context.Context, cny int64, user *db.User, payWay string, couponsID int64) (string, errors.WTError) {
	if payWay == alipay.PayWayPC || payWay == alipay.PayWayWap || payWay == wechatpay.PayWayNative || payWay == wechatpay.PayWayH5 || payWay == wechatpay.PayWayJSAPI {
		return "", BadPayWay.New()
	}

	if cny < 0 { // 小于0，大于等于十万元
		return "", BadCNY.New()
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		payModel := db.NewPayModel(mysql.MySQLConn)
		_, err := payModel.FindByPayID(ctx, SelfpayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := SelfpayID(OutTradeNoUUID.String())

	payModel := db.NewPayModel(mysql.MySQLConn)
	_, err := payModel.Insert(ctx, &db.Pay{
		UserId:      user.Id,
		WalletId:    user.WalletId,
		PayId:       OutTradeNo,
		Subject:     fmt.Sprintf("%s：%.2f", config.BackendConfig.Coin.Name, float64(cny)/100.00),
		PayWay:      payWay,
		Cny:         cny,
		Get:         cny,
		TradeStatus: db.PayWait,
		CouponsId: sql.NullInt64{
			Valid: couponsID != 0,
			Int64: couponsID,
		},
	})
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	logger.Logger.WXInfo("收到用户自支付（%.2f）申请，单号：%s", float64(cny)/100.0, OutTradeNo)
	_ = LogMsg(true, "收到用户自支付（%.2f）申请，单号：%s", float64(cny)/100.0, OutTradeNo)

	audit.NewUserAudit(user.Id, "自支付（%s）发起成功（%.2f）", payWay, float64(cny)/100.0)
	return OutTradeNo, nil
}

func NewAdminPay(ctx context.Context, get int64, user *db.User, payWay string, subject string) errors.WTError {
	if payWay == alipay.PayWayPC || payWay == alipay.PayWayWap || payWay == wechatpay.PayWayNative || payWay == wechatpay.PayWayH5 || payWay == wechatpay.PayWayJSAPI {
		return BadPayWay.New()
	}

	if get < 0 { // 小于0，大于等于十万元
		return BadCNY.New()
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		payModel := db.NewPayModel(mysql.MySQLConn)
		_, err := payModel.FindByPayID(ctx, SelfpayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := SelfpayID(OutTradeNoUUID.String())

	pay := &db.Pay{
		UserId:      user.Id,
		WalletId:    user.WalletId,
		PayId:       OutTradeNo,
		Subject:     subject,
		PayWay:      payWay,
		Cny:         0,
		Get:         get,
		TradeStatus: db.PayWait,
		CouponsId: sql.NullInt64{
			Valid: false,
		},
		PayAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	_, err := balance.Pay(ctx, user, pay)
	if err != nil {
		logger.Logger.Error("add balance error: %s", err.Error())
		return errors.WarpQuick(err)
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	_, mysqlErr := payModel.Insert(ctx, pay)
	if mysqlErr != nil {
		logger.Logger.Error("mysql error: %s", mysqlErr.Error())
		return errors.WarpQuick(mysqlErr)
	}

	logger.Logger.WXInfo("管理员录入支付：%.2f", float64(pay.Cny)/100.0)
	_ = LogMsg(true, "管理员录入支付：%.2f", float64(pay.Cny)/100.0)

	sender.PhoneSendChange(pay.UserId, "余额（管理员录入充值）")
	sender.EmailSendChange(pay.UserId, "余额（管理员录入充值）")
	sender.MessageSendRecharge(pay.UserId, pay.Get, "管理员录入充值")
	sender.WxrobotSendRecharge(pay.UserId, pay.Get, "管理员录入充值")
	sender.FuwuhaoSendNotCnyPay(pay)
	audit.NewUserAudit(pay.UserId, "管理员录入支付（%s）发起成功（%.2f）", payWay, float64(pay.Cny)/100.0)
	return nil
}

func NewRefund(ctx context.Context, user *db.User, pay *db.Pay) errors.WTError {
	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap || pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		return BadPayWay.New()
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return errors.Errorf("pay lock error")
	}
	defer redis.ReleaseLock(keyPay)

	payModel := db.NewPayModel(mysql.MySQLConn)
	_, err := balance.PayRefund(ctx, user, pay)
	if errors.Is(err, balance.Insufficient) {
		return Insufficient.New()
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	mysqlErr := payModel.Update(ctx, pay)
	if mysqlErr != nil {
		return errors.WarpQuick(mysqlErr)
	}

	logger.Logger.WXInfo("收到用户自支付退款申请，单号：%s", pay.PayId)

	sender.PhoneSendChange(pay.UserId, "余额（支付自充值退款）")
	sender.EmailSendChange(pay.UserId, "余额（支付自充值退款）")
	sender.MessageSendRefundPay(pay.UserId, pay.Cny, "支付自充值退款")
	sender.WxrobotSendRefundPay(pay.UserId, pay.Cny, "支付自充值退款")
	sender.FuwuhaoSendRefundPay(pay)
	audit.NewUserAudit(pay.UserId, "支付自充值退款（%s）发起（%.2f）", pay.PayWay, float64(pay.Cny)/100.0)
	return nil
}
