package sender

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/wxrobot"
	"github.com/wuntsong-org/wterrors"
)

func sendWxrobot(userID int64, title string, content string) {
	wxrobotModel := db.NewWxrobotModel(mysql.MySQLConn)
	w, err := wxrobotModel.FindByUserID(context.Background(), userID)
	if errors.Is(err, db.ErrNotFound) {
		return
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	if !w.Webhook.Valid {
		return
	}

	_ = wxrobot.Send(context.Background(), w.Webhook.String, fmt.Sprintf("%s: %s", title, content), false, 0, config.BackendConfig.WXRobot.Sender)
}

func WxrobotSendDelete(userID int64) {
	go sendWxrobot(userID, "用户信息变更提示", fmt.Sprintf("你的用户已注销。请知悉！"))
}

func WxrobotSendChange(userID int64, project string) {
	go sendWxrobot(userID, "用户信息变更提示", fmt.Sprintf("你的用户信息已变更，变更内容为：%s。请知悉！", project))
}

func WxrobotSendRefundPay(userID int64, cny int64, payWay string) {
	go sendWxrobot(userID, "用户充值提示", fmt.Sprintf("你的退款已到账，退款额度：%.2f，退款方式：%s。", float64(cny)/100.0, payWay))
}

func WxrobotSendRecharge(userID int64, cny int64, subject string) {
	go sendWxrobot(userID, "用户充值提示", fmt.Sprintf("你已获得返现，返现额度：%.2f，充值方式：%s。", float64(cny)/100.0, subject))
}

func WxrobotSendPay(userID int64, cny int64, subject string) {
	go sendWxrobot(userID, "用户支付提示", fmt.Sprintf("你已成功支付，消耗额度：%.2f，商品名称：%s。", float64(cny)/100.0, subject))
}

func WxrobotSendPayAdmin(userID int64, cny int64, subject string) {
	go sendWxrobot(userID, "用户支付提示", fmt.Sprintf("你已成功扣费，消耗额度：%.2f，商品名称：%s。", float64(cny)/100.0, subject))
}

func WxrobotSendPayReturn(userID int64, cny int64, subject string) {
	go sendWxrobot(userID, "用户支付提示", fmt.Sprintf("你已成功退款，返回额度：%.2f，商品名称：%s。", float64(cny)/100.0, subject))
}

func WxrobotSendLoginCenter(userID int64, ctx context.Context) {
	geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok {
		geo = "未知"
	}

	ip, ok := ctx.Value("X-Real-IP").(string)
	if !ok {
		ip = "未知"
	}

	go sendWxrobot(userID, "用户登录提示", fmt.Sprintf("用户成功登录%s，登录IP为：%s，登录地为：%s", config.BackendConfig.User.ReadableName, ip, geo))
}

func WxrobotSendOauth2(userID int64, webName string, ctx context.Context) {
	geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok {
		geo = "未知"
	}

	ip, ok := ctx.Value("X-Real-IP").(string)
	if !ok {
		ip = "未知"
	}

	go sendWxrobot(userID, "用户授权提示", fmt.Sprintf("你以成功授权登录网站（%s），授权IP为：%s，授权地为：%s", webName, ip, geo))
}
