package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ CouponsModel = (*customCouponsModel)(nil)

type (
	// CouponsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponsModel.
	CouponsModel interface {
		couponsModel
		couponsModelSelf
	}

	customCouponsModel struct {
		*defaultCouponsModel
	}
)

// NewCouponsModel returns a model for the database table.
func NewCouponsModel(conn sqlx.SqlConn) CouponsModel {
	return &customCouponsModel{
		defaultCouponsModel: newCouponsModel(conn),
	}
}
