package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WorkOrderModel = (*customWorkOrderModel)(nil)

type (
	// WorkOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkOrderModel.
	WorkOrderModel interface {
		workOrderModel
		workOrderModelSelf
	}

	customWorkOrderModel struct {
		*defaultWorkOrderModel
	}
)

// NewWorkOrderModel returns a model for the database table.
func NewWorkOrderModel(conn sqlx.SqlConn) WorkOrderModel {
	return &customWorkOrderModel{
		defaultWorkOrderModel: newWorkOrderModel(conn),
	}
}
