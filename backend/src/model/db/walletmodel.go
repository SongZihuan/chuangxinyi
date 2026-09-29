package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ WalletModel = (*customWalletModel)(nil)

type (
	// WalletModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWalletModel.
	WalletModel interface {
		walletModel
		walletModelSelf
	}

	customWalletModel struct {
		*defaultWalletModel
	}
)

// NewWalletModel returns a model for the database table.
func NewWalletModel(conn sqlx.SqlConn) WalletModel {
	return &customWalletModel{
		defaultWalletModel: newWalletModel(conn),
	}
}

func NewWalletModelWithSession(session sqlx.Session) WalletModel {
	return &customWalletModel{
		defaultWalletModel: newWalletModel(sqlx.NewSqlConnFromSession(session)),
	}
}
