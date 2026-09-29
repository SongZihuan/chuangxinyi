package config

type UserConfig struct {
	GetBusySecond    int64    `json:"getBusySecond" yaml:"getBusySecond" mapstructure:"getBusySecond"`
	GetBusyLimit     int64    `json:"getBusyLimit" yaml:"getBusyLimit" mapstructure:"getBusyLimit"`
	PostBusySecond   int64    `json:"postBusySecond" yaml:"postBusySecond" mapstructure:"postBusySecond"`
	PostBusyLimit    int64    `json:"postBusyLimit" yaml:"postBusyLimit" mapstructure:"postBusyLimit"`
	Origin           []string `json:"origin" yaml:"origin" mapstructure:"origin"`
	ReadableName     string   `json:"readableName" yaml:"readableName" mapstructure:"readableName"`
	WebsiteUID       string   `json:"websiteUID" yaml:"websiteUID" mapstructure:"websiteUID"`
	AllowMethod      []string `json:"allowMethod" yaml:"allowMethod" mapstructure:"allowMethod"`
	AllowHeader      []string `json:"allowHeader" yaml:"allowHeader" mapstructure:"allowHeader"`
	Url              string   `json:"url" yaml:"url" mapstructure:"url"`
	Group            string   `json:"group" yaml:"group" mapstructure:"group"`
	GoZeroHttpConfig `yaml:",inline" mapstructure:",squash"`
}
