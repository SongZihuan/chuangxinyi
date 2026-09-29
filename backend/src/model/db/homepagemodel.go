package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ HomepageModel = (*customHomepageModel)(nil)

type (
	// HomepageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageModel.
	HomepageModel interface {
		homepageModel
		homepageModelSelf
	}

	customHomepageModel struct {
		*defaultHomepageModel
	}
)

// NewHomepageModel returns a model for the database table.
func NewHomepageModel(conn sqlx.SqlConn) HomepageModel {
	return &customHomepageModel{
		defaultHomepageModel: newHomepageModel(conn),
	}
}
