package config

type WeChatPayConfig struct {
	AppID       string `json:"appID" yaml:"appID" mapstructure:"appID"`
	MchID       string `json:"mchID" yaml:"mchID" mapstructure:"mchID"`
	MchAPIv3Key string `json:"mchAPIv3Key" yaml:"mchAPIv3Key" mapstructure:"mchAPIv3Key"`
	PublicCert  string `json:"publicCert" yaml:"publicCert" mapstructure:"publicCert"`
	PrivateKey  string `json:"privateKey" yaml:"privateKey" mapstructure:"privateKey"`
	ReturnURL   string `json:"returnURL" yaml:"returnURL" mapstructure:"returnURL"`

	UseNativePay bool `json:"useNativePay" yaml:"useNativePay" mapstructure:"useNativePay"`
	UseH5Pay     bool `json:"useH5Pay" yaml:"useH5Pay" mapstructure:"useH5Pay"`
	UseReturnPay bool `json:"useReturnPay" yaml:"useReturnPay" mapstructure:"useReturnPay"`
	UseWithdraw  bool `json:"useWithdraw" yaml:"useWithdraw" mapstructure:"useWithdraw"`
}
