package config

type PasswordConfig struct {
	Salt      string `json:"salt" yaml:"salt" mapstructure:"salt"`
	FrontSalt string `json:"frontSalt" yaml:"frontSalt" mapstructure:"frontSalt"`
}
