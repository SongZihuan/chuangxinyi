package config

type TotpConfig struct {
	IssuerName string `json:"issuerName" yaml:"issuerName" mapstructure:"issuerName"`
}
