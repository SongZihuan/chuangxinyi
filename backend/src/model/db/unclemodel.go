package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ UncleModel = (*customUncleModel)(nil)

type (
	// UncleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUncleModel.
	UncleModel interface {
		uncleModel
		uncleModelSelf
	}

	customUncleModel struct {
		*defaultUncleModel
	}
)

// NewUncleModel returns a model for the database table.
func NewUncleModel(conn sqlx.SqlConn) UncleModel {
	return &customUncleModel{
		defaultUncleModel: newUncleModel(conn),
	}
}
