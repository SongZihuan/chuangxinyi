package alipay

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
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/SuperH-0630/gopay/alipay"
	"github.com/wuntsong-org/wterrors"
	"net/http"
	"time"
)

func Notification(w http.ResponseWriter, r *http.Request) {
	record := record.GetRecord(r.Context())

	notifyReq, err := alipay.ParseNotifyToBodyMap(r) // c.Request 是 gin 框架的写法
	if err != nil {
		utils.Forbidden(w, r, err, false, record.RequestsID)
		return
	}

	ok, err := alipay.VerifySignWithCert(AlipayPublicCert, notifyReq)
	if err != nil {
		record.Msg = err.Error()
		utils.Forbidden(w, r, err, false, record.RequestsID)
		return
	} else if !ok {
		record.Msg = "not alipay sign"
		utils.Forbidden(w, r, nil, false, record.RequestsID)
		return
	}

	switch notifyReq.Get("notify_type") {
	case "trade_status_sync":
		payModel := db.NewPayModel(mysql.MySQLConn)
		userModel := db.NewUserModel(mysql.MySQLConn)

		payID := notifyReq.Get("out_trade_no")
		pay, err := payModel.FindByPayID(r.Context(), payID)
		if errors.Is(err, db.ErrNotFound) {
			record.Msg = "pay not found"
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		} else if err != nil {
			record.Msg = err.Error()
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		}

		user, err := userModel.FindOneByIDWithoutDelete(r.Context(), pay.UserId)
		if errors.Is(err, db.ErrNotFound) {
			record.Msg = "user not found"
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		} else if err != nil {
			record.Msg = err.Error()
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		}

		keyPay := fmt.Sprintf("pay:%s", pay.PayId)
		if !redis.AcquireLockMore(r.Context(), keyPay, time.Minute*2) {
			record.Msg = fmt.Sprintf("can not get lock for %s", pay.PayId)
			utils.Forbidden(w, r, nil, true, record.RequestsID)
			return
		}
		defer redis.ReleaseLock(keyPay)

		tradeStatus := notifyReq.Get("trade_status")
		switch tradeStatus {
		default:
			pay.TradeStatus = db.PayWait
		case "WAIT_BUYER_PAY":
			pay.TradeStatus = db.PayWait
		case "TRADE_CLOSED":
			pay.TradeStatus = db.PayClose
		case "TRADE_SUCCESS": // 除了PayWait不改变状态
			payTime, err := time.Parse("2006-01-02 15:04:05", notifyReq.Get("gmt_payment"))
			if err != nil {
				record.Msg = err.Error()
				utils.Forbidden(w, r, err, true, record.RequestsID)
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

			buyerID := notifyReq.Get("buyer_id")
			if len(notifyReq.Get("buyer_id")) == 0 {
				buyerID = notifyReq.Get("buyer_open_id")
			}

			pay.BuyerId = sql.NullString{
				Valid:  true,
				String: buyerID,
			}
			pay.TradeNo = sql.NullString{
				Valid:  true,
				String: notifyReq.Get("trade_no"),
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
					utils.Forbidden(w, r, err, true, record.RequestsID)
					return
				}
				logger.Logger.WXInfo("支付宝支付成功：%.2f", float64(pay.Cny)/100.0)
				_ = LogMsg(true, "支付宝支付成功：%.2f", float64(pay.Cny)/100.0)
			}
		case "TRADE_FINISHED": // 不可退款
			pay.TradeStatus = db.PayFinish
		}

		defer func() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success"))
		}()

		err = payModel.Update(r.Context(), pay)
		if err != nil {
			record.Msg = err.Error()
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		}

		if tradeStatus == "TRADE_SUCCESS" {
			sender.PhoneSendChange(pay.UserId, "余额（支付宝充值）")
			sender.EmailSendChange(pay.UserId, "余额（支付宝充值）")
			sender.MessageSendRecharge(pay.UserId, pay.Cny, "支付宝")
			sender.WxrobotSendRecharge(pay.UserId, pay.Cny, "支付宝")
			sender.FuwuhaoSendRecharge(pay)
			audit.NewUserAudit(pay.UserId, "支付宝充值已到账（%.2f）", pay.Cny)
		}
	default:
		record.Msg = "bad notify type"
		utils.Forbidden(w, r, nil, true, record.RequestsID)
		return
	}

}

