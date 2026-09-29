package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ Oauth2RecordModel = (*customOauth2RecordModel)(nil)

type (
	// Oauth2RecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOauth2RecordModel.
	Oauth2RecordModel interface {
		oauth2RecordModel
		oauth2RecordModelSelf
	}

	customOauth2RecordModel struct {
		*defaultOauth2RecordModel
	}
)

// NewOauth2RecordModel returns a model for the database table.
func NewOauth2RecordModel(conn sqlx.SqlConn) Oauth2RecordModel {
	return &customOauth2RecordModel{
		defaultOauth2RecordModel: newOauth2RecordModel(conn),
	}
}
