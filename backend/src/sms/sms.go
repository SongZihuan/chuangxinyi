package sms

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wuntsong-org/wterrors"
	"time"
)

var SMSClient *dysmsapi20170525.Client

var SMSSendLimit = errors.NewClass("send limit")
var SMSError = errors.NewClass("sms error")

func InitSMS() errors.WTError {
	if len(config.BackendConfig.Aliyun.AccessKeyId) == 0 {
		return errors.Errorf("aliyun AccessKeyId must be given")
	}

	if len(config.BackendConfig.Aliyun.AccessKeySecret) == 0 {
		return errors.Errorf("aliyun AccessKeySecret must be given")
	}

	if len(config.BackendConfig.Aliyun.ImportCode.Sig) == 0 {
		return errors.Errorf("aliyun import code sig must be given")
	}

	if len(config.BackendConfig.Aliyun.ImportCode.Sig) == 0 {
		return errors.Errorf("aliyun import code sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Code.Sig) == 0 {
		return errors.Errorf("aliyun code sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Code.Sig) == 0 {
		return errors.Errorf("aliyun code sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Change.Sig) == 0 {
		return errors.Errorf("aliyun change sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Change.Sig) == 0 {
		return errors.Errorf("aliyun change sig must be given")
	}

	if len(config.BackendConfig.Aliyun.ChangePhone.Sig) == 0 {
		return errors.Errorf("aliyun change phone sig must be given")
	}

	if len(config.BackendConfig.Aliyun.ChangePhone.Sig) == 0 {
		return errors.Errorf("aliyun change phone sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Delete.Sig) == 0 {
		return errors.Errorf("aliyun delete sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Delete.Sig) == 0 {
		return errors.Errorf("aliyun delete sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Register.Sig) == 0 {
		return errors.Errorf("aliyun register sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Register.Sig) == 0 {
		return errors.Errorf("aliyun register sig must be given")
	}

	AccessKeyId := tea.String(config.BackendConfig.Aliyun.AccessKeyId)
	AccessKeySecret := tea.String(config.BackendConfig.Aliyun.AccessKeySecret)
	Endpoint := tea.String("dysmsapi.aliyuncs.com")

	var redisErr error
	SMSClient, redisErr = dysmsapi20170525.NewClient(&openapi.Config{
		AccessKeyId:     AccessKeyId,
		AccessKeySecret: AccessKeySecret,
		Endpoint:        Endpoint,
	})
	if redisErr == nil {
		return nil
	}

	return errors.WarpQuick(redisErr)
}

func SendCode(code int64, phone string) errors.WTError {
	if len(config.BackendConfig.Aliyun.Code.Sig) == 0 {
		return errors.Errorf("aliyun code sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Code.Sig) == 0 {
		return errors.Errorf("aliyun code sig must be given")
	}

	if code > 999999 || code < 0 {
		return errors.Errorf("bad code")
	}

	TemplateParam := map[string]string{
		"code": fmt.Sprintf("%06d", code),
	}

	return Send(context.Background(), TemplateParam, config.BackendConfig.Aliyun.Code.Sig, config.BackendConfig.Aliyun.Code.Template, phone, warp.UserCenterWebsite)
}

func SendImportCode(code int64, phone string) errors.WTError {
	if len(config.BackendConfig.Aliyun.ImportCode.Sig) == 0 {
		return errors.Errorf("aliyun import code sig must be given")
	}

	if len(config.BackendConfig.Aliyun.ImportCode.Sig) == 0 {
		return errors.Errorf("aliyun import code sig must be given")
	}

	if code > 999999 || code < 0 {
		return errors.Errorf("bad code")
	}

	TemplateParam := map[string]string{
		"code": fmt.Sprintf("%06d", code),
	}

	return Send(context.Background(), TemplateParam, config.BackendConfig.Aliyun.ImportCode.Sig, config.BackendConfig.Aliyun.ImportCode.Template, phone, warp.UserCenterWebsite)
}

func SendChange(changeProject string, phone string) errors.WTError {
	if len(config.BackendConfig.Aliyun.Change.Sig) == 0 {
		return errors.Errorf("aliyun change sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Change.Sig) == 0 {
		return errors.Errorf("aliyun change sig must be given")
	}

	TemplateParam := map[string]string{
		"project": changeProject,
	}

	return Send(context.Background(), TemplateParam, config.BackendConfig.Aliyun.Change.Sig, config.BackendConfig.Aliyun.Change.Template, phone, warp.UserCenterWebsite)
}

func SendChangePhone(newPhone string, phone string) errors.WTError {
	if len(config.BackendConfig.Aliyun.ChangePhone.Sig) == 0 {
		return errors.Errorf("aliyun change phone sig must be given")
	}

	if len(config.BackendConfig.Aliyun.ChangePhone.Sig) == 0 {
		return errors.Errorf("aliyun change phone sig must be given")
	}

	TemplateParam := map[string]string{
		"phone": newPhone,
	}

	return Send(context.Background(), TemplateParam, config.BackendConfig.Aliyun.ChangePhone.Sig, config.BackendConfig.Aliyun.ChangePhone.Template, phone, warp.UserCenterWebsite)
}

func SendDelete(phone string) errors.WTError {
	if len(config.BackendConfig.Aliyun.Delete.Sig) == 0 {
		return errors.Errorf("aliyun delete sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Delete.Sig) == 0 {
		return errors.Errorf("aliyun delete sig must be given")
	}

	TemplateParam := map[string]string{}

	return Send(context.Background(), TemplateParam, config.BackendConfig.Aliyun.Delete.Sig, config.BackendConfig.Aliyun.Delete.Template, phone, warp.UserCenterWebsite)
}

func SendRegister(phone string) errors.WTError {
	if len(config.BackendConfig.Aliyun.Register.Sig) == 0 {
		return errors.Errorf("aliyun register sig must be given")
	}

	if len(config.BackendConfig.Aliyun.Register.Sig) == 0 {
		return errors.Errorf("aliyun register sig must be given")
	}

	TemplateParam := map[string]string{}

	return Send(context.Background(), TemplateParam, config.BackendConfig.Aliyun.Register.Sig, config.BackendConfig.Aliyun.Register.Template, phone, warp.UserCenterWebsite)
}

func Send(ctx context.Context, TemplateParam map[string]string, sig, template string, phone string, senderID int64) errors.WTError {
	data, err := utils.JsonMarshal(TemplateParam)
	if err != nil {
		return errors.WarpQuick(err)
	}
	templateParamJson := string(data)

	smsMessageModel := db.NewSmsMessageModel(mysql.MySQLConn)
	smsMessage := &db.SmsMessage{
		Phone:         phone,
		Sig:           sig,
		Template:      template,
		TemplateParam: templateParamJson,
		SenderId:      senderID,
		Success:       true,
	}

	limitKey := fmt.Sprintf("sms:sendlimit:%s", phone)
	limitRes := redis.Get(ctx, limitKey)
	res, redisErr := limitRes.Result()
	if redisErr == nil && len(res) != 0 {
		smsMessage.Success = false
		smsMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: fmt.Sprintf("系统限制：%s", limitRes),
		}
		_, _ = smsMessageModel.Insert(ctx, smsMessage)
		return SMSSendLimit.New()
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(sig),
		TemplateCode:  tea.String(template),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(templateParamJson),
	}

	err = (func() errors.WTError {
		defer utils.Recover(logger.Logger, &err, "sms error")

		resp, err := SMSClient.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
		var sdkError *tea.SDKError
		if errors.As(err, &sdkError) {
			errorCode := *sdkError.Code
			if errorCode == "DayLimitControl" || errorCode == "isv.DAY_LIMIT_CONTROL" || errorCode == "MonthLimitControl" {
				_ = redis.Set(ctx, limitKey, fmt.Sprintf("%s: %s", errorCode, *sdkError.Message), time.Hour*5)
				smsMessage.Success = false
				smsMessage.ErrorMsg = sql.NullString{
					Valid:  true,
					String: fmt.Sprintf("%s: %s", errorCode, *sdkError.Message),
				}
				_, _ = smsMessageModel.Insert(ctx, smsMessage)
				return SMSSendLimit.New()
			} else if errorCode == "isv.BUSINESS_LIMIT_CONTROL" {
				_ = redis.Set(ctx, limitKey, fmt.Sprintf("%s: %s", errorCode, *sdkError.Message), time.Minute*5)
				smsMessage.Success = false
				smsMessage.ErrorMsg = sql.NullString{
					Valid:  true,
					String: fmt.Sprintf("%s: %s", errorCode, *sdkError.Message),
				}
				_, _ = smsMessageModel.Insert(ctx, smsMessage)
				return SMSSendLimit.New()
			} else {
				smsMessage.Success = false
				smsMessage.ErrorMsg = sql.NullString{
					Valid:  true,
					String: fmt.Sprintf("%s: %s", errorCode, *sdkError.Message),
				}
				_, _ = smsMessageModel.Insert(ctx, smsMessage)
				return SMSError.New()
			}
		} else if err != nil {
			return errors.WarpQuick(err)
		}

		if *resp.Body.Code == "DayLimitControl" || *resp.Body.Code == "isv.DAY_LIMIT_CONTROL" || *resp.Body.Code == "MonthLimitControl" || *resp.Body.Code == "isv.BUSINESS_LIMIT_CONTROL" {
			_ = redis.Set(ctx, limitKey, fmt.Sprintf("%s: %s", *resp.Body.Code, *resp.Body.Message), time.Hour*5)
			smsMessage.Success = false
			smsMessage.ErrorMsg = sql.NullString{
				Valid:  true,
				String: fmt.Sprintf("%s: %s", *resp.Body.Code, *resp.Body.Message),
			}
			_, _ = smsMessageModel.Insert(ctx, smsMessage)
			return SMSSendLimit.New()
		} else if *resp.Body.Code == "isv.BUSINESS_LIMIT_CONTROL" {
			_ = redis.Set(ctx, limitKey, fmt.Sprintf("%s: %s", *resp.Body.Code, *resp.Body.Message), time.Minute*5)
			smsMessage.Success = false
			smsMessage.ErrorMsg = sql.NullString{
				Valid:  true,
				String: fmt.Sprintf("%s: %s", *resp.Body.Code, *resp.Body.Message),
			}
			_, _ = smsMessageModel.Insert(ctx, smsMessage)
			return SMSSendLimit.New()
		} else if *resp.Body.Code != "OK" {
			smsMessage.Success = false
			smsMessage.ErrorMsg = sql.NullString{
				Valid:  true,
				String: fmt.Sprintf("%s: %s", *resp.Body.Code, *resp.Body.Message),
			}
			_, _ = smsMessageModel.Insert(ctx, smsMessage)
			return SMSError.New()
		}

		_, _ = smsMessageModel.Insert(ctx, smsMessage)
		return nil
	})()
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
