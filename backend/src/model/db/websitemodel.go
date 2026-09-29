package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WebsiteModel = (*customWebsiteModel)(nil)

type (
	// WebsiteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWebsiteModel.
	WebsiteModel interface {
		websiteModel
		websiteModelSelf
	}

	customWebsiteModel struct {
		*defaultWebsiteModel
	}
)

// NewWebsiteModel returns a model for the database table.
func NewWebsiteModel(conn sqlx.SqlConn) WebsiteModel {
	return &customWebsiteModel{
		defaultWebsiteModel: newWebsiteModel(conn),
	}
}
