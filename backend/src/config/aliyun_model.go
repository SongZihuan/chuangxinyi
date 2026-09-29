package config

type SMSConfig struct {
	Sig      string `json:"sig" yaml:"sig" mapstructure:"sig"`
	Template string `json:"template" yaml:"template" mapstructure:"template"`
}

type YunDunConfig struct {
	AppKey    string `json:"appKey" yaml:"appKey" mapstructure:"appKey"`
	AppSecret string `json:"appSecret" yaml:"appSecret" mapstructure:"appSecret"`
	AppCode   string `json:"appCode" yaml:"appCode" mapstructure:"appCode"`
}

type OSSConfig struct {
	Endpoint   string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	BucketName string `json:"bucketName" yaml:"bucketName" mapstructure:"bucketName"`
}

type ImportCodeConfig struct {
	SMSConfig `yaml:",inline" mapstructure:",squash"`
}

type CodeConfig struct {
	SMSConfig `yaml:",inline" mapstructure:",squash"`
}

type ChangeConfig struct {
	SMSConfig `yaml:",inline" mapstructure:",squash"`
}

type ChangePhoneConfig struct {
	SMSConfig `yaml:",inline" mapstructure:",squash"`
}

type DeleteConfig struct {
	SMSConfig `yaml:",inline" mapstructure:",squash"`
}

type RegisterConfig struct {
	SMSConfig `yaml:",inline" mapstructure:",squash"`
}

type IdentityConfig struct {
	YunDunConfig `yaml:",inline" mapstructure:",squash"`
	OSSConfig    `yaml:",inline" mapstructure:",squash"`
	Sign         OSSConfig `json:"sign" yaml:"sign" mapstructure:"sign"`
}

type HeaderNicknameConfig struct {
	YunDunConfig `yaml:",inline" mapstructure:",squash"`
	OSSConfig    `yaml:",inline" mapstructure:",squash"`
	ImageStyle   string    `json:"imageStyle" yaml:"imageStyle" mapstructure:"imageStyle"`
	Sign         OSSConfig `json:"sign" yaml:"sign" mapstructure:"sign"`
}

type AFSConfig struct {
	CAPTCHAStatus        bool   `json:"captchaStatus" yaml:"captchaStatus" mapstructure:"captchaStatus"`
	CAPTCHAAppKey        string `json:"captchaAppKey" yaml:"captchaAppKey" mapstructure:"captchaAppKey"`
	CAPTCHAScene         string `json:"captchaScene" yaml:"captchaScene" mapstructure:"captchaScene"`
	SilenceCAPTCHAStatus bool   `json:"silenceCAPTCHAStatus" yaml:"silenceCAPTCHAStatus" mapstructure:"silenceCAPTCHAStatus"`
	SilenceCAPTCHAAppKey string `json:"silenceCAPTCHAAppKey" yaml:"silenceCAPTCHAAppKey" mapstructure:"silenceCAPTCHAAppKey"`
	SilenceCAPTCHAScene  string `json:"silenceCAPTCHAScene" yaml:"silenceCAPTCHAScene" mapstructure:"silenceCAPTCHAScene"`
}

type IPConfig struct {
	AppKey        string `json:"appKey" yaml:"appKey" mapstructure:"appKey"`
	AppCode       string `json:"appCode" yaml:"appCode" mapstructure:"appCode"`
	AppSecret     string `json:"appSecret" yaml:"appSecret" mapstructure:"appSecret"`
	ExpiresSecond int64  `json:"expiresSecond" yaml:"expiresSecond" mapstructure:"expiresSecond"`
}

type FileConfig struct {
	OSSConfig `yaml:",inline" mapstructure:",squash"`
	Sign      OSSConfig `json:"sign" yaml:"sign" mapstructure:"sign"`
}

type WorkOrderFileConfig struct {
	OSSConfig `yaml:",inline" mapstructure:",squash"`
	Sign      OSSConfig `json:"sign" yaml:"sign" mapstructure:"sign"`
}

type InvoiceFileConfig struct {
	OSSConfig `yaml:",inline" mapstructure:",squash"`
	Sign      OSSConfig `json:"sign" yaml:"sign" mapstructure:"sign"`
}

type AliyunConfig struct {
	AccessKeyId     string `json:"accessKeyId" yaml:"accessKeyId" mapstructure:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret" mapstructure:"accessKeySecret"`

	Identity    IdentityConfig       `json:"identity" yaml:"identity" mapstructure:"identity"`
	OcrEndpoint string               `json:"ocrEndpoint" yaml:"ocrEndpoint" mapstructure:"ocrEndpoint"`
	ImportCode  ImportCodeConfig     `json:"smsImportCode" yaml:"smsImportCode" mapstructure:"smsImportCode"`
	Code        CodeConfig           `json:"smsCode" yaml:"smsCode" mapstructure:"smsCode"`
	Change      ChangeConfig         `json:"smsChange" yaml:"smsChange" mapstructure:"smsChange"`
	ChangePhone ChangePhoneConfig    `json:"smsChangePhone" yaml:"smsChangePhone" mapstructure:"smsChangePhone"`
	Delete      DeleteConfig         `json:"smsDelete" yaml:"smsDelete" mapstructure:"smsDelete"`
	Register    RegisterConfig       `json:"smsRegister" yaml:"smsRegister" mapstructure:"smsRegister"`
	AFS         AFSConfig            `json:"afs" yaml:"afs" mapstructure:"afs"`
	IP          IPConfig             `json:"ip" yaml:"ip" mapstructure:"ip"`
	Header      HeaderNicknameConfig `json:"header" yaml:"header" mapstructure:"header"`
	File        FileConfig           `json:"file" yaml:"file" mapstructure:"file"`
	WorkOrder   WorkOrderFileConfig  `json:"workOrder" yaml:"workOrder" mapstructure:"workOrder"`
	Invoice     InvoiceFileConfig    `json:"invoice" yaml:"invoice" mapstructure:"invoice"`
}
