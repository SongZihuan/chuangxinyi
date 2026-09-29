package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WebsiteUrlPathModel = (*customWebsiteUrlPathModel)(nil)

type (
	// WebsiteUrlPathModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWebsiteUrlPathModel.
	WebsiteUrlPathModel interface {
		websiteUrlPathModel
		websiteUrlPathModelSelf
	}

	customWebsiteUrlPathModel struct {
		*defaultWebsiteUrlPathModel
	}
)

// NewWebsiteUrlPathModel returns a model for the database table.
func NewWebsiteUrlPathModel(conn sqlx.SqlConn) WebsiteUrlPathModel {
	return &customWebsiteUrlPathModel{
		defaultWebsiteUrlPathModel: newWebsiteUrlPathModel(conn),
	}
}
