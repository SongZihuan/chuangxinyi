package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"time"
)

type (
	announcementModelSelf interface {
		GetList(ctx context.Context, src string, show bool, page int64, pageSize int64) ([]Announcement, error)
		GetCount(ctx context.Context, src string, show bool) (int64, error)
		GetNewSortNumber(ctx context.Context) (res int64, err error)
		FindNear(ctx context.Context, sort int64, isUp bool) (res *Announcement, err error)
		FindOneWithoutDelete(ctx context.Context, id int64) (*Announcement, error)
		InsertCh(ctx context.Context, data *Announcement) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Announcement) error
		UpdateDeleteCh(ctx context.Context, data *Announcement) error
	}
)

func IsNowShow(data *Announcement) bool {
	now := time.Now()
	if data.StartAt.After(now) {
		return false
	}

	if data.StopAt.Before(now) {
		return false
	}

	return true
}

func (m *defaultAnnouncementModel) InsertCh(ctx context.Context, data *Announcement) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	if IsNowShow(data) && err == nil {
		id, err := ret.LastInsertId()
		if err == nil {
			data.Id = id
			NewAnnouncement(data)
		}
	}
	return ret, err
}

func (m *defaultAnnouncementModel) UpdateCh(ctx context.Context, data *Announcement) error {
	err := m.Update(ctx, data)
	if IsNowShow(data) && err == nil && !data.DeleteAt.Valid {
		UpdateAnnouncement(data)
	}
	return err
}

func (m *defaultAnnouncementModel) UpdateDeleteCh(ctx context.Context, data *Announcement) error {
	err := m.Update(ctx, data)
	if IsNowShow(data) && err == nil {
		DeleteAnnouncement(data)
	}
	return err
}

func (m *defaultAnnouncementModel) FindOneWithoutDelete(ctx context.Context, id int64) (*Announcement, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `delete_at` is null order by id desc limit 1", announcementRows, m.table)
	var resp Announcement
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAnnouncementModel) GetList(ctx context.Context, src string, show bool, page int64, pageSize int64) (resp []Announcement, err error) {
	cond := where.NewCond(m.table, announcementFieldNames).Like(src, true, "title")
	if show {
		now := time.Now()
		cond = cond.Add("create_at <= %s", where.GetTime(now)).Add("stop_at > %s", where.GetTime(now))
	}

	query := fmt.Sprintf("select %s from %s where %s order by sort %s", announcementRows, m.table, cond, where.NewPage(page, pageSize))
	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Announcement{}, nil
	default:
		return nil, err
	}
}

func (m *defaultAnnouncementModel) GetCount(ctx context.Context, src string, show bool) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, announcementFieldNames).Like(src, true, "title")
	if show {
		now := time.Now()
		cond = cond.Add("create_at <= %s", where.GetTime(now)).Add("stop_at > %s", where.GetTime(now))
	}

	query := fmt.Sprintf("select count(id) as res from %s where %s", m.table, cond)
	err = m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp.Res, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}

func (m *defaultAnnouncementModel) GetNewSortNumber(ctx context.Context) (res int64, err error) {
	var resp OneIntOrNull
	query := fmt.Sprintf("select max(sort) as res from %s where delete_at is null", m.table)
	err = m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp.Res.Int64 + 1, nil
	case sqlc.ErrNotFound:
		return 1, nil
	default:
		return 0, err
	}
}

func (m *defaultAnnouncementModel) FindNear(ctx context.Context, sort int64, isUp bool) (res *Announcement, err error) {
	var resp Announcement
	if isUp {
		query := fmt.Sprintf("select %s from %s where `sort` < ? and delete_at is null order by `sort` desc limit 1", announcementRows, m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, sort)
	} else {
		query := fmt.Sprintf("select %s from %s where `sort` > ? and delete_at is null order by `sort` asc limit 1", announcementRows, m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, sort)
	}

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
