package config

type SmtpConfig struct {
	Addr             string `json:"addr" yaml:"addr" mapstructure:"addr"`
	UserName         string `json:"userName" yaml:"userName" mapstructure:"userName"`
	Password         string `json:"password" yaml:"password" mapstructure:"password"`
	FromEmail        string `json:"fromEmail" yaml:"fromEmail" mapstructure:"fromEmail"`
	Sender           string `json:"sender" yaml:"sender" mapstructure:"sender"`
	TemplateFilePath string `json:"templateFilePath" yaml:"templateFilePath" mapstructure:"templateFilePath"`
	Sig              string `json:"sig" yaml:"sig" mapstructure:"sig"`
}
