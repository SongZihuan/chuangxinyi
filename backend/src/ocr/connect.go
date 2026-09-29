package ocr

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ocr_api "github.com/alibabacloud-go/ocr-api-20210707/client"
	"github.com/alibabacloud-go/tea/tea"
	errors "github.com/wuntsong-org/wterrors"
)

var OcrClient *ocr_api.Client

func InitOcr() errors.WTError {
	if len(config.BackendConfig.Aliyun.AccessKeyId) == 0 {
		return errors.Errorf("aliyun AccessKeyId must be given")
	}

	if len(config.BackendConfig.Aliyun.AccessKeySecret) == 0 {
		return errors.Errorf("aliyun AccessKeySecret must be given")
	}

	if len(config.BackendConfig.Aliyun.OcrEndpoint) == 0 {
		return errors.Errorf("aliyun ocr endpoint must be given")
	}

	conf := &openapi.Config{
		AccessKeyId:     tea.String(config.BackendConfig.Aliyun.AccessKeyId),
		AccessKeySecret: tea.String(config.BackendConfig.Aliyun.AccessKeySecret),
		Endpoint:        tea.String(config.BackendConfig.Aliyun.OcrEndpoint),
	}

	var err error
	OcrClient, err = ocr_api.NewClient(conf)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
