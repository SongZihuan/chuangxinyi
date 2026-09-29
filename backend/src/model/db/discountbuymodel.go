package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ DiscountBuyModel = (*customDiscountBuyModel)(nil)

type (
	// DiscountBuyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDiscountBuyModel.
	DiscountBuyModel interface {
		discountBuyModel
		discountBuyModelSelf
	}

	customDiscountBuyModel struct {
		*defaultDiscountBuyModel
	}
)

// NewDiscountBuyModel returns a model for the database table.
func NewDiscountBuyModel(conn sqlx.SqlConn) DiscountBuyModel {
	return &customDiscountBuyModel{
		defaultDiscountBuyModel: newDiscountBuyModel(conn),
	}
}
