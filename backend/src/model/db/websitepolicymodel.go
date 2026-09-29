package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WebsitePolicyModel = (*customWebsitePolicyModel)(nil)

type (
	// WebsitePolicyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWebsitePolicyModel.
	WebsitePolicyModel interface {
		websitePolicyModel
		websitePolicyModelSelf
	}

	customWebsitePolicyModel struct {
		*defaultWebsitePolicyModel
	}
)

// NewWebsitePolicyModel returns a model for the database table.
func NewWebsitePolicyModel(conn sqlx.SqlConn) WebsitePolicyModel {
	return &customWebsitePolicyModel{
		defaultWebsitePolicyModel: newWebsitePolicyModel(conn),
	}
}

func NewWebsitePolicyModelWithSession(session sqlx.Session) WebsitePolicyModel {
	return &customWebsitePolicyModel{
		defaultWebsitePolicyModel: newWebsitePolicyModel(sqlx.NewSqlConnFromSession(session)),
	}
}
