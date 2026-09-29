package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ TokenRecordModel = (*customTokenRecordModel)(nil)

type (
	// TokenRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTokenRecordModel.
	TokenRecordModel interface {
		tokenRecordModel
		tokenRecordModelSelf
	}

	customTokenRecordModel struct {
		*defaultTokenRecordModel
	}
)

// NewTokenRecordModel returns a model for the database table.
func NewTokenRecordModel(conn sqlx.SqlConn) TokenRecordModel {
	return &customTokenRecordModel{
		defaultTokenRecordModel: newTokenRecordModel(conn),
	}
}
