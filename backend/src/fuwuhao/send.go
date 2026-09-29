package fuwuhao

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/fastwego/offiaccount/apis/message/template"
	"github.com/wuntsong-org/wterrors"
	"strings"
	"time"
)

const Recharge = 1
const Defray = 2
const RefundPay = 3
const RefundDefray = 4
const Back = 5
const Withdraw = 6

func SendDeleteUser(ctx context.Context, userID int64) errors.WTError {
	userModel := db.NewUserModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)

	w, err := wechatModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if !w.Fuwuhao.Valid {
		return nil
	}

	_, err = userModel.FindOneByIDWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	p, err := phoneModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	return Send(context.Background(),
		w.Fuwuhao.String,
		config.BackendConfig.FuWuHao.UserDelete.TemplateID,
		config.BackendConfig.FuWuHao.UserDelete.Url,
		map[string]TemplateValue{
			"thing2": {
				Value: p.Phone,
			},
			"phrase3": {
				Value: "已注销",
			},
			"time4": {
				Value: time.Now().Format("2006-01-02 15:04:05"),
			},
		},
		warp.UserCenterWebsite)
}

func SendOauth2(ctx context.Context, userID int64, webName string, url string) errors.WTError {
	userModel := db.NewUserModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)

	w, err := wechatModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if !w.Fuwuhao.Valid {
		return nil
	}

	_, err = userModel.FindOneByIDWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	p, err := phoneModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	return Send(context.Background(),
		w.Fuwuhao.String,
		config.BackendConfig.FuWuHao.Oauth2.TemplateID,
		url,
		map[string]TemplateValue{
			"thing1": {
				Value: p.Phone,
			},
			"thing2": {
				Value: webName,
			},
		},
		warp.UserCenterWebsite)
}

func SendLoginFail(ctx context.Context, userID int64, reason string) errors.WTError {
	userModel := db.NewUserModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)

	w, err := wechatModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if !w.Fuwuhao.Valid {
		return nil
	}

	_, err = userModel.FindOneByIDWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	p, err := phoneModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	return Send(context.Background(),
		w.Fuwuhao.String,
		config.BackendConfig.FuWuHao.LoginFail.TemplateID,
		config.BackendConfig.FuWuHao.LoginFail.Url,
		map[string]TemplateValue{
			"thing2": {
				Value: p.Phone,
			},
			"time3": {
				Value: time.Now().Format("2006-01-02 15:04:05"),
			},
			"thing6": {
				Value: reason,
			},
		},
		warp.UserCenterWebsite)
}

func SendLoginSuccess(ctx context.Context, userID int64, webName string, url string) errors.WTError {
	userModel := db.NewUserModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)

	w, err := wechatModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if !w.Fuwuhao.Valid {
		return nil
	}

	_, err = userModel.FindOneByIDWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	p, err := phoneModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	return Send(context.Background(),
		w.Fuwuhao.String,
		config.BackendConfig.FuWuHao.LoginSuccess.TemplateID,
		url,
		map[string]TemplateValue{
			"thing6": {
				Value: p.Phone,
			},
			"time3": {
				Value: time.Now().Format("2006-01-02 15:04:05"),
			},
			"thing4": {
				Value: webName,
			},
		},
		warp.UserCenterWebsite)
}

func SendPay(ctx context.Context, userID int64, cny int64, payType int, id string) errors.WTError {
	userModel := db.NewUserModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)

	w, err := wechatModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if !w.Fuwuhao.Valid {
		return nil
	}

	_, err = userModel.FindOneByIDWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	p, err := phoneModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	var pt string
	if payType == Recharge {
		pt = "充值订单"
	} else if payType == Defray {
		pt = "支付订单"
	} else if payType == RefundPay {
		pt = "退款订单"
	} else if payType == RefundDefray {
		pt = "退款消费"
	} else if payType == Back {
		pt = "优惠返现"
	} else if payType == Withdraw {
		pt = "提现订单"
	} else {
		pt = "充值订单"
	}

	return Send(context.Background(),
		w.Fuwuhao.String,
		config.BackendConfig.FuWuHao.Pay.TemplateID,
		config.BackendConfig.FuWuHao.Pay.Url,
		map[string]TemplateValue{
			"thing1": {
				Value: p.Phone,
			},
			"const2": {
				Value: pt,
			},
			"character_string3": {
				Value: id,
			},
			"amount4": {
				Value: fmt.Sprintf("%.2f", float64(cny)/100.00),
			},
		},
		warp.UserCenterWebsite)
}

