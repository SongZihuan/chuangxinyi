package wechatpay

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/SuperH-0630/gopay/wechat/v3"
	errors "github.com/wuntsong-org/wterrors"
	"os"
)

var WeChatPayClient *wechat.ClientV3
var MchPrivateKey []byte

func InitWeChatPay() errors.WTError {
	if len(config.BackendConfig.WeChatPay.AppID) == 0 {
		return errors.Errorf("wechat pay app id must be given")
	}

	if len(config.BackendConfig.WeChatPay.MchID) == 0 {
		return errors.Errorf("wechat pay mch id id must be given")
	}

	if len(config.BackendConfig.WeChatPay.MchAPIv3Key) == 0 {
		return errors.Errorf("wechat pay api v3 key id must be given")
	}

	if len(config.BackendConfig.WeChatPay.PublicCert) == 0 {
		return errors.Errorf("wechat pay public cert id must be given")
	}

	if len(config.BackendConfig.WeChatPay.PrivateKey) == 0 {
		return errors.Errorf("wechat pay privare key id must be given")
	}

	if len(config.BackendConfig.WeChatPay.ReturnURL) == 0 {
		return errors.Errorf("wechat pay return url id must be given")
	}

	MchPrivateKey, err := os.ReadFile(config.BackendConfig.WeChatPay.PrivateKey)
	if err != nil {
		return errors.WarpQuick(err)
	}

	mchCert, err := utils.LoadCertificateWithPath(config.BackendConfig.WeChatPay.PublicCert)
	if err != nil {
		return errors.WarpQuick(err)
	}

	// 使用商户私钥等初始化 client并使它具有自动定时获取微信支付平台证书的能力
	WeChatPayClient, err = wechat.NewClientV3(config.BackendConfig.WeChatPay.MchID, utils.GetCertificateSerialNumber(*mchCert), config.BackendConfig.WeChatPay.MchAPIv3Key, string(MchPrivateKey))
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = WeChatPayClient.AutoVerifySign()
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
