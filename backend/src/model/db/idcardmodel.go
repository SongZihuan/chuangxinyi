package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ IdcardModel = (*customIdcardModel)(nil)

type (
	// IdcardModel is an interface to be customized, add more methods here,
	// and implement the added methods in customIdcardModel.
	IdcardModel interface {
		idcardModel
		idcardModelSelf
	}

	customIdcardModel struct {
		*defaultIdcardModel
	}
)

// NewIdcardModel returns a model for the database table.
func NewIdcardModel(conn sqlx.SqlConn) IdcardModel {
	return &customIdcardModel{
		defaultIdcardModel: newIdcardModel(conn),
	}
}

func NewIdcardModelWithSession(session sqlx.Session) IdcardModel {
	return &customIdcardModel{
		defaultIdcardModel: newIdcardModel(sqlx.NewSqlConnFromSession(session)),
	}
}
