package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ DiscountModel = (*customDiscountModel)(nil)

type (
	// DiscountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDiscountModel.
	DiscountModel interface {
		discountModel
		discountModelSelf
	}

	customDiscountModel struct {
		*defaultDiscountModel
	}
)

// NewDiscountModel returns a model for the database table.
func NewDiscountModel(conn sqlx.SqlConn) DiscountModel {
	return &customDiscountModel{
		defaultDiscountModel: newDiscountModel(conn),
	}
}
