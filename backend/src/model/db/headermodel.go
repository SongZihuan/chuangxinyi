package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ HeaderModel = (*customHeaderModel)(nil)

type (
	// HeaderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHeaderModel.
	HeaderModel interface {
		headerModel
		headerModelSelf
	}

	customHeaderModel struct {
		*defaultHeaderModel
	}
)

// NewHeaderModel returns a model for the database table.
func NewHeaderModel(conn sqlx.SqlConn) HeaderModel {
	return &customHeaderModel{
		defaultHeaderModel: newHeaderModel(conn),
	}
}
