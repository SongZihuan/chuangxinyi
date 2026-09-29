package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ LoginControllerModel = (*customLoginControllerModel)(nil)

type (
	// LoginControllerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoginControllerModel.
	LoginControllerModel interface {
		loginControllerModel
		loginControllerModelSelf
	}

	customLoginControllerModel struct {
		*defaultLoginControllerModel
	}
)

// NewLoginControllerModel returns a model for the database table.
func NewLoginControllerModel(conn sqlx.SqlConn) LoginControllerModel {
	return &customLoginControllerModel{
		defaultLoginControllerModel: newLoginControllerModel(conn),
	}
}
