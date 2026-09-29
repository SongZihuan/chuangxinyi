package config

type CoinConfig struct {
	ID             string `json:"id" yaml:"id" mapstructure:"id"`
	Name           string `json:"name" yaml:"name" mapstructure:"name"`
	Price          int64  `json:"price" yaml:"price" mapstructure:"price"`
	ShowUrl        string `json:"showUrl" yaml:"showUrl" mapstructure:"showUrl"`
	TimeExpireSec  int64  `json:"timeExpireSec" yaml:"timeExpireSec" mapstructure:"timeExpireSec"`
	RefundDayLimit int64  `json:"refundDayLimit" yaml:"refundDayLimit" mapstructure:"refundDayLimit"`
	RefundReason   string `json:"refundReason" yaml:"refundReason" mapstructure:"refundReason"`
	WithdrawMin    int64  `json:"withdrawMin" yaml:"withdrawMin" mapstructure:"withdrawMin"`
}
