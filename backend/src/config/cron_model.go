package config

type CronConfig struct {
	PermissionUpdate        string `json:"permissionUpdate" yaml:"permissionUpdate" mapstructure:"permissionUpdate"`
	UrlPathUpdate           string `json:"urlPathUpdate" yaml:"urlPathUpdate" mapstructure:"urlPathUpdate"`
	RoleUpdate              string `json:"roleUpdate" yaml:"roleUpdate" mapstructure:"roleUpdate"`
	MenuUpdate              string `json:"menuUpdate" yaml:"menuUpdate" mapstructure:"menuUpdate"`
	WebsitePermissionUpdate string `json:"websitePermissionUpdate" yaml:"websitePermissionUpdate" mapstructure:"websitePermissionUpdate"`
	WebsiteUrlPathUpdate    string `json:"websiteUrlPathUpdate" yaml:"websiteUrlPathUpdate" mapstructure:"websiteUrlPathUpdate"`
	WebsiteUpdate           string `json:"websiteUpdate" yaml:"websiteUpdate" mapstructure:"websiteUpdate"`
	ApplicationUpdate       string `json:"applicationUpdate" yaml:"applicationUpdate" mapstructure:"applicationUpdate"`
}
