package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ EmailMessageModel = (*customEmailMessageModel)(nil)

type (
	// EmailMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailMessageModel.
	EmailMessageModel interface {
		emailMessageModel
		emailMessageModelSelf
	}

	customEmailMessageModel struct {
		*defaultEmailMessageModel
	}
)

// NewEmailMessageModel returns a model for the database table.
func NewEmailMessageModel(conn sqlx.SqlConn) EmailMessageModel {
	return &customEmailMessageModel{
		defaultEmailMessageModel: newEmailMessageModel(conn),
	}
}
