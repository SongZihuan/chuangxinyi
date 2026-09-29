package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ InvoiceModel = (*customInvoiceModel)(nil)

type (
	// InvoiceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInvoiceModel.
	InvoiceModel interface {
		invoiceModel
		invoiceModelSelf
	}

	customInvoiceModel struct {
		*defaultInvoiceModel
	}
)

// NewInvoiceModel returns a model for the database table.
func NewInvoiceModel(conn sqlx.SqlConn) InvoiceModel {
	return &customInvoiceModel{
		defaultInvoiceModel: newInvoiceModel(conn),
	}
}

func NewInvoiceModelWithSession(session sqlx.Session) InvoiceModel {
	return &customInvoiceModel{
		defaultInvoiceModel: newInvoiceModel(sqlx.NewSqlConnFromSession(session)),
	}
}
