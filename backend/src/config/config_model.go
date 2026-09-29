package config

type RunMode int

const (
	RunModeDevelop RunMode = iota
	RunModeRelease RunMode = iota
)

type Config struct {
	Mode string `json:"mode" yaml:"mode" mapstructure:"mode"`

	User      UserConfig      `json:"user" yaml:"user" mapstructure:"user"`
	MySQL     MySQLConfig     `json:"mysql" yaml:"mysql" mapstructure:"mysql"`
	Redis     RedisConfig     `json:"redis" yaml:"redis" mapstructure:"redis"`
	Cache     RedisConfig     `json:"cache" yaml:"cache" mapstructure:"cache"`
	Aliyun    AliyunConfig    `json:"aliyun" yaml:"aliyun" mapstructure:"aliyun"`
	Smtp      SmtpConfig      `json:"smtp" yaml:"smtp" mapstructure:"smtp"`
	Jwt       JWTConfig       `json:"jwt" yaml:"jwt" mapstructure:"jwt"`
	Cron      CronConfig      `json:"cron" yaml:"cron" mapstructure:"cron"`
	Admin     AdminConfig     `json:"admin" yaml:"admin" mapstructure:"admin"`
	WeChat    WeChatConfig    `json:"wechat" yaml:"wechat" mapstructure:"wechat"`
	Totp      TotpConfig      `json:"totp" yaml:"totp" mapstructure:"totp"`
	Alipay    AlipayConfig    `json:"alipay" yaml:"alipay" mapstructure:"alipay"`
	WeChatPay WeChatPayConfig `json:"wechatpay" yaml:"wechatpay" mapstructure:"wechatpay"`
	Coin      CoinConfig      `json:"coin" yaml:"coin" mapstructure:"coin"`
	Defray    DefrayConfig    `json:"defray" yaml:"defray" mapstructure:"defray"`
	Sign      SignConfig      `json:"sign" yaml:"sign" mapstructure:"sign"`
	Message   MessageConfig   `json:"message" yaml:"message" mapstructure:"message"`
	WXRobot   WXRobotConfig   `json:"wxrobot" yaml:"wxrobot" mapstructure:"wxrobot"`
	FuWuHao   FuWuHaoConfig   `json:"fuWuHao" yaml:"fuWuHao" mapstructure:"fuWuHao"`
	Password  PasswordConfig  `json:"password" yaml:"password" mapstructure:"password"`
	SqlClear  SqlClearConfig  `json:"sqlClear" yaml:"sqlClear" mapstructure:"sqlClear"`
}
