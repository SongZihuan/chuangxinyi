package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ UrlPathModel = (*customUrlPathModel)(nil)

type (
	// UrlPathModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUrlPathModel.
	UrlPathModel interface {
		urlPathModel
		urlPathModelSelf
	}

	customUrlPathModel struct {
		*defaultUrlPathModel
	}
)

// NewUrlPathModel returns a model for the database table.
func NewUrlPathModel(conn sqlx.SqlConn) UrlPathModel {
	return &customUrlPathModel{
		defaultUrlPathModel: newUrlPathModel(conn),
	}
}
