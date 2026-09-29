package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ ApplicationModel = (*customApplicationModel)(nil)

type (
	// ApplicationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApplicationModel.
	ApplicationModel interface {
		applicationModel
		applicationModelSelf
	}

	customApplicationModel struct {
		*defaultApplicationModel
	}
)

// NewApplicationModel returns a model for the database table.
func NewApplicationModel(conn sqlx.SqlConn) ApplicationModel {
	return &customApplicationModel{
		defaultApplicationModel: newApplicationModel(conn),
	}
}

func NewApplicationModelWithSession(session sqlx.Session) ApplicationModel {
	return &customApplicationModel{
		defaultApplicationModel: newApplicationModel(sqlx.NewSqlConnFromSession(session)),
	}
}
