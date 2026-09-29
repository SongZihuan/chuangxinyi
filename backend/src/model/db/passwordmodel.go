package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ PasswordModel = (*customPasswordModel)(nil)

type (
	// PasswordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPasswordModel.
	PasswordModel interface {
		passwordModel
		passwordModelSelf
	}

	customPasswordModel struct {
		*defaultPasswordModel
	}
)

// NewPasswordModel returns a model for the database table.
func NewPasswordModel(conn sqlx.SqlConn) PasswordModel {
	return &customPasswordModel{
		defaultPasswordModel: newPasswordModel(conn),
	}
}

func NewPasswordModelWithSession(session sqlx.Session) PasswordModel {
	return &customPasswordModel{
		defaultPasswordModel: newPasswordModel(sqlx.NewSqlConnFromSession(session)),
	}
}
