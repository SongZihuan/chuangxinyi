package db

import "github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

var _ AnnouncementModel = (*customAnnouncementModel)(nil)

type (
	// AnnouncementModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAnnouncementModel.
	AnnouncementModel interface {
		announcementModel
		announcementModelSelf
	}

	customAnnouncementModel struct {
		*defaultAnnouncementModel
	}
)

// NewAnnouncementModel returns a model for the database table.
func NewAnnouncementModel(conn sqlx.SqlConn) AnnouncementModel {
	return &customAnnouncementModel{
		defaultAnnouncementModel: newAnnouncementModel(conn),
	}
}

func NewAnnouncementModelWithSession(session sqlx.Session) AnnouncementModel {
	return &customAnnouncementModel{
		defaultAnnouncementModel: newAnnouncementModel(sqlx.NewSqlConnFromSession(session)),
	}
}
