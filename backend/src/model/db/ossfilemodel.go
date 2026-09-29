package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ OssFileModel = (*customOssFileModel)(nil)

type (
	// OssFileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOssFileModel.
	OssFileModel interface {
		ossFileModel
		ossFileModelSelf
	}

	customOssFileModel struct {
		*defaultOssFileModel
	}
)

// NewOssFileModel returns a model for the database table.
func NewOssFileModel(conn sqlx.SqlConn) OssFileModel {
	return &customOssFileModel{
		defaultOssFileModel: newOssFileModel(conn),
	}
}
