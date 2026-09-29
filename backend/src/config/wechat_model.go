package config

type WeChatConfig struct {
	AppID     string `json:"appID" yaml:"appID" mapstructure:"appID"`
	AppSecret string `json:"appSecret" yaml:"appSecret" mapstructure:"appSecret"`
}
