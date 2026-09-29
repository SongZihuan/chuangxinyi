package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ AddressModel = (*customAddressModel)(nil)

type (
	// AddressModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAddressModel.
	AddressModel interface {
		addressModel
		addressModelSelf
	}

	customAddressModel struct {
		*defaultAddressModel
	}
)

// NewAddressModel returns a model for the database table.
func NewAddressModel(conn sqlx.SqlConn) AddressModel {
	return &customAddressModel{
		defaultAddressModel: newAddressModel(conn),
	}
}
