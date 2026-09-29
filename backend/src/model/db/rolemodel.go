package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ RoleModel = (*customRoleModel)(nil)

type (
	// RoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleModel.
	RoleModel interface {
		roleModel
		roleModelSelf
	}

	customRoleModel struct {
		*defaultRoleModel
	}
)

// NewRoleModel returns a model for the database table.
func NewRoleModel(conn sqlx.SqlConn) RoleModel {
	return &customRoleModel{
		defaultRoleModel: newRoleModel(conn),
	}
}

func NewRoleModelWithSession(session sqlx.Session) RoleModel {
	return &customRoleModel{
		defaultRoleModel: newRoleModel(sqlx.NewSqlConnFromSession(session)),
	}
}
