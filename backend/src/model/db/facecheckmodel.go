package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ FaceCheckModel = (*customFaceCheckModel)(nil)

type (
	// FaceCheckModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFaceCheckModel.
	FaceCheckModel interface {
		faceCheckModel
		faceCheckModelSelf
	}

	customFaceCheckModel struct {
		*defaultFaceCheckModel
	}
)

// NewFaceCheckModel returns a model for the database table.
func NewFaceCheckModel(conn sqlx.SqlConn) FaceCheckModel {
	return &customFaceCheckModel{
		defaultFaceCheckModel: newFaceCheckModel(conn),
	}
}
