package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ Oauth2BanedModel = (*customOauth2BanedModel)(nil)

type (
	// Oauth2BanedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOauth2BanedModel.
	Oauth2BanedModel interface {
		oauth2BanedModel
		oauth2BanedModelSelf
	}

	customOauth2BanedModel struct {
		*defaultOauth2BanedModel
	}
)

// NewOauth2BanedModel returns a model for the database table.
func NewOauth2BanedModel(conn sqlx.SqlConn) Oauth2BanedModel {
	return &customOauth2BanedModel{
		defaultOauth2BanedModel: newOauth2BanedModel(conn),
	}
}
