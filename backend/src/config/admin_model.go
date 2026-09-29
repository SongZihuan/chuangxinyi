package config

type Role struct {
	RoleName             string `json:"roleName" yaml:"roleName" mapstructure:"roleName"`
	RoleSign             string `json:"roleSign" yaml:"roleSign" mapstructure:"roleSign"`
	RoleDescribe         string `json:"roleDescribe" yaml:"roleDescribe" mapstructure:"roleDescribe"`
	NotDelete            bool   `json:"notDelete" yaml:"notDelete" mapstructure:"notDelete"`
	NotChangeSign        bool   `json:"notChangeSign" yaml:"notChangeSign" yaml:"notChangeSign"`
	NotChangePermissions bool   `json:"notChangePermissions" yaml:"notChangePermissions" yaml:"notChangePermissions"`
	ResetPermission      bool   `json:"resetPermission"  yaml:"resetPermission" yaml:"resetPermission"`
}

type AdminConfig struct {
	AdminPhone    string `json:"adminPhone" yaml:"adminPhone" mapstructure:"adminPhone"`
	RootRole      Role   `json:"rootRole" yaml:"rootRole" mapstructure:"rootRole"`
	UserRole      Role   `json:"userRole" yaml:"userRole" mapstructure:"userRole"`
	AnonymousRole Role   `json:"anonymousRole" yaml:"anonymousRole" mapstructure:"anonymousRole"`

	MenuDepth int64 `json:"menuDepth" yaml:"menuDepth" mapstructure:"menuDepth"`

	ICP1      string `json:"icp1" yaml:"icp1" mapstructure:"icp1"`
	ICP2      string `json:"icp2" yaml:"icp2" mapstructure:"icp2"`
	Gongan    string `json:"gongan" yaml:"icp1" mapstructure:"gongan"`
	Copyright string `json:"copyright" yaml:"copyright" mapstructure:"copyright"`
}
