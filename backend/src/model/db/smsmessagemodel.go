package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ SmsMessageModel = (*customSmsMessageModel)(nil)

type (
	// SmsMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSmsMessageModel.
	SmsMessageModel interface {
		smsMessageModel
		smsMessageModelSelf
	}

	customSmsMessageModel struct {
		*defaultSmsMessageModel
	}
)

// NewSmsMessageModel returns a model for the database table.
func NewSmsMessageModel(conn sqlx.SqlConn) SmsMessageModel {
	return &customSmsMessageModel{
		defaultSmsMessageModel: newSmsMessageModel(conn),
	}
}