func SendPing(fuwuhao string) errors.WTError {
	return Send(context.Background(),
		fuwuhao,
		config.BackendConfig.FuWuHao.Project.TemplateID,
		config.BackendConfig.FuWuHao.Project.Url,
		map[string]TemplateValue{
			"thing9": {
				Value: "Ping业务处理完成",
			},
			"thing4": {
				Value: "未知",
			},
			"time7": {
				Value: time.Now().Format("2006-01-02 15:04:05"),
			},
		},
		warp.UserCenterWebsite)
}

func SendProject(ctx context.Context, userID int64, project string) errors.WTError {
	userModel := db.NewUserModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)

	w, err := wechatModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if !w.Fuwuhao.Valid {
		return nil
	}

	_, err = userModel.FindOneByIDWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	p, err := phoneModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	return Send(context.Background(),
		w.Fuwuhao.String,
		config.BackendConfig.FuWuHao.Project.TemplateID,
		config.BackendConfig.FuWuHao.Project.Url,
		map[string]TemplateValue{
			"thing9": {
				Value: project,
			},
			"thing4": {
				Value: p.Phone,
			},
			"time7": {
				Value: time.Now().Format("2006-01-02 15:04:05"),
			},
		},
		warp.UserCenterWebsite)
}

func SendVal(ctx context.Context, templateID string, url string, fuwuhao string, val map[string]string, senderID int64) errors.WTError {
	valTemplateValue := make(map[string]TemplateValue, len(val))
	for k, v := range val {
		valTemplateValue[k] = TemplateValue{Value: v}
	}

	return Send(ctx, fuwuhao, templateID, url, valTemplateValue, senderID)
}

func FuwuhaoUserID(userID string) string {
	return strings.Replace(userID, "-", "", -1)
}

func SendBindSuccess(ctx context.Context, fuwuhao string, userID int64) errors.WTError {
	userModel := db.NewUserModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)

	u, err := userModel.FindOneByIDWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return nil
	}

	p, err := phoneModel.FindByUserID(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if err != nil {
		return nil
	}

	return Send(context.Background(),
		fuwuhao,
		config.BackendConfig.FuWuHao.Register.TemplateID,
		config.BackendConfig.FuWuHao.Register.Url,
		map[string]TemplateValue{
			"character_string6": {
				Value: FuwuhaoUserID(u.Uid),
			},
			"phone_number1": {
				Value: p.Phone,
			},
			"time3": {
				Value: u.CreateAt.Format("2006-01-02 15:04:05"),
			},
		},
		0)
}

func Send(ctx context.Context, toUser string, templateID string, url string, data map[string]TemplateValue, senderID int64) errors.WTError {
	dataEasy := make(map[string]interface{}, len(data))
	for k, v := range data {
		dataEasy[k] = v.Value
	}
	dataEasyByte, err := utils.JsonMarshal(dataEasy)
	if err != nil {
		return errors.WarpQuick(err)
	}

	fuwuhaoMessageModel := db.NewFuwuhaoMessageModel(mysql.MySQLConn)
	fuwuhaoMessage := &db.FuwuhaoMessage{
		OpenId:   toUser,
		Template: templateID,
		Url:      url,
		Val:      string(dataEasyByte),
		SenderId: senderID,
		Success:  true,
	}

	d, err := utils.JsonMarshal(TemplateMsgReq{
		ToUser:     toUser,
		TemplateId: templateID,
		Url:        url,
		Data:       data,
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	r, sendErr := template.Send(OffiAccount, d)
	if sendErr != nil {
		fuwuhaoMessage.Success = false
		fuwuhaoMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: sendErr.Error(),
		}
		_, _ = fuwuhaoMessageModel.Insert(ctx, fuwuhaoMessage)
		return errors.WarpQuick(sendErr)
	}

	var resp Resp
	err = utils.JsonUnmarshal(r, &resp)
	if err != nil {
		return errors.WarpQuick(err)
	}

	if resp.Errcode != 0 {
		fuwuhaoMessage.Success = false
		fuwuhaoMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: fmt.Sprintf("%d: %s", resp.Errcode, resp.Errmsg),
		}
		_, _ = fuwuhaoMessageModel.Insert(ctx, fuwuhaoMessage)
		return errors.Errorf("%d: %s", resp.Errcode, resp.Errmsg)
	}

	_, _ = fuwuhaoMessageModel.Insert(ctx, fuwuhaoMessage)
	return nil
}
