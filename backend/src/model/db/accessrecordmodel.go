package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ AccessRecordModel = (*customAccessRecordModel)(nil)

type (
	// AccessRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccessRecordModel.
	AccessRecordModel interface {
		accessRecordModel
		accessRecordModelSelf
	}

	customAccessRecordModel struct {
		*defaultAccessRecordModel
	}
)

// NewAccessRecordModel returns a model for the database table.
func NewAccessRecordModel(conn sqlx.SqlConn) AccessRecordModel {
	return &customAccessRecordModel{
		defaultAccessRecordModel: newAccessRecordModel(conn),
	}
}
