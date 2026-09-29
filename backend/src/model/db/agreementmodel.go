package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ AgreementModel = (*customAgreementModel)(nil)

type (
	// AgreementModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAgreementModel.
	AgreementModel interface {
		agreementModel
		agreementModelSelf
	}

	customAgreementModel struct {
		*defaultAgreementModel
	}
)

// NewAgreementModel returns a model for the database table.
func NewAgreementModel(conn sqlx.SqlConn) AgreementModel {
	return &customAgreementModel{
		defaultAgreementModel: newAgreementModel(conn),
	}
}
