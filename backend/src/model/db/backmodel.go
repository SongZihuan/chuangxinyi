package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ BackModel = (*customBackModel)(nil)

type (
	// BackModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBackModel.
	BackModel interface {
		backModel
		backModelSelf
	}

	customBackModel struct {
		*defaultBackModel
	}
)

// NewBackModel returns a model for the database table.
func NewBackModel(conn sqlx.SqlConn) BackModel {
	return &customBackModel{
		defaultBackModel: newBackModel(conn),
	}
}

func NewBackModelWithSession(session sqlx.Session) BackModel {
	return &customBackModel{
		defaultBackModel: newBackModel(sqlx.NewSqlConnFromSession(session)),
	}
}
