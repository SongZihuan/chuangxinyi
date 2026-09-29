package config

type RedisConfig struct {
	Addr     string `json:"addr" yaml:"addr" mapstructure:"addr"`
	UserName string `json:"userName" yaml:"userName" mapstructure:"userName"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
	DB       int64  `json:"db" yaml:"db" mapstructure:"db"`
}
