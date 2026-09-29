package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ NicknameModel = (*customNicknameModel)(nil)

type (
	// NicknameModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNicknameModel.
	NicknameModel interface {
		nicknameModel
		nicknameModelSelf
	}

	customNicknameModel struct {
		*defaultNicknameModel
	}
)

// NewNicknameModel returns a model for the database table.
func NewNicknameModel(conn sqlx.SqlConn) NicknameModel {
	return &customNicknameModel{
		defaultNicknameModel: newNicknameModel(conn),
	}
}
