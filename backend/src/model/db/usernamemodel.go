package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ UsernameModel = (*customUsernameModel)(nil)

type (
	// UsernameModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsernameModel.
	UsernameModel interface {
		usernameModel
		usernameModelSelf
	}

	customUsernameModel struct {
		*defaultUsernameModel
	}
)

// NewUsernameModel returns a model for the database table.
func NewUsernameModel(conn sqlx.SqlConn) UsernameModel {
	return &customUsernameModel{
		defaultUsernameModel: newUsernameModel(conn),
	}
}

func NewUsernameModelWithSession(session sqlx.Session) UsernameModel {
	return &customUsernameModel{
		defaultUsernameModel: newUsernameModel(sqlx.NewSqlConnFromSession(session)),
	}
}
