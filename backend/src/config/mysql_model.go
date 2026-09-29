package config

type ClearInfo struct {
	TableName string `json:"tableName" yaml:"tableName" mapstructure:"tableName"`
	SaveDay   int64  `json:"saveDay" yaml:"saveDay" mapstructure:"saveDay"`
}

type MySQLConfig struct {
	DSN         string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
	SQLFilePath string `json:"sqlFilePath" yaml:"sqlFilePath" mapstructure:"sqlFilePath"`

	SystemResourceLimit int64 `json:"systemResourceLimit" yaml:"systemResourceLimit" mapstructure:"systemResourceLimit"`
	WorkOrderFileLimit  int64 `json:"workOrderFileLimit" yaml:"workOrderFileLimit" mapstructure:"workOrderFileLimit"`
	SameWalletUserLimit int64 `json:"sameWalletUserLimit" yaml:"sameWalletUserLimit" mapstructure:"sameWalletUserLimit"`
	SonUserLimit        int64 `json:"sonUserLimit" yaml:"sonUserLimit" mapstructure:"sonUserLimit"`
	UncleUserLimit      int64 `json:"uncleUserLimit" yaml:"uncleUserLimit" mapstructure:"uncleUserLimit"`
	NephewLimit         int64 `json:"nephewLimit" yaml:"nephewLimit" mapstructure:"nephewLimit"`
	SonLevelLimit       int64 `json:"sonLevelLimit" yaml:"sonLevelLimit" mapstructure:"sonLevelLimit"`

	ClearDeleteAt []ClearInfo `json:"clearDeleteAt" yaml:"clearDeleteAt" mapstructure:"clearDeleteAt"`
	ClearCreateAt []ClearInfo `json:"clearCreateAt" yaml:"clearCreateAt" mapstructure:"clearCreateAt"`
}
