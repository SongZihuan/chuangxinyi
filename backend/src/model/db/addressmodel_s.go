package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/ip"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"time"
)

type (
	addressModelSelf interface {
		FindByUserID(ctx context.Context, userID int64) (*Address, error)
		InsertWithDelete(ctx context.Context, data *Address) (sql.Result, error)
		InsertCh(ctx context.Context, data *Address) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Address) error
	}
)

func (m *defaultAddressModel) InsertCh(ctx context.Context, data *Address) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultAddressModel) UpdateCh(ctx context.Context, data *Address) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultAddressModel) FindByUserID(ctx context.Context, userID int64) (*Address, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", addressRows, m.table)
	var resp Address
	err := m.conn.QueryRowCtx(ctx, &resp, query, userID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAddressModel) InsertWithDelete(ctx context.Context, data *Address) (sql.Result, error) {
	key := fmt.Sprintf("db:insert:%s:%d", m.table, data.UserId)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key)

	updateQuery := fmt.Sprintf("update %s set delete_at=? where user_id = ? and delete_at is null", m.table)
	retUpdate, err := m.conn.ExecCtx(ctx, updateQuery, time.Now(), data.UserId)
	if err != nil {
		return retUpdate, err
	}

	return m.InsertCh(ctx, data)
}

func (a *Address) GetAreaList() []string {
	switch true {
	case a.CountryCode.Valid && a.ProvinceCode.Valid && a.CityCode.Valid && a.DistrictCode.Valid:
		return []string{a.CountryCode.String, a.ProvinceCode.String, a.CityCode.String, a.DistrictCode.String}
	case a.CountryCode.Valid && a.ProvinceCode.Valid && a.CityCode.Valid && !a.DistrictCode.Valid:
		return []string{a.CountryCode.String, a.ProvinceCode.String, a.CityCode.String}
	case a.CountryCode.Valid && a.ProvinceCode.Valid && !a.CityCode.Valid && !a.DistrictCode.Valid:
		return []string{a.CountryCode.String, a.ProvinceCode.String}
	case a.CountryCode.Valid && !a.ProvinceCode.Valid && !a.CityCode.Valid && !a.DistrictCode.Valid:
		return []string{a.CountryCode.String}
	default:
		return []string{}
	}
}

var BadArea = fmt.Errorf("bad area")

func (a *Address) SetAreaCode(area ...string) error {
	if area == nil || len(area) == 0 {
		return nil
	}

	for _, a := range area {
		if len(a) != ip.CityGeoCodeLen || len(a) != ip.CountryGeoCodeLen {
			return BadArea
		}
	}

	if len(area) == 4 {
		a.CountryCode = sql.NullString{
			Valid:  true,
			String: area[0],
		}
		a.ProvinceCode = sql.NullString{
			Valid:  true,
			String: area[1],
		}
		a.CityCode = sql.NullString{
			Valid:  true,
			String: area[2],
		}
		a.DistrictCode = sql.NullString{
			Valid:  true,
			String: area[3],
		}
	} else if len(area) == 3 {
		a.CountryCode = sql.NullString{
			Valid:  true,
			String: area[0],
		}
		a.ProvinceCode = sql.NullString{
			Valid:  true,
			String: area[2],
		}
		a.CityCode = sql.NullString{
			Valid:  true,
			String: area[3],
		}
	} else if len(area) == 2 {
		a.CountryCode = sql.NullString{
			Valid:  true,
			String: area[0],
		}
		a.ProvinceCode = sql.NullString{
			Valid:  true,
			String: area[2],
		}
	} else if len(area) == 1 {
		a.CountryCode = sql.NullString{
			Valid:  true,
			String: area[0],
		}
	} else {
		return BadArea
	}

	return nil
}
