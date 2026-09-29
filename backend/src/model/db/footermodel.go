package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ FooterModel = (*customFooterModel)(nil)

type (
	// FooterModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFooterModel.
	FooterModel interface {
		footerModel
		footerModelSelf
	}

	customFooterModel struct {
		*defaultFooterModel
	}
)

// NewFooterModel returns a model for the database table.
func NewFooterModel(conn sqlx.SqlConn) FooterModel {
	return &customFooterModel{
		defaultFooterModel: newFooterModel(conn),
	}
}
