package sender

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/fuwuhao"
	"gitee.com/wuntsong-auth/backend/src/model/db"
)

func FuwuhaoSendRefundPay(pay *db.Pay) {
	go func() {
		_ = fuwuhao.SendPay(context.Background(), pay.UserId, pay.Cny, fuwuhao.RefundPay, pay.PayId)
	}()
}

func FuwuhaoSendNotCnyPay(pay *db.Pay) {
	go func() {
		_ = fuwuhao.SendPay(context.Background(), pay.UserId, pay.Get, fuwuhao.Recharge, pay.PayId)
	}()
}

func FuwuhaoSendRecharge(pay *db.Pay) {
	go func() {
		_ = fuwuhao.SendPay(context.Background(), pay.UserId, pay.Cny, fuwuhao.Recharge, pay.PayId)
	}()
}

func FuwuhaoSendReturnDefray(defray *db.Defray) {
	go func() {
		_ = fuwuhao.SendPay(context.Background(), defray.UserId.Int64, defray.Price, fuwuhao.RefundDefray, defray.DefrayId)
	}()
}

func FuwuhaoSendDefray(defray *db.Defray) {
	go func() {
		_ = fuwuhao.SendPay(context.Background(), defray.UserId.Int64, defray.Price, fuwuhao.Defray, defray.DefrayId)
	}()
}

func FuwuhaoSendDeleteUser(userID int64) {
	go func() {
		_ = fuwuhao.SendDeleteUser(context.Background(), userID)
	}()
}

func FuwuhaoSendChange(userID int64, project string) {
	go func() {
		_ = fuwuhao.SendProject(context.Background(), userID, fmt.Sprintf("信息（%s）变更已完成", project))
	}()
}

func FuwuhaoSendLoginCenter(userID int64) {
	go func() {
		_ = fuwuhao.SendLoginSuccess(context.Background(), userID, config.BackendConfig.User.ReadableName, config.BackendConfig.User.Url)
	}()
}

func FuwuhaoSendOauth2(userID int64, webName string) {
	go func() {
		_ = fuwuhao.SendOauth2(context.Background(), userID, webName, "")
	}()
}
