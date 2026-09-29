package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WebsiteFundingModel = (*customWebsiteFundingModel)(nil)

type (
	// WebsiteFundingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWebsiteFundingModel.
	WebsiteFundingModel interface {
		websiteFundingModel
		websiteFundingModelSelf
	}

	customWebsiteFundingModel struct {
		*defaultWebsiteFundingModel
	}
)

// NewWebsiteFundingModel returns a model for the database table.
func NewWebsiteFundingModel(conn sqlx.SqlConn) WebsiteFundingModel {
	return &customWebsiteFundingModel{
		defaultWebsiteFundingModel: newWebsiteFundingModel(conn),
	}
}

func NewWebsiteFundingModelWithSession(session sqlx.Session) WebsiteFundingModel {
	return &customWebsiteFundingModel{
		defaultWebsiteFundingModel: newWebsiteFundingModel(sqlx.NewSqlConnFromSession(session)),
	}
}
