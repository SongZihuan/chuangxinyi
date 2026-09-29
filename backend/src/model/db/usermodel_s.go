package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"github.com/wuntsong-org/go-zero-plus/core/stringx"
	"strings"
	"time"
)

const (
	UserStatus_Register int64 = 1
	UserStatus_Normal   int64 = 2
	UserStatus_Banned   int64 = 3
	UserStatus_Delete   int64 = 4
	UserStatus_Freeze   int64 = 5
)

var UserStatusMap = map[int64]string{
	UserStatus_Register: "REGISTER",
	UserStatus_Normal:   "NORMAL",
	UserStatus_Banned:   "BANNED",
	UserStatus_Delete:   "DELETE",
	UserStatus_Freeze:   "FREEZE",
}

var UserStatusStartMap = map[string]int64{
	"REGISTER": UserStatus_Register,
	"NORMAL":   UserStatus_Normal,
	"BANNED":   UserStatus_Banned,
	"DELETE":   UserStatus_Delete,
	"FREEZE":   UserStatus_Freeze,
}

func IsUserStatus(userStatus string) bool {
	_, ok := UserStatusStartMap[userStatus]
	return ok
}

func IsKeepInfoStatus(status int64) bool {
	return IsOKStatus(status) || status == UserStatus_Freeze
}

func IsOKStatus(status int64) bool {
	return status == UserStatus_Normal || status == UserStatus_Register
}

func IsBannedStatus(status int64) bool {
	return !IsOKStatus(status)
}

func IsBanned(user *User) bool {
	if user.IsAdmin {
		return false
	}
	return IsBannedStatus(user.Status)
}

type UserEasy struct {
	ID             int64          `db:"id"`
	UID            string         `db:"uid"`
	RoleID         int64          `db:"role_id"`
	IsAdmin        bool           `db:"is_admin"`
	Phone          sql.NullString `db:"phone"`
	UserName       sql.NullString `db:"user_name"`
	NickName       sql.NullString `db:"nickname"`
	Header         sql.NullString `db:"header"`
	Email          sql.NullString `db:"email"`
	UserRealName   sql.NullString `db:"user_real_name"`
	CompanyName    sql.NullString `db:"company_name"`
	UnionID        sql.NullString `db:"union_id"`
	WeChatNickName sql.NullString `db:"wxname"`
	WeChatHeader   sql.NullString `db:"wxheader"`
	SignIn         bool           `db:"signin"`
	Status         int64          `db:"status"`
	TokenExpire    int64          `db:"token_expire"`
	WalletID       int64          `db:"wallet_id"`
	CreateAt       time.Time      `db:"create_at"`
}

type UncleUserEasy struct {
	ID             int64          `db:"id"`
	UID            string         `db:"uid"`
	RoleID         int64          `db:"role_id"`
	IsAdmin        bool           `db:"is_admin"`
	Phone          sql.NullString `db:"phone"`
	UserName       sql.NullString `db:"user_name"`
	NickName       sql.NullString `db:"nickname"`
	Header         sql.NullString `db:"header"`
	Email          sql.NullString `db:"email"`
	UserRealName   sql.NullString `db:"user_real_name"`
	CompanyName    sql.NullString `db:"company_name"`
	UnionID        sql.NullString `db:"union_id"`
	WeChatNickName sql.NullString `db:"wxname"`
	WeChatHeader   sql.NullString `db:"wxheader"`
	SignIn         bool           `db:"signin"`
	Status         int64          `db:"status"`
	TokenExpire    int64          `db:"token_expire"`
	WalletID       int64          `db:"wallet_id"`
	CreateAt       time.Time      `db:"create_at"`
	UncleTag       string         `db:"uncle_tag"`    // 仅限叔叔 侄子有用
	UncleStatus    int64          `db:"uncle_status"` // 仅限叔叔 侄子有用
}

