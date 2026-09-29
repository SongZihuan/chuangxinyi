package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ SecondfaModel = (*customSecondfaModel)(nil)

type (
	// SecondfaModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSecondfaModel.
	SecondfaModel interface {
		secondfaModel
		secondfaModelSelf
	}

	customSecondfaModel struct {
		*defaultSecondfaModel
	}
)

// NewSecondfaModel returns a model for the database table.
func NewSecondfaModel(conn sqlx.SqlConn) SecondfaModel {
	return &customSecondfaModel{
		defaultSecondfaModel: newSecondfaModel(conn),
	}
}
