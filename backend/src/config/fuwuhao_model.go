package config

type TemplateMsgConfig struct {
	TemplateID string `json:"templateID"`
	Url        string `json:"url"`
}

type KefuConfig struct {
	HuanChuang string `json:"huanChuang" yaml:"huanChuang" mapstructure:"huanChuang"`
	Vxwk       string `json:"vxwk" yaml:"vxwk" mapstructure:"vxwk"`
}

type MenuConfig struct {
	UpdateMenu     bool   `json:"updateMenu" yaml:"updateMenu" mapstructure:"updateMenu"`
	AboutUsWebsite string `json:"aboutUsWebsite" yaml:"aboutUsWebsite" mapstructure:"aboutUsWebsite"`
	AboutUsContact string `json:"aboutUsContact" yaml:"aboutUsContact" mapstructure:"aboutUsContact"`
	AboutUsKefu    string `json:"aboutUsKefu" yaml:"aboutUsKefu" mapstructure:"aboutUsKefu"`

	ProductVxwk      string `json:"productVXWK" yaml:"productVXWK" mapstructure:"productVXWK"`
	ProductWallpaper string `json:"ProductWallpaper" yaml:"ProductWallpaper" mapstructure:"ProductWallpaper"`

	Auth string `json:"auth" yaml:"auth" mapstructure:"auth"`
}

type FuWuHaoConfig struct {
	AppID          string `json:"appID" yaml:"appID" mapstructure:"appID"`
	Secret         string `json:"secret" yaml:"secret" mapstructure:"secret"`
	Token          string `json:"token" yaml:"token" mapstructure:"token"`
	EncodingAESKey string `json:"encodingAESKey" yaml:"encodingAESKey" mapstructure:"encodingAESKey"`

	Register     TemplateMsgConfig `json:"register" yaml:"register" mapstructure:"register"`
	Project      TemplateMsgConfig `json:"project" yaml:"project" mapstructure:"project"`
	Pay          TemplateMsgConfig `json:"pay" yaml:"pay" mapstructure:"pay"`
	Oauth2       TemplateMsgConfig `json:"oauth2" yaml:"oauth2" mapstructure:"oauth2"`
	LoginSuccess TemplateMsgConfig `json:"loginSuccess" yaml:"loginSuccess" mapstructure:"loginSuccess"`
	LoginFail    TemplateMsgConfig `json:"loginFail" yaml:"loginFail" mapstructure:"loginFail"`
	UserDelete   TemplateMsgConfig `json:"userDelete" yaml:"userDelete" mapstructure:"userDelete"`

	Kefu KefuConfig `json:"kefu" yaml:"kefu" mapstructure:"kefu"`
	Menu MenuConfig `json:"menu" yaml:"menu" mapstructure:"menu"`
}