type (
	userModelSelf interface {
		GetSonList(ctx context.Context, fatherID int64, status []string, src string) ([]User, error)
		FindAdminWithoutDelete(ctx context.Context, limit int64) ([]User, error)
		HaveAny(ctx context.Context) (bool, error)
		FindOneByUidWithoutDelete(ctx context.Context, uid string) (*User, error)
		FindOneByIDWithoutDelete(ctx context.Context, id int64) (*User, error)
		FindOneByFatherIDWithoutDelete(ctx context.Context, fatherID int64) ([]User, error)
		GetUserEasyList(ctx context.Context, status []string, src string, page int64, pageSize int64, startTime int64, endTime int64) ([]UserEasy, error)
		GetUserEasyCount(ctx context.Context, status []string, src string, startTime int64, endTime int64) (int64, error)
		FindUserEasyByIDWithoutDelete(ctx context.Context, id int64) (*UserEasy, error)
		FindUserEasyByUidWithoutDelete(ctx context.Context, uid string) (*UserEasy, error)
		GetSonUserEasyList(ctx context.Context, fatherID int64, status []string, src string) ([]UserEasy, error)
		GetNephewUserEasyList(ctx context.Context, uncleID int64, status []string, src string) ([]UncleUserEasy, error)
		GetUncleUserEasyList(ctx context.Context, userID int64, status []string, src string) ([]UncleUserEasy, error)
		GetUserInviteCount(ctx context.Context, userID int64) (int64, error)
		GetInviteUserEasyList(ctx context.Context, inviteID int64, status []string, src string, page int64, pageSize int64) ([]UserEasy, error)
		CountInviteUserEasyList(ctx context.Context, inviteID int64, status []string, src string) (int64, error)
		GetSameWalletUserList(ctx context.Context, walletID int64) ([]User, error)
		CountSameWalletUser(ctx context.Context, walletID int64) (int64, error)
		GetSonCount(ctx context.Context, fatherID int64) (int64, error)
		UpdateWithoutStatus(ctx context.Context, newData *User) error
		UpdateChWithoutStatus(ctx context.Context, newData *User) error
		InsertCh(ctx context.Context, data *User) (sql.Result, error)
		UpdateCh(ctx context.Context, newData *User) error
	}
)

