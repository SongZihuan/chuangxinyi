package afs

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	afs "github.com/alibabacloud-go/afs-20180112/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	errors "github.com/wuntsong-org/wterrors"
)

var AFSConfig = new(rpc.Config)
var AFSClient *afs.Client

func InitAFS() errors.WTError {
	if len(config.BackendConfig.Aliyun.AccessKeyId) == 0 {
		return errors.Errorf("aliyun AccessKeyId must be given")
	}

	if len(config.BackendConfig.Aliyun.AccessKeySecret) == 0 {
		return errors.Errorf("aliyun AccessKeySecret must be given")
	}

	if len(config.BackendConfig.Aliyun.AFS.CAPTCHAAppKey) == 0 {
		return errors.Errorf("aliyun CAPTCHAAppKey must be given")
	}

	if len(config.BackendConfig.Aliyun.AFS.SilenceCAPTCHAAppKey) == 0 {
		return errors.Errorf("aliyun SilenceCAPTCHAAppKey must be given")
	}

	var err error
	AFSConfig.SetAccessKeyId(config.BackendConfig.Aliyun.AccessKeyId).
		SetAccessKeySecret(config.BackendConfig.Aliyun.AccessKeySecret).
		SetRegionId("cn-hangzhou").
		SetEndpoint("afs.aliyuncs.com")
	AFSClient, err = afs.NewClient(AFSConfig)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
