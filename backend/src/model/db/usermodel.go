package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		userModelSelf
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func NewUserModelWithSession(session sqlx.Session) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(sqlx.NewSqlConnFromSession(session)),
	}
}
