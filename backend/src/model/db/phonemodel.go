package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ PhoneModel = (*customPhoneModel)(nil)

type (
	// PhoneModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPhoneModel.
	PhoneModel interface {
		phoneModel
		phoneModelSelf
	}

	customPhoneModel struct {
		*defaultPhoneModel
	}
)

// NewPhoneModel returns a model for the database table.
func NewPhoneModel(conn sqlx.SqlConn) PhoneModel {
	return &customPhoneModel{
		defaultPhoneModel: newPhoneModel(conn),
	}
}

func NewPhoneModelWithSession(session sqlx.Session) PhoneModel {
	return &customPhoneModel{
		defaultPhoneModel: newPhoneModel(sqlx.NewSqlConnFromSession(session)),
	}
}
