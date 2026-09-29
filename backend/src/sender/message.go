package sender

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/msg"
)

func sendMessage(userID int64, title string, content string) {
	_ = msg.SendMessage(userID, title, content, config.BackendConfig.Message.Sender, 0, config.BackendConfig.Message.SenderLink)
}

func MessageSendRegister(userID int64, phone string) {
	go sendMessage(userID, "欢迎信", fmt.Sprintf("你已成功注册账号，欢迎投入使用！手机号：%s。", phone))
}

func MessageSendSonRegister(userID int64, phone string) {
	go sendMessage(userID, "欢迎信", fmt.Sprintf("你已成功注册子账号，欢迎投入使用！手机号：%s。", phone))
}

func MessageSendChange(userID int64, project string) {
	go sendMessage(userID, "用户信息变更提示", fmt.Sprintf("你的用户信息已变更，变更内容为：%s。请知悉！", project))
}

func MessageSendRefundPay(userID int64, cny int64, payWay string) {
	go sendMessage(userID, "用户退款提示", fmt.Sprintf("你的退款已到账，退款额度：%.2f，退款方式：%s。", float64(cny)/100.0, payWay))
}

func MessageSendWithdraw(userID int64, cny int64, payWay string) {
	go sendMessage(userID, "用户提现提示", fmt.Sprintf("你的提现已到账，提现额度：%.2f，提现方式：%s。", float64(cny)/100.0, payWay))
}

func MessageSendRecharge(userID int64, cny int64, payWay string) {
	go sendMessage(userID, "用户充值提示", fmt.Sprintf("你的充值已到账，充值额度：%.2f，充值方式：%s。", float64(cny)/100.0, payWay))
}

func MessageSendBack(userID int64, cny int64, subject string) {
	go sendMessage(userID, "用户返现提示", fmt.Sprintf("你已获得返现，返现额度：%.2f，充值方式：%s。", float64(cny)/100.0, subject))
}

func MessageSendPay(userID int64, cny int64, subject string) {
	go sendMessage(userID, "用户支付提示", fmt.Sprintf("你已成功支付，消耗额度：%.2f，商品名称：%s。", float64(cny)/100.0, subject))
}

func MessageSendPayAdmin(userID int64, cny int64, subject string) {
	go sendMessage(userID, "用户支付提示", fmt.Sprintf("你已成功扣费，消耗额度：%.2f，商品名称：%s。", float64(cny)/100.0, subject))
}

func MessageSendPayReturn(userID int64, cny int64, subject string) {
	go sendMessage(userID, "用户支付提示", fmt.Sprintf("你已成功退款，返回额度：%.2f，商品名称：%s。", float64(cny)/100.0, subject))
}

func MessageSendLoginCenter(userID int64, ctx context.Context) {
	geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok {
		geo = "未知"
	}

	ip, ok := ctx.Value("X-Real-IP").(string)
	if !ok {
		ip = "未知"
	}

	go sendMessage(userID, "用户登录提示", fmt.Sprintf("用户成功登录%s，登录IP为：%s，登录地为：%s", config.BackendConfig.User.ReadableName, ip, geo))
}

func MessageSendOauth2(userID int64, webName string, ctx context.Context) {
	geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok {
		geo = "未知"
	}

	ip, ok := ctx.Value("X-Real-IP").(string)
	if !ok {
		ip = "未知"
	}

	go sendMessage(userID, "用户授权提示", fmt.Sprintf("你以成功授权登录网站（%s），授权IP为：%s，授权地为：%s", webName, ip, geo))
}

func MessageSendUncle(userID string, uncleID int64) {
	go sendMessage(uncleID, "新增协作账号提醒", fmt.Sprintf("账号（%s）将你设置为协作人，请前往同意后可共享信息和控制其账号", userID))
}

func MessageSend(userID int64, title, content string) {
	go sendMessage(userID, title, content)
}
