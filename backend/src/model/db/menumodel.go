package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ MenuModel = (*customMenuModel)(nil)

type (
	// MenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMenuModel.
	MenuModel interface {
		menuModel
		menuModelSelf
	}

	customMenuModel struct {
		*defaultMenuModel
	}
)

// NewMenuModel returns a model for the database table.
func NewMenuModel(conn sqlx.SqlConn) MenuModel {
	return &customMenuModel{
		defaultMenuModel: newMenuModel(conn),
	}
}

func NewMenuModelWithSession(session sqlx.Session) MenuModel {
	return &customMenuModel{
		defaultMenuModel: newMenuModel(sqlx.NewSqlConnFromSession(session)),
	}
}
