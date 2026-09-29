package config

type WXRobotConfig struct {
	Log        string `json:"log" yaml:"log" mapstructure:"log"`
	NewUserLog string `json:"newUserLog" yaml:"newUserLog" mapstructure:"newUserLog"`
	PayLog     string `json:"payLog" yaml:"payLog" mapstructure:"payLog"`
	Sender     string `json:"sender" yaml:"sender" mapstructure:"sender"`
}
