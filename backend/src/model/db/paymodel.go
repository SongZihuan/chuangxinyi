package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ PayModel = (*customPayModel)(nil)

type (
	// PayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPayModel.
	PayModel interface {
		payModel
		payModelSelf
	}

	customPayModel struct {
		*defaultPayModel
	}
)

// NewPayModel returns a model for the database table.
func NewPayModel(conn sqlx.SqlConn) PayModel {
	return &customPayModel{
		defaultPayModel: newPayModel(conn),
	}
}

func NewPayModelWithSession(session sqlx.Session) PayModel {
	return &customPayModel{
		defaultPayModel: newPayModel(sqlx.NewSqlConnFromSession(session)),
	}
}
