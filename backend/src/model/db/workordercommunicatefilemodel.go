package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WorkOrderCommunicateFileModel = (*customWorkOrderCommunicateFileModel)(nil)

type (
	// WorkOrderCommunicateFileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkOrderCommunicateFileModel.
	WorkOrderCommunicateFileModel interface {
		workOrderCommunicateFileModel
		workOrderCommunicateFileModelSelf
	}

	customWorkOrderCommunicateFileModel struct {
		*defaultWorkOrderCommunicateFileModel
	}
)

// NewWorkOrderCommunicateFileModel returns a model for the database table.
func NewWorkOrderCommunicateFileModel(conn sqlx.SqlConn) WorkOrderCommunicateFileModel {
	return &customWorkOrderCommunicateFileModel{
		defaultWorkOrderCommunicateFileModel: newWorkOrderCommunicateFileModel(conn),
	}
}
