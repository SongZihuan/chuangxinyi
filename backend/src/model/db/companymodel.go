package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ CompanyModel = (*customCompanyModel)(nil)

type (
	// CompanyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCompanyModel.
	CompanyModel interface {
		companyModel
		companyModelSelf
	}

	customCompanyModel struct {
		*defaultCompanyModel
	}
)

// NewCompanyModel returns a model for the database table.
func NewCompanyModel(conn sqlx.SqlConn) CompanyModel {
	return &customCompanyModel{
		defaultCompanyModel: newCompanyModel(conn),
	}
}

func NewCompanyModelWithSession(session sqlx.Session) CompanyModel {
	return &customCompanyModel{
		defaultCompanyModel: newCompanyModel(sqlx.NewSqlConnFromSession(session)),
	}
}