func NotificationWangguan(w http.ResponseWriter, r *http.Request) {
	record := record.GetRecord(r.Context())

	notifyReq, err := alipay.ParseNotifyToBodyMap(r) // c.Request 是 gin 框架的写法
	if err != nil {
		utils.Forbidden(w, r, err, false, record.RequestsID)
		return
	}

	ok, err := alipay.VerifySignWithCert(AlipayPublicCert, notifyReq)
	if err != nil {
		record.Msg = err.Error()
		utils.Forbidden(w, r, err, false, record.RequestsID)
		return
	} else if !ok {
		record.Msg = "not alipay sign"
		utils.Forbidden(w, r, nil, false, record.RequestsID)
		return
	}

	switch notifyReq.Get("msg_method") {
	case "alipay.fund.trans.order.changed":
		contentData := struct {
			WithdrawID string `json:"out_biz_no"`
			Date       string `json:"pay_date"`
		}{}
		contentString := notifyReq.Get("biz_content")

		jsonErr := utils.JsonUnmarshal([]byte(contentString), &contentData)
		if err != nil {
			record.Msg = jsonErr.Error()
			utils.Forbidden(w, r, jsonErr, true, record.RequestsID)
			return
		}

		withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
		withdraw, err := withdrawModel.FindByWithdrawID(r.Context(), contentData.WithdrawID)
		if errors.Is(err, db.ErrNotFound) {
			record.Msg = "withdraw not found"
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		} else if err != nil {
			record.Msg = err.Error()
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		}

		withdraw.Status = db.WithdrawOK
		payDate, err := time.Parse("2006-01-02 15:04:05", contentData.Date)
		if err == nil {
			withdraw.PayAt = sql.NullTime{
				Valid: true,
				Time:  payDate,
			}
		}

		err = withdrawModel.Update(r.Context(), withdraw)
		if err != nil {
			record.Msg = err.Error()
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		}

		defer func() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success"))
		}()

		sender.PhoneSendChange(withdraw.UserId, "提现额度（支付宝快捷提现）")
		sender.EmailSendChange(withdraw.UserId, "提现额度（支付宝快捷提现）")
		sender.MessageSendWithdraw(withdraw.UserId, withdraw.Cny, "支付宝")
		audit.NewUserAudit(withdraw.UserId, "支付宝提现已到账（%.2f）", float64(withdraw.Cny)/100.0)
		logger.Logger.WXInfo("支付宝提现成功：%.2f", float64(withdraw.Cny)/100.0)
		_ = LogMsg(true, "支付宝提现成功：%.2f", float64(withdraw.Cny)/100.0)
	case "alipay.trade.refund.depositback.completed":
		contentData := struct {
			PayID string `json:"out_trade_no"`
		}{}
		contentString := notifyReq.Get("biz_content")

		jsonErr := utils.JsonUnmarshal([]byte(contentString), &contentData)
		if err != nil {
			record.Msg = jsonErr.Error()
			utils.Forbidden(w, r, jsonErr, true, record.RequestsID)
			return
		}

		payModel := db.NewPayModel(mysql.MySQLConn)
		pay, err := payModel.FindByPayID(r.Context(), contentData.PayID)
		if errors.Is(err, db.ErrNotFound) {
			record.Msg = "pay not found"
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		} else if err != nil {
			record.Msg = err.Error()
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		}

		pay.TradeStatus = db.PaySuccessRefund

		err = payModel.Update(r.Context(), pay)
		if err != nil {
			record.Msg = err.Error()
			utils.Forbidden(w, r, err, true, record.RequestsID)
			return
		}

		defer func() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success"))
		}()

		sender.PhoneSendChange(pay.UserId, "余额（支付宝退款）")
		sender.EmailSendChange(pay.UserId, "余额（支付宝退款）")
		sender.MessageSendRefundPay(pay.UserId, pay.Cny, "支付宝")
		sender.WxrobotSendRefundPay(pay.UserId, pay.Cny, "支付宝")
		sender.FuwuhaoSendRefundPay(pay)
		audit.NewUserAudit(pay.UserId, "支付宝退款已到账（%.2f）", float64(pay.Cny)/100.0)
		logger.Logger.WXInfo("支付宝退款成功：%.2f", float64(pay.Cny)/100.0)
		_ = LogMsg(true, "支付宝退款成功：%.2f", float64(pay.Cny)/100.0)
	default:
		record.Msg = "bad msg_method"
		utils.Forbidden(w, r, nil, true, record.RequestsID)
		return
	}

}

func LogMsg(atall bool, text string, args ...any) errors.WTError {
	return logger.WxRobotSendNotRecord(config.BackendConfig.WXRobot.PayLog, fmt.Sprintf(text, args...), atall)
}
