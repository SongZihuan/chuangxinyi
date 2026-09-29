package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ AuditModel = (*customAuditModel)(nil)

type (
	// AuditModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuditModel.
	AuditModel interface {
		auditModel
		auditModelSelf
	}

	customAuditModel struct {
		*defaultAuditModel
	}
)

// NewAuditModel returns a model for the database table.
func NewAuditModel(conn sqlx.SqlConn) AuditModel {
	return &customAuditModel{
		defaultAuditModel: newAuditModel(conn),
	}
}
