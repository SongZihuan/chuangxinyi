package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WalletRecordModel = (*customWalletRecordModel)(nil)

type (
	// WalletRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWalletRecordModel.
	WalletRecordModel interface {
		walletRecordModel
		walletRecordModelSelf
	}

	customWalletRecordModel struct {
		*defaultWalletRecordModel
	}
)

// NewWalletRecordModel returns a model for the database table.
func NewWalletRecordModel(conn sqlx.SqlConn) WalletRecordModel {
	return &customWalletRecordModel{
		defaultWalletRecordModel: newWalletRecordModel(conn),
	}
}

func NewWalletRecordModelWithSession(session sqlx.Session) WalletRecordModel {
	return &customWalletRecordModel{
		defaultWalletRecordModel: newWalletRecordModel(sqlx.NewSqlConnFromSession(session)),
	}
}
