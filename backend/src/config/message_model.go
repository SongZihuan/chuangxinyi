package config

type MessageConfig struct {
	Sender     string `json:"sender" yaml:"sender" mapstructure:"sender"`
	SenderLink string `json:"senderLink" yaml:"senderLink" mapstructure:"senderLink"`
}
