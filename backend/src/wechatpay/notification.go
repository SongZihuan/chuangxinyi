package wechatpay

import (
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/coupons"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/record"
	utils2 "gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/SuperH-0630/gopay"
	"github.com/SuperH-0630/gopay/wechat/v3"
	"github.com/wuntsong-org/wterrors"
	"net/http"
	"time"
)

type Resource struct {
	Algorithm      string `json:"algorithm"`
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
}

type Event struct {
	ID           string   `json:"id"`
	ResourceType string   `json:"resource_type"`
	EventType    string   `json:"event_type"`
	Summary      string   `json:"summary"`
	Resource     Resource `json:"resource"`
}

type Payer struct {
	OpenID string `json:"openid"`
}

type Trade struct {
	TransactionID string `json:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no"`
	Payer         Payer  `json:"payer"`
	SuccessTime   string `json:"success_time"`
}

func Notification(w http.ResponseWriter, r *http.Request) {
	record := record.GetRecord(r.Context())

	notifyReq, err := wechat.V3ParseNotify(r)
	if err != nil {
		record.Msg = err.Error()
		utils2.Forbidden(w, r, err, false, record.RequestsID)
		return
	}

	// 获取微信平台证书
	certMap := WeChatPayClient.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		record.Msg = err.Error()
		utils2.Forbidden(w, r, err, false, record.RequestsID)
		return
	}

	switch notifyReq.EventType {
	case "REFUND.SUCCESS":
		result, err := notifyReq.DecryptRefundCipherText(config.BackendConfig.WeChatPay.MchAPIv3Key)
		if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		payModel := db.NewPayModel(mysql.MySQLConn)
		userModel := db.NewUserModel(mysql.MySQLConn)

		payID := result.OutTradeNo
		pay, err := payModel.FindByPayID(r.Context(), payID)
		if errors.Is(err, db.ErrNotFound) {
			record.Msg = "pay not found"
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		} else if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		user, err := userModel.FindOneByIDWithoutDelete(r.Context(), pay.UserId)
		if errors.Is(err, db.ErrNotFound) {
			record.Msg = "user not found"
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		} else if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		} else if user.WalletId != pay.WalletId {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		switch result.RefundStatus {
		case "SUCCESS":
			if pay.TradeStatus == db.PayWaitRefund {
				logger.Logger.WXInfo("微信支付退款成功：%.2f", float64(pay.Cny)/100.0)
				_ = LogMsg(true, "微信支付退款成功：%.2f", float64(pay.Cny)/100.0)
			}
			pay.TradeStatus = db.PaySuccessRefund
		case "CLOSE":
			pay.TradeStatus = db.PayCloseRefund
			_, err = balance.PayRefundFail(r.Context(), user, pay, db.PayCloseRefund)
		case "ABNORMAL":
			pay.TradeStatus = db.PayWaitRefund
			logger.Logger.Error("微信退款异常，需要前往处理：%s", pay.PayId)
		default:
			pay.TradeStatus = db.PayWaitRefund
		}

		err = payModel.Update(r.Context(), pay)
		if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		defer func() {
			resByte, err := utils2.JsonMarshal(wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
			if err != nil {
				record.Msg = err.Error()
				utils2.Forbidden(w, r, err, false, record.RequestsID)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(resByte)
		}()

		if result.RefundStatus == "SUCCESS" {
			sender.PhoneSendChange(pay.UserId, "余额（微信支付退款）")
			sender.EmailSendChange(pay.UserId, "余额（微信支付退款）")
			sender.MessageSendRefundPay(pay.UserId, pay.Cny, "微信支付")
			sender.WxrobotSendRefundPay(pay.UserId, pay.Cny, "微信支付")
			sender.FuwuhaoSendRefundPay(pay)
			audit.NewUserAudit(pay.UserId, "微信支付退款已到账（%.2f）", float64(pay.Cny)/100.0)
		}
	case "TRANSACTION.SUCCESS":
		result, err := notifyReq.DecryptCipherText(config.BackendConfig.WeChatPay.MchAPIv3Key)
		if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		payModel := db.NewPayModel(mysql.MySQLConn)
		userModel := db.NewUserModel(mysql.MySQLConn)

		payID := result.OutTradeNo
		pay, err := payModel.FindByPayID(r.Context(), payID)
		if errors.Is(err, db.ErrNotFound) {
			record.Msg = "pay not found"
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		} else if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		user, err := userModel.FindOneByIDWithoutDelete(r.Context(), pay.UserId)
		if errors.Is(err, db.ErrNotFound) {
			record.Msg = "user not found"
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		} else if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		} else if user.WalletId != pay.WalletId {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		keyPay := fmt.Sprintf("pay:%s", pay.PayId)
		if !redis.AcquireLockMore(r.Context(), keyPay, time.Minute*2) {
			record.Msg = fmt.Sprintf("can not get lock for %s", pay.PayId)
			utils2.Forbidden(w, r, nil, true, record.RequestsID)
			return
		}
		defer redis.ReleaseLock(keyPay)

		payTime, err := time.Parse(time.RFC3339, result.SuccessTime)
		if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		var get int64
		if pay.CouponsId.Valid {
			get, err = coupons.Recharge(r.Context(), pay.CouponsId.Int64, pay.Cny)
			if err != nil {
				get = pay.Get
			}
		} else {
			get = pay.Get
		}

		pay.BuyerId = sql.NullString{
			Valid:  true,
			String: result.Payer.Openid,
		}
		pay.TradeNo = sql.NullString{
			Valid:  true,
			String: result.TransactionId,
		}
		pay.PayAt = sql.NullTime{
			Valid: true,
			Time:  payTime,
		}
		pay.Get = get

		if pay.TradeStatus == db.PayWait {
			_, err = balance.Pay(r.Context(), user, pay)
			if err != nil {
				record.Msg = err.Error()
				utils2.Forbidden(w, r, err, false, record.RequestsID)
				return
			}
			logger.Logger.WXInfo("微信支付支付成功：%.2f", float64(pay.Cny)/100.0)
			_ = LogMsg(true, "微信支付支付成功：%.2f", float64(pay.Cny)/100.0)
		}

		err = payModel.Update(r.Context(), pay)
		if err != nil {
			record.Msg = err.Error()
			utils2.Forbidden(w, r, err, false, record.RequestsID)
			return
		}

		defer func() {
			resByte, err := utils2.JsonMarshal(wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
			if err != nil {
				record.Msg = err.Error()
				utils2.Forbidden(w, r, err, false, record.RequestsID)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(resByte)
		}()

		sender.PhoneSendChange(pay.UserId, "余额（微信支付充值）")
		sender.EmailSendChange(pay.UserId, "余额（微信支付充值）")
		sender.MessageSendRecharge(pay.UserId, pay.Cny, "微信支付")
		sender.WxrobotSendRecharge(pay.UserId, pay.Cny, "微信支付")
		sender.FuwuhaoSendRecharge(pay)
		audit.NewUserAudit(pay.UserId, "微信支付充值已到账（%.2f）", float64(pay.Cny)/100.0)
	default:
		record.Msg = "bad event type"
		utils2.Forbidden(w, r, nil, false, record.RequestsID)
		return
	}
}

func LogMsg(atall bool, text string, args ...any) errors.WTError {
	return logger.WxRobotSendNotRecord(config.BackendConfig.WXRobot.PayLog, fmt.Sprintf(text, args...), atall)
}
