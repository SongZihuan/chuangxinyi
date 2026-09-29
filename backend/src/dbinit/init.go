package dbinit

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/defaultuid"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/google/uuid"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	"github.com/wuntsong-org/wterrors"
	"time"
)

func CreateFooter() errors.WTError {
	footerModel := db.NewFooterModel(mysql.MySQLConn)
	_, err := footerModel.FindTheNew(context.Background())
	if errors.Is(err, db.ErrNotFound) {
		_, err = footerModel.Insert(context.Background(), &db.Footer{
			Icp1:      config.BackendConfig.Admin.ICP1,
			Icp2:      config.BackendConfig.Admin.ICP2,
			Gongan:    config.BackendConfig.Admin.Gongan,
			Copyright: config.BackendConfig.Admin.Copyright,
		})
		if err != nil {
			return errors.WarpQuick(err)
		}
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func ResetRoleAdmin() errors.WTError {
	var adminID int64

	if len(config.BackendConfig.Admin.RootRole.RoleName) == 0 {
		return errors.Errorf("bad role name for admin")
	}

	if len(config.BackendConfig.Admin.UserRole.RoleName) == 0 {
		return errors.Errorf("bad role name for user")
	}

	if len(config.BackendConfig.Admin.AnonymousRole.RoleName) == 0 {
		return errors.Errorf("bad role name for anonymous")
	}

	if len(config.BackendConfig.Admin.RootRole.RoleSign) == 0 {
		return errors.Errorf("bad role sign for admin")
	}

	if len(config.BackendConfig.Admin.UserRole.RoleSign) == 0 {
		return errors.Errorf("bad role sign for user")
	}

	if len(config.BackendConfig.Admin.AnonymousRole.RoleSign) == 0 {
		return errors.Errorf("bad role sign for anonymous")
	}

	if !config.BackendConfig.Admin.RootRole.ResetPermission {
		return errors.Errorf("root role must reset permission")
	}

	if !utils.IsPhoneNumber(config.BackendConfig.Admin.AdminPhone) {
		return errors.Errorf("bad phone for admin")
	}

	key := "db:init:admin"
	if !redis.AcquireLock(context.Background(), key, time.Second*30) {
		return nil
	}
	defer redis.ReleaseLock(key)

	roleModel := db.NewRoleModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	userModel := db.NewUserModel(mysql.MySQLConn)

	root, err := roleModel.FindBySignWithoutDelete(context.Background(), config.BackendConfig.Admin.RootRole.RoleSign)
	if errors.Is(err, db.ErrNotFound) {
		resAdmin, err := roleModel.Insert(context.Background(), &db.Role{
			Name:                 config.BackendConfig.Admin.RootRole.RoleName,
			Describe:             config.BackendConfig.Admin.RootRole.RoleDescribe,
			Sign:                 config.BackendConfig.Admin.RootRole.RoleSign,
			NotDelete:            config.BackendConfig.Admin.RootRole.NotDelete,
			NotChangeSign:        config.BackendConfig.Admin.RootRole.NotChangeSign,
			NotChangePermissions: config.BackendConfig.Admin.RootRole.NotChangePermissions,
			Status:               db.RoleStatusOK,
			Permissions:          model.AllPermission().Text(16),
		})
		if err != nil {
			return errors.WarpQuick(err)
		}

		adminID, err = resAdmin.LastInsertId()
		if err != nil {
			return errors.WarpQuick(err)
		}
	} else { // root必须重置
		root.Permissions = model.AllPermission().Text(16)

		err = roleModel.Update(context.Background(), root)
		if err != nil {
			return errors.WarpQuick(err)
		}

		adminID = root.Id
	}

	user, err := roleModel.FindBySignWithoutDelete(context.Background(), config.BackendConfig.Admin.UserRole.RoleSign)
	if errors.Is(err, db.ErrNotFound) {
		_, err = roleModel.Insert(context.Background(), &db.Role{
			Name:                 config.BackendConfig.Admin.UserRole.RoleName,
			Describe:             config.BackendConfig.Admin.UserRole.RoleDescribe,
			Sign:                 config.BackendConfig.Admin.UserRole.RoleSign,
			NotDelete:            config.BackendConfig.Admin.UserRole.NotDelete,
			NotChangeSign:        config.BackendConfig.Admin.UserRole.NotChangeSign,
			NotChangePermissions: config.BackendConfig.Admin.UserRole.NotChangePermissions,
			Status:               db.RoleStatusOK,
			Permissions:          model.UserPermission().Text(16),
		})
		if err != nil {
			return errors.WarpQuick(err)
		}
	} else if config.BackendConfig.Admin.UserRole.ResetPermission {
		user.Permissions = model.UserPermission().Text(16)

		err = roleModel.Update(context.Background(), user)
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	anonymous, err := roleModel.FindBySignWithoutDelete(context.Background(), config.BackendConfig.Admin.AnonymousRole.RoleSign)
	if errors.Is(err, db.ErrNotFound) {
		_, err = roleModel.Insert(context.Background(), &db.Role{
			Name:                 config.BackendConfig.Admin.AnonymousRole.RoleName,
			Describe:             config.BackendConfig.Admin.AnonymousRole.RoleDescribe,
			Sign:                 config.BackendConfig.Admin.AnonymousRole.RoleSign,
			NotDelete:            config.BackendConfig.Admin.AnonymousRole.NotDelete,
			NotChangeSign:        config.BackendConfig.Admin.AnonymousRole.NotChangeSign,
			NotChangePermissions: config.BackendConfig.Admin.AnonymousRole.NotChangePermissions,
			Status:               db.RoleStatusOK,
			Permissions:          model.AnonymousPermission().Text(16),
		})
		if err != nil {
			return errors.WarpQuick(err)
		}
	} else if config.BackendConfig.Admin.AnonymousRole.ResetPermission {
		anonymous.Permissions = model.AnonymousPermission().Text(16)

		err = roleModel.Update(context.Background(), anonymous)
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	lock := fmt.Sprintf("phone:%s", config.BackendConfig.Admin.AdminPhone)
	if !redis.AcquireLock(context.Background(), lock, time.Minute*2) {
		return nil
	}
	defer redis.ReleaseLock(lock)

	adminList, err := userModel.FindAdminWithoutDelete(context.Background(), 1)
	if len(adminList) == 0 {
		phone, err := phoneModel.FindByPhone(context.Background(), config.BackendConfig.Admin.AdminPhone)
		if errors.Is(err, db.ErrNotFound) {
			err := createAdmin(adminID)
			if err != nil {
				return errors.WarpQuick(err)
			}
		} else if err != nil {
			return errors.WarpQuick(err)
		} else {
			user, err := userModel.FindOneByIDWithoutDelete(context.Background(), phone.UserId)
			if errors.Is(err, db.ErrNotFound) {
				phone.DeleteAt = sql.NullTime{
					Valid: true,
					Time:  time.Now(),
				}

				err := phoneModel.Update(context.Background(), phone)
				if err != nil {
					return errors.WarpQuick(err)
				}

				err = createAdmin(adminID)
				if err != nil {
					return errors.WarpQuick(err)
				}
			} else if err != nil {
				return errors.WarpQuick(err)
			} else if !user.IsAdmin {
				user.IsAdmin = true
				defaultuid.DefaultUID = user.Uid
				err := userModel.UpdateWithoutStatus(context.Background(), user)
				if err != nil {
					return errors.WarpQuick(err)
				}
			}
		}
	} else {
		_, err := phoneModel.FindByPhone(context.Background(), config.BackendConfig.Admin.AdminPhone)
		if errors.Is(err, db.ErrNotFound) {
			err := createAdmin(adminID)
			if err != nil {
				return errors.WarpQuick(err)
			}
		} else if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return nil
}

func createAdmin(adminRoleID int64) errors.WTError {
	userIDByte, success := redis.GenerateUUIDMore(context.Background(), "user", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		userModel := db.NewUserModel(mysql.MySQLConn)
		_, err := userModel.FindOneByUidWithoutDelete(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}
		return false
	})
	if !success {
		return errors.Errorf("generate user uid fail")
	}

	defaultuid.DefaultUID = userIDByte.String()

	err := mysql.MySQLConn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		userModel := db.NewUserModelWithSession(session)
		phoneModel := db.NewPhoneModelWithSession(session)

		newWallet := &db.Wallet{
			Balance:   0,
			NotBilled: 0,
			HasBilled: 0,
			Billed:    0,
		}

		resWallet, err := walletModel.Insert(context.Background(), newWallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		walletID, err := resWallet.LastInsertId()
		if err != nil {
			return errors.WarpQuick(err)
		}

		newUser := &db.User{
			Uid:             userIDByte.String(),
			SonLevel:        0,
			Status:          db.UserStatus_Normal,
			Signin:          true,
			IsAdmin:         true,
			RoleId:          adminRoleID,
			WalletId:        walletID,
			TokenExpiration: config.BackendConfig.Jwt.User.ExpiresSecond,
		}

		resUser, err := userModel.InsertCh(context.Background(), newUser)
		if err != nil {
			return errors.WarpQuick(err)
		}

		userID, err := resUser.LastInsertId()
		if err != nil {
			return errors.WarpQuick(err)
		}

		// 手机最后插入
		newPhone := &db.Phone{
			UserId: userID,
			Phone:  config.BackendConfig.Admin.AdminPhone,
		}

		_, err = phoneModel.InsertWithDelete(context.Background(), newPhone) // 采用强制插入
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
