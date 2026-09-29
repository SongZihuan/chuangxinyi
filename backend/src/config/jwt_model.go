package config

import "github.com/spf13/viper"

type JWTDataConfig struct {
	Subject       string `json:"subject" yaml:"subject" mapstructure:"subject"`
	Issuer        string `json:"issuer" yaml:"issuer" mapstructure:"issuer"`
	ExpiresSecond int64  `json:"expiresSecond" yaml:"expiresSecond" mapstructure:"expiresSecond"`
}

type JWTConfig struct {
	Phone         JWTDataConfig `json:"phone" yaml:"phone" mapstructure:"phone"`
	Email         JWTDataConfig `json:"email" yaml:"email" mapstructure:"email"`
	Login         JWTDataConfig `json:"login" yaml:"login" mapstructure:"login"`
	SecondFA      JWTDataConfig `json:"2fa" yaml:"2fa" mapstructure:"2fa"`
	PassSecondFA  JWTDataConfig `json:"pass2fa" yaml:"pass2fa" mapstructure:"pass2fa"`
	Defray        JWTDataConfig `json:"defray" yaml:"defray" mapstructure:"defray"`
	User          JWTDataConfig `json:"user" yaml:"user" mapstructure:"user"`
	WeChat        JWTDataConfig `json:"wechat" yaml:"wechat" mapstructure:"wechat"`
	IDCard        JWTDataConfig `json:"idcard" yaml:"idcard" mapstructure:"idcard"`
	Company       JWTDataConfig `json:"company" yaml:"company" mapstructure:"company"`
	Face          JWTDataConfig `json:"face" yaml:"face" mapstructure:"face"`
	Delete        JWTDataConfig `json:"delete" yaml:"delete" mapstructure:"delete"`
	CheckSecondFA JWTDataConfig `json:"check2fa" yaml:"check2fa" mapstructure:"check2fa"`
	ExpireSecond  int64         `json:"expireSecond" yaml:"expireSecond" mapstructure:"expireSecond"`
}

func JwtSetDefaultValue(v *viper.Viper) {
	v.SetDefault("jwt.phone.subject", "PhoneCheck")
	v.SetDefault("jwt.email.subject", "EmailCheck")
	v.SetDefault("jwt.login.subject", "Login")
	v.SetDefault("jwt.2fa.subject", "Login2FA")
	v.SetDefault("jwt.user.subject", "WunTsongUser")
	v.SetDefault("jwt.admin2fa.subject", "Admin2FA")
	v.SetDefault("jwt.admin.subject", "WunTsongAdmin")
	v.SetDefault("jwt.wechat.subject", "WeChat")
	v.SetDefault("jwt.idcard.subject", "IDCard")
	v.SetDefault("jwt.company.subject", "Company")
	v.SetDefault("jwt.face.subject", "Face")
	v.SetDefault("jwt.phone.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.email.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.login.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.2fa.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.user.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.admin2fa.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.admin.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.wechat.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.idcard.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.company.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
	v.SetDefault("jwt.face.issuer", "GuangZhou WunTsong Information Technology Co,. Ltd")
}
