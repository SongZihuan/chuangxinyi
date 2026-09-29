package config

type AlipayConfig struct {
	Sandbox          bool   `json:"sandbox" yaml:"sandbox" mapstructure:"sandbox"`
	AppID            string `json:"appID" yaml:"appID" mapstructure:"appID"`
	PrivateKey       string `json:"privateKey" yaml:"privateKey" mapstructure:"privateKey"`
	PublicCert       string `json:"publicCert" yaml:"publicCert" mapstructure:"publicCert"`
	AlipayPublicCert string `json:"AlipayPublicCert" yaml:"AlipayPublicCert" mapstructure:"AlipayPublicCert"`
	AlipayRootCert   string `json:"alipayRootCert" yaml:"alipayRootCert" mapstructure:"alipayRootCert"`
	EncryptKey       string `json:"encryptKey" yaml:"encryptKey" mapstructure:"encryptKey"`
	NotifyUrl        string `json:"notifyUrl" yaml:"notifyUrl" mapstructure:"notifyUrl"`
	ReturnUrl        string `json:"returnUrl" yaml:"returnUrl" mapstructure:"returnUrl"`
	FaceReturnUrl    string `json:"faceReturnUrl" yaml:"faceReturnUrl" mapstructure:"faceReturnUrl"`
	WapQuitUrl       string `json:"wapQuitUrl" yaml:"wapQuitUrl" mapstructure:"wapQuitUrl"`

	UsePCPay     bool `json:"usePCPay" yaml:"usePCPay" mapstructure:"usePCPay"`
	UseWapPay    bool `json:"useWapPay" yaml:"useWapPay" mapstructure:"useWapPay"`
	UseFaceCheck bool `json:"useFaceCheck" yaml:"useFaceCheck" mapstructure:"useFaceCheck"`
	UseReturnPay bool `json:"useReturnPay" yaml:"useReturnPay" mapstructure:"useReturnPay"`
	UseWithdraw  bool `json:"useWithdraw" yaml:"useWithdraw" mapstructure:"useWithdraw"`
}
