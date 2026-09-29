package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WebsiteIpModel = (*customWebsiteIpModel)(nil)

type (
	// WebsiteIpModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWebsiteIpModel.
	WebsiteIpModel interface {
		websiteIpModel
		websiteIpModelSelf
	}

	customWebsiteIpModel struct {
		*defaultWebsiteIpModel
	}
)

// NewWebsiteIpModel returns a model for the database table.
func NewWebsiteIpModel(conn sqlx.SqlConn) WebsiteIpModel {
	return &customWebsiteIpModel{
		defaultWebsiteIpModel: newWebsiteIpModel(conn),
	}
}
