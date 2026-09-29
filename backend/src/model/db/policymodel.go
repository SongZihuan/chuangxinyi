package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ PolicyModel = (*customPolicyModel)(nil)

type (
	// PolicyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPolicyModel.
	PolicyModel interface {
		policyModel
		policyModelSelf
	}

	customPolicyModel struct {
		*defaultPolicyModel
	}
)

// NewPolicyModel returns a model for the database table.
func NewPolicyModel(conn sqlx.SqlConn) PolicyModel {
	return &customPolicyModel{
		defaultPolicyModel: newPolicyModel(conn),
	}
}

func NewPolicyModelWithSession(session sqlx.Session) PolicyModel {
	return &customPolicyModel{
		defaultPolicyModel: newPolicyModel(sqlx.NewSqlConnFromSession(session)),
	}
}
