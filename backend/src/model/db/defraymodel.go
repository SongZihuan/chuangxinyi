package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ DefrayModel = (*customDefrayModel)(nil)

type (
	// DefrayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDefrayModel.
	DefrayModel interface {
		defrayModel
		defrayModelSelf
	}

	customDefrayModel struct {
		*defaultDefrayModel
	}
)

// NewDefrayModel returns a model for the database table.
func NewDefrayModel(conn sqlx.SqlConn) DefrayModel {
	return &customDefrayModel{
		defaultDefrayModel: newDefrayModel(conn),
	}
}

func NewDefrayModelWithSession(session sqlx.Session) DefrayModel {
	return &customDefrayModel{
		defaultDefrayModel: newDefrayModel(sqlx.NewSqlConnFromSession(session)),
	}
}
