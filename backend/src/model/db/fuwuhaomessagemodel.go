package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ FuwuhaoMessageModel = (*customFuwuhaoMessageModel)(nil)

type (
	// FuwuhaoMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFuwuhaoMessageModel.
	FuwuhaoMessageModel interface {
		fuwuhaoMessageModel
		fuwuhaoMessageModelSelf
	}

	customFuwuhaoMessageModel struct {
		*defaultFuwuhaoMessageModel
	}
)

// NewFuwuhaoMessageModel returns a model for the database table.
func NewFuwuhaoMessageModel(conn sqlx.SqlConn) FuwuhaoMessageModel {
	return &customFuwuhaoMessageModel{
		defaultFuwuhaoMessageModel: newFuwuhaoMessageModel(conn),
	}
}
