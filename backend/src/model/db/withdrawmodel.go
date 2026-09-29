package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WithdrawModel = (*customWithdrawModel)(nil)

type (
	// WithdrawModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWithdrawModel.
	WithdrawModel interface {
		withdrawModel
		withdrawModelSelf
	}

	customWithdrawModel struct {
		*defaultWithdrawModel
	}
)

// NewWithdrawModel returns a model for the database table.
func NewWithdrawModel(conn sqlx.SqlConn) WithdrawModel {
	return &customWithdrawModel{
		defaultWithdrawModel: newWithdrawModel(conn),
	}
}

func NewWithdrawModelWithSession(session sqlx.Session) WithdrawModel {
	return &customWithdrawModel{
		defaultWithdrawModel: newWithdrawModel(sqlx.NewSqlConnFromSession(session)),
	}
}
