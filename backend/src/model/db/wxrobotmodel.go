package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WxrobotModel = (*customWxrobotModel)(nil)

type (
	// WxrobotModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWxrobotModel.
	WxrobotModel interface {
		wxrobotModel
		wxrobotModelSelf
	}

	customWxrobotModel struct {
		*defaultWxrobotModel
	}
)

// NewWxrobotModel returns a model for the database table.
func NewWxrobotModel(conn sqlx.SqlConn) WxrobotModel {
	return &customWxrobotModel{
		defaultWxrobotModel: newWxrobotModel(conn),
	}
}
