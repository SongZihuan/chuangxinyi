package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WxrobotMessageModel = (*customWxrobotMessageModel)(nil)

type (
	// WxrobotMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWxrobotMessageModel.
	WxrobotMessageModel interface {
		wxrobotMessageModel
		wxrobotMessageModelSelf
	}

	customWxrobotMessageModel struct {
		*defaultWxrobotMessageModel
	}
)

// NewWxrobotMessageModel returns a model for the database table.
func NewWxrobotMessageModel(conn sqlx.SqlConn) WxrobotMessageModel {
	return &customWxrobotMessageModel{
		defaultWxrobotMessageModel: newWxrobotMessageModel(conn),
	}
}
