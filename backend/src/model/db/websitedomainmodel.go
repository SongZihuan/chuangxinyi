package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WebsiteDomainModel = (*customWebsiteDomainModel)(nil)

type (
	// WebsiteDomainModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWebsiteDomainModel.
	WebsiteDomainModel interface {
		websiteDomainModel
		websiteDomainModelSelf
	}

	customWebsiteDomainModel struct {
		*defaultWebsiteDomainModel
	}
)

// NewWebsiteDomainModel returns a model for the database table.
func NewWebsiteDomainModel(conn sqlx.SqlConn) WebsiteDomainModel {
	return &customWebsiteDomainModel{
		defaultWebsiteDomainModel: newWebsiteDomainModel(conn),
	}
}
