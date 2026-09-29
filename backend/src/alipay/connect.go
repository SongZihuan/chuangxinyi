package alipay

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"github.com/SuperH-0630/gopay/alipay"
	errors "github.com/wuntsong-org/wterrors"
	"os"
)

var AlipayClient *alipay.Client
var AlipayPublicCert []byte
var AlipayRootCert []byte
var PublicCert []byte

func InitAlipay() errors.WTError {
	if len(config.BackendConfig.Alipay.AppID) == 0 {
		return errors.Errorf("alipay app id key must be given")
	}

	if len(config.BackendConfig.Alipay.PrivateKey) == 0 {
		return errors.Errorf("alipay private key must be given")
	}

	if len(config.BackendConfig.Alipay.PublicCert) == 0 {
		return errors.Errorf("alipay app public cert must be given")
	}

	if len(config.BackendConfig.Alipay.AlipayPublicCert) == 0 {
		return errors.Errorf("alipay cert key must be given")
	}

	if len(config.BackendConfig.Alipay.AlipayRootCert) == 0 {
		return errors.Errorf("alipay root cert key must be given")
	}

	if len(config.BackendConfig.Alipay.NotifyUrl) == 0 {
		return errors.Errorf("alipay return url must be given")
	}

	PublicCert, err := os.ReadFile(config.BackendConfig.Alipay.PublicCert)
	if err != nil {
		return errors.WarpQuick(err)
	}

	AlipayRootCert, err = os.ReadFile(config.BackendConfig.Alipay.AlipayRootCert)
	if err != nil {
		return errors.WarpQuick(err)
	}

	AlipayPublicCert, err = os.ReadFile(config.BackendConfig.Alipay.AlipayPublicCert)
	if err != nil {
		return errors.WarpQuick(err)
	}

	AlipayClient, err = alipay.NewClient(config.BackendConfig.Alipay.AppID, config.BackendConfig.Alipay.PrivateKey, !config.BackendConfig.Alipay.Sandbox)
	if err != nil {
		return errors.WarpQuick(err)
	}

	AlipayClient.
		SetLocation(alipay.LocationShanghai).                // 设置时区，不设置或出错均为默认服务器时间
		SetCharset(alipay.UTF8).                             // 设置字符编码，不设置默认 utf-8
		SetSignType(alipay.RSA2).                            // 设置签名类型，不设置默认 RSA2
		SetNotifyUrl(config.BackendConfig.Alipay.NotifyUrl). // 设置异步通知URL
		SetReturnUrl(config.BackendConfig.Alipay.ReturnUrl).
		AutoVerifySign(AlipayPublicCert)

	err = AlipayClient.SetCertSnByContent(PublicCert, AlipayRootCert, AlipayPublicCert)
	if err != nil {
		return errors.WarpQuick(err)
	}

	//
	//if len(config.BackendConfig.Alipay.EncryptKey) != 0 {
	//	err = AlipayClient.SetEncryptKey(config.BackendConfig.Alipay.EncryptKey)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}