var userRowsWithPlaceHolderWithStatus = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`", "`status`"), "=?,") + "=?"

func (m *defaultUserModel) UpdateWithoutStatus(ctx context.Context, newData *User) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolderWithStatus)
	_, err := m.conn.ExecCtx(ctx, query, newData.Uid, newData.Signin, newData.SonLevel, newData.FatherId, newData.RootFatherId, newData.InviteId, newData.WalletId, newData.TokenExpiration, newData.RoleId, newData.IsAdmin, newData.Remark, newData.DeleteAt, newData.Id)
	return err
}

func (m *defaultUserModel) InsertCh(ctx context.Context, data *User) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	return ret, err
}

func (m *defaultUserModel) UpdateCh(ctx context.Context, newData *User) error {
	err := m.Update(ctx, newData)
	UpdateUser(newData.Id, m.conn, nil)
	return err
}

func (m *defaultUserModel) UpdateChWithoutStatus(ctx context.Context, newData *User) error {
	err := m.UpdateWithoutStatus(ctx, newData)
	UpdateUser(newData.Id, m.conn, nil)
	return err
}

func (m *defaultUserModel) HaveAny(ctx context.Context) (bool, error) {
	query := fmt.Sprintf("select %s from %s where delete_at is null order by id desc limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return true, nil
	case sqlc.ErrNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (m *defaultUserModel) FindOneByIDWithoutDelete(ctx context.Context, id int64) (*User, error) {
	var resp User
	query := fmt.Sprintf("select %s from %s where `id` = ? and `delete_at` is null order by id desc limit 1", userRows, m.table)
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

func (m *defaultUserModel) FindAdminWithoutDelete(ctx context.Context, limit int64) ([]User, error) {
	var resp []User
	cond := where.NewCond(m.table, userFieldNames).Add("`is_admin`=true")
	query := fmt.Sprintf("select %s from %s where %s order by id desc %s", userRows, m.table, cond, where.NewLimit(limit))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []User{}, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByUidWithoutDelete(ctx context.Context, uid string) (*User, error) {
	var resp User
	query := fmt.Sprintf("select %s from %s where `uid` = ? and delete_at is null order by id desc limit 1", userRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) GetSameWalletUserList(ctx context.Context, walletID int64) ([]User, error) {
	var resp []User
	cond := where.NewCond(m.table, userFieldNames).LinkID(walletID, "wallet_id")
	query := fmt.Sprintf("select %s from %s where %s order by id desck %s", userRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SameWalletUserLimit*2))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []User{}, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) CountSameWalletUser(ctx context.Context, walletID int64) (int64, error) {
	var resp OneInt
	query := fmt.Sprintf("select COUNT(id) as res from %s where `wallet_id` = ? and delete_at is null", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, walletID)
	switch err {
	case nil:
		return resp.Res, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}

func (m *defaultUserModel) FindOneByFatherIDWithoutDelete(ctx context.Context, fatherID int64) ([]User, error) {
	var resp []User
	query := fmt.Sprintf("select %s from %s where u.delete_at is null and u.father_id = ? order by id desc limit 1", userRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, fatherID)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []User{}, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) GetSonList(ctx context.Context, fatherID int64, status []string, src string) ([]User, error) {
	var resp []User
	var err error

	cond := where.NewCond(m.table, userFieldNames).HasFatherID(fatherID).Like(src, true, "uid", "phone", "user_name", "nickname", "email", "user_real_name", "company_name", "union_id", "wxname").Int64In("status", utils.ListConvert(status, UserStatusStartMap))
	query := fmt.Sprintf("select %s from %s where %s %s", userRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SonUserLimit*2))
	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []User{}, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) GetSonCount(ctx context.Context, fatherID int64) (int64, error) {
	var resp OneInt
	var err error

	cond := where.NewCond(m.table, userFieldNames).HasFatherID(fatherID)
	query := fmt.Sprintf("select COUNT(id) as res from %s where %s", m.table, cond)
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

const UserEasyQuery = "u.id as id, u.wallet_id as wallet_id, u.token_expiration as token_expire, u.signin as signin, u.uid as uid, u.role_id as role_id, u.is_admin as is_admin, u.status as status, u.create_at as create_at, un.username as user_name, p.phone as phone, e.email as email, w.union_id as union_id, w.nickname as wxname, w.headimgurl as wxheader, i.user_name as user_real_name, cm.company_name as company_name, n.nickname as nickname, h.header as header"
const UserEasyJoin = "left join phone p on u.id = p.user_id and p.delete_at is null left join email e on u.id = e.user_id and e.delete_at is null left join idcard i on u.id = i.user_id and i.delete_at is null left join wechat w on u.id = w.user_id and w.delete_at is null left join company cm on u.id = cm.user_id and cm.delete_at is null left join nickname n on u.id = n.user_id and n.delete_at is null left join header h on u.id = h.user_id and h.delete_at is null left join username un on u.id = un.user_id and un.delete_at is null "
const UserEasyOrderBy = "u.create_at desc, p.create_at desc, e.create_at desc, i.create_at desc, w.create_at desc, cm.create_at desc, n.create_at desc, h.create_at desc, un.create_at desc "

func (m *defaultUserModel) GetUserEasyList(ctx context.Context, status []string, src string, page int64, pageSize int64, startTime int64, endTime int64) ([]UserEasy, error) {
	var resp []UserEasy
	var err error

	cond := where.NewCondWithoutDeleteAtWithStruct(m.table, UserEasy{}).Add("u.delete_at is null").TimeBetween("u.create_at", startTime, endTime).Like(src, true, "uid", "phone", "user_name", "nickname", "email", "user_real_name", "company_name", "union_id", "wxname").Int64In("status", utils.ListConvert(status, UserStatusStartMap))
	query := fmt.Sprintf("select %s from `user` as u %s where %s order by %s %s", UserEasyQuery, UserEasyJoin, cond, UserEasyOrderBy, where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) GetUserEasyCount(ctx context.Context, status []string, src string, startTime int64, endTime int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCondWithoutDeleteAtWithStruct(m.table, UserEasy{}).Add("u.delete_at is null").TimeBetween("create_at", startTime, endTime).Like(src, true, "uid", "phone", "user_name", "nickname", "email", "user_real_name", "company_name", "union_id", "wxname").Int64In("status", utils.ListConvert(status, UserStatusStartMap))
	query := fmt.Sprintf("select count(u.id) as res from user as u where %s", cond)

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

func (m *defaultUserModel) FindUserEasyByIDWithoutDelete(ctx context.Context, id int64) (*UserEasy, error) {
	var resp UserEasy
	var err error

	query := fmt.Sprintf("select %s from `user` as u %s where u.delete_at is null and u.id = ? order by %s limit 1", UserEasyQuery, UserEasyJoin, UserEasyOrderBy)
	err = m.conn.QueryRowCtx(ctx, &resp, query, id)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindUserEasyByUidWithoutDelete(ctx context.Context, uid string) (*UserEasy, error) {
	var resp UserEasy
	var err error

	query := fmt.Sprintf("select %s from `user` as u %s where u.delete_at is null and u.uid = ? order by %s limit 1", UserEasyQuery, UserEasyJoin, UserEasyOrderBy)
	err = m.conn.QueryRowCtx(ctx, &resp, query, uid)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) GetInviteUserEasyList(ctx context.Context, inviteID int64, status []string, src string, page int64, pageSize int64) ([]UserEasy, error) {
	var resp []UserEasy
	var err error

	cond := where.NewCondWithoutDeleteAtWithStruct(m.table, UserEasy{}).Add("u.delete_at is null").Add("u.invite_id = %d", inviteID).Like(src, true, "uid", "phone").Int64In("status", utils.ListConvert(status, UserStatusStartMap))
	query := fmt.Sprintf("select %s from `user` as u %s where %s order by %s %s", UserEasyQuery, UserEasyJoin, cond, UserEasyOrderBy, where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []UserEasy{}, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) CountInviteUserEasyList(ctx context.Context, inviteID int64, status []string, src string) (int64, error) {
	var resp OneInt
	var err error

	cond := where.NewCondWithoutDeleteAtWithStruct(m.table, UserEasy{}).Add("u.delete_at is null").Add("u.invite_id = %d", inviteID).Like(src, true, "uid", "phone").Int64In("status", utils.ListConvert(status, UserStatusStartMap))
	query := fmt.Sprintf("select count(u.id) as res from `user` as u %s where %s", UserEasyJoin, cond)

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

func (m *defaultUserModel) GetSonUserEasyList(ctx context.Context, fatherID int64, status []string, src string) ([]UserEasy, error) {
	var resp []UserEasy
	var err error

	cond := where.NewCondWithoutDeleteAtWithStruct(m.table, UserEasy{}).Add("u.delete_at is null").Add("u.father_id = %d", fatherID).Like(src, true, "uid", "phone", "user_name", "nickname", "email", "user_real_name", "company_name", "union_id", "wxname").Int64In("status", utils.ListConvert(status, UserStatusStartMap))
	query := fmt.Sprintf("select %s from `user` as u %s where %s order by %s %s", UserEasyQuery, UserEasyJoin, cond, UserEasyOrderBy, where.NewLimit(config.BackendConfig.MySQL.SonUserLimit*2))
	err = m.conn.QueryRowsCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []UserEasy{}, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) GetNephewUserEasyList(ctx context.Context, uncleID int64, status []string, src string) ([]UncleUserEasy, error) {
	var resp []UncleUserEasy
	var err error

	cond := where.NewCondWithoutDeleteAtWithStruct(m.table, UserEasy{}).Add("u.delete_at is null").Like(src, true, "uid", "phone", "user_name", "nickname", "email", "user_real_name", "company_name", "union_id", "wxname").Int64In("status", utils.ListConvert(status, UserStatusStartMap))
	query := fmt.Sprintf("select %s, uc.uncle_tag as uncle_tag, uc.status as uncle_status from `user` as u  %s inner join uncle uc on u.id = uc.user_id and uc.delete_at is null and uc.uncle_id = ? where %s order by %s %s", UserEasyQuery, UserEasyJoin, cond, UserEasyOrderBy, where.NewLimit(config.BackendConfig.MySQL.NephewLimit*2))

	err = m.conn.QueryRowsCtx(ctx, &resp, query, uncleID)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []UncleUserEasy{}, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) GetUncleUserEasyList(ctx context.Context, userID int64, status []string, src string) ([]UncleUserEasy, error) {
	var resp []UncleUserEasy
	var err error

	cond := where.NewCondWithoutDeleteAtWithStruct(m.table, UserEasy{}).Add("u.delete_at is null").Like(src, true, "uid", "phone", "user_name", "nickname", "email", "user_real_name", "company_name", "union_id", "wxname").Int64In("status", utils.ListConvert(status, UserStatusStartMap))
	query := fmt.Sprintf("select %s, uc.uncle_tag as uncle_tag, uc.status as uncle_status from `user` as u %s inner join uncle uc on u.id = uc.uncle_id and uc.delete_at is null and uc.user_id = ? where %s order by %s %s", UserEasyQuery, UserEasyJoin, cond, UserEasyOrderBy, where.NewLimit(config.BackendConfig.MySQL.UncleUserLimit*2))

	err = m.conn.QueryRowsCtx(ctx, &resp, query, userID)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []UncleUserEasy{}, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) GetUserInviteCount(ctx context.Context, userID int64) (int64, error) {
	var err error
	var resp OneInt

	query := fmt.Sprintf("select count(u.id) as res from user as u %s where u.delete_at is null and invite_id = ?", UserEasyJoin)
	err = m.conn.QueryRowCtx(ctx, &resp, query, userID)

	switch err {
	case nil:
		return resp.Res, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}
