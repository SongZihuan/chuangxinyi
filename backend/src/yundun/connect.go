package yundun

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/client"
	"github.com/alibabacloud-go/tea/tea"
	errors "github.com/wuntsong-org/wterrors"
)

var YunDunClient *green20220302.Client

func InitYunDun() errors.WTError {
	var err error
	if len(config.BackendConfig.Aliyun.Identity.AppCode) == 0 {
		return errors.Errorf("aliyun identity app code must be given")
	}

	if len(config.BackendConfig.Aliyun.Identity.AppSecret) == 0 {
		return errors.Errorf("aliyun identity app secret must be given")
	}

	if len(config.BackendConfig.Aliyun.Identity.AppKey) == 0 {
		return errors.Errorf("aliyun identity app key must be given")
	}

	if len(config.BackendConfig.Aliyun.Header.AppCode) == 0 {
		return errors.Errorf("aliyun header(nickname) app code must be given")
	}

	if len(config.BackendConfig.Aliyun.Header.AppSecret) == 0 {
		return errors.Errorf("aliyun header(nickname) app secret must be given")
	}

	if len(config.BackendConfig.Aliyun.Header.AppKey) == 0 {
		return errors.Errorf("aliyun header(nickname) app key must be given")
	}

	YunDunClient, err = green20220302.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(config.BackendConfig.Aliyun.AccessKeyId),
		AccessKeySecret: tea.String(config.BackendConfig.Aliyun.AccessKeySecret),
		Endpoint:        tea.String("green-cip.cn-shenzhen.aliyuncs.com"),
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
