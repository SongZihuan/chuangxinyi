package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ TitleModel = (*customTitleModel)(nil)

type (
	// TitleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTitleModel.
	TitleModel interface {
		titleModel
		titleModelSelf
	}

	customTitleModel struct {
		*defaultTitleModel
	}
)

// NewTitleModel returns a model for the database table.
func NewTitleModel(conn sqlx.SqlConn) TitleModel {
	return &customTitleModel{
		defaultTitleModel: newTitleModel(conn),
	}
}
