package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ EmailModel = (*customEmailModel)(nil)

type (
	// EmailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailModel.
	EmailModel interface {
		emailModel
		emailModelSelf
	}

	customEmailModel struct {
		*defaultEmailModel
	}
)

// NewEmailModel returns a model for the database table.
func NewEmailModel(conn sqlx.SqlConn) EmailModel {
	return &customEmailModel{
		defaultEmailModel: newEmailModel(conn),
	}
}

func NewEmailModelWithSession(session sqlx.Session) EmailModel {
	return &customEmailModel{
		defaultEmailModel: newEmailModel(sqlx.NewSqlConnFromSession(session)),
	}
}
