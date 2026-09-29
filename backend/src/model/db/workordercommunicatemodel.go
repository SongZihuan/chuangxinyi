package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WorkOrderCommunicateModel = (*customWorkOrderCommunicateModel)(nil)

type (
	// WorkOrderCommunicateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkOrderCommunicateModel.
	WorkOrderCommunicateModel interface {
		workOrderCommunicateModel
		workOrderCommunicateModelSelf
	}

	customWorkOrderCommunicateModel struct {
		*defaultWorkOrderCommunicateModel
	}
)

// NewWorkOrderCommunicateModel returns a model for the database table.
func NewWorkOrderCommunicateModel(conn sqlx.SqlConn) WorkOrderCommunicateModel {
	return &customWorkOrderCommunicateModel{
		defaultWorkOrderCommunicateModel: newWorkOrderCommunicateModel(conn),
	}
}
