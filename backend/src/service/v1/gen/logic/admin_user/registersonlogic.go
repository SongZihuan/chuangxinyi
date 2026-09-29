package admin_user

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"github.com/google/uuid"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type RegisterSonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterSonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterSonLogic {
	return &RegisterSonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterSonLogic) RegisterSon(req *types.AdminRegisterSonReq) (resp *types.SuccessResp, err error) {
	father, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if father.SonLevel >= config.BackendConfig.MySQL.SonLevelLimit {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "子用户层级过多"),
		}, nil
	}

	fatherToken, ok := l.ctx.Value("X-Token").(string)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token")
	}

	phone, err := jwt.ParserPhoneToken(req.PhoneToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if phone.WebID != warp.UserCenterWebsite {
		return nil, respmsg.JWTError.New("bad website")
	}

	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	userModel := db.NewUserModel(mysql.MySQLConn)

	count, err := userModel.GetSonCount(l.ctx, father.Id)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if count > config.BackendConfig.MySQL.SonUserLimit {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "子用户过多"),
		}, nil
	}

	if !req.NewWallet {
		countWallet, err := userModel.CountSameWalletUser(l.ctx, father.WalletId)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		} else if countWallet > config.BackendConfig.MySQL.SameWalletUserLimit {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "共享钱包用户过多"),
			}, nil
		}
	}

	lock := fmt.Sprintf("phone:%s", phone)
	if !redis.AcquireLock(l.ctx, lock, time.Minute*2) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneHasBeenRegister, "无法上锁，手机号重复注册"),
		}, nil
	}
	defer redis.ReleaseLock(lock)

	_, err = phoneModel.FindByPhone(l.ctx, phone.Phone)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if err == nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneHasBeenRegister, "数据库存在手机号，手机号重复注册"),
		}, nil
	}

	userIDUUID, success := redis.GenerateUUIDMore(l.ctx, "user", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		userModel := db.NewUserModel(mysql.MySQLConn)
		_, err := userModel.FindOneByUidWithoutDelete(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserRegisterAgain, "生成用户uuid错误"),
		}, nil
	}

	userID := userIDUUID.String()

	var roleID int64
	fatherRole := action.GetRoleWithoutBanned(father.RoleId, father.IsAdmin)
	r := action.GetRole(req.RoleID, false)
	if r.Sign != config.BackendConfig.Admin.AnonymousRole.RoleSign {
		roleID = r.ID
	} else {
		roleID = fatherRole.ID
	}

	var walletID int64
	var newUserID int64
	var son *db.User

	err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		userModel := db.NewUserModelWithSession(session)
		phoneModel := db.NewPhoneModelWithSession(session)

		if req.NewWallet {
			newWallet := &db.Wallet{
				Balance:   0,
				NotBilled: 0,
				HasBilled: 0,
				Billed:    0,
			}

			res1, err := walletModel.Insert(l.ctx, newWallet)
			if err != nil {
				return err
			}

			walletID, err = res1.LastInsertId()
			if err != nil {
				return err
			}
		} else {
			walletID = father.WalletId
		}

		var rootFather int64
		if father.RootFatherId.Valid {
			rootFather = father.RootFatherId.Int64
		} else {
			rootFather = father.Id
		}

		son = &db.User{
			Uid:             userID,
			SonLevel:        father.SonLevel + 1,
			Status:          db.UserStatus_Normal,
			Signin:          false,
			TokenExpiration: config.BackendConfig.Jwt.User.ExpiresSecond,
			IsAdmin:         false,
			RoleId:          roleID,
			FatherId: sql.NullInt64{
				Valid: true,
				Int64: father.Id,
			},
			RootFatherId: sql.NullInt64{
				Valid: true,
				Int64: rootFather,
			},
			InviteId: father.InviteId,
			WalletId: walletID,
		}

		res1, err := userModel.InsertCh(l.ctx, son)
		if err != nil {
			return err
		}

		newUserID, err = res1.LastInsertId()
		if err != nil {
			return err
		}

		newPhone := &db.Phone{
			UserId: newUserID,
			Phone:  phone.Phone,
		}

		_, err = phoneModel.InsertSafe(l.ctx, newPhone)
		if err != nil {
			return err
		}

		return nil
	})
	if errors.Is(err, db.PhoneRepeat) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneRepeat, "手机号重复注册"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	var token string
	var subtype string
	if son.RootFatherId.Int64 == father.Id {
		token, err = jwt.CreateUserToken(l.ctx, userID, son.Signin, son.TokenExpiration, jwt2.UserRootFatherToken, fatherToken, 0)
		subtype = jwt2.UserRootFatherTokenString
	} else {
		token, err = jwt.CreateUserToken(l.ctx, userID, son.Signin, son.TokenExpiration, jwt2.UserFatherToken, fatherToken, 0)
		subtype = jwt2.UserFatherTokenString
	}

	sender.PhoneSendBind(phone.Phone)

	// 因为MessageSendLoginCenter不能用go，为保证顺序MessageSendRegister也不用
	sender.MessageSendSonRegister(newUserID, phone.Phone)

	sender.MessageSendLoginCenter(newUserID, l.ctx)

	audit.NewUserAudit(father.Id, "管理员注册子用户成功")
	audit.NewUserAudit(son.Id, "子用户注册成功")

	logger.Logger.WXInfo("新注册子用户：%s", phone)
	_ = LogMsg(true, "%s收到新子用户注册成功：%s", config.BackendConfig.User.ReadableName, phone)

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:    UserToken,
			Token:   token,
			SubType: subtype,
		},
	}, nil
}

func LogMsg(atall bool, text string, args ...any) errors.WTError {
	return logger.WxRobotSendNotRecord(config.BackendConfig.WXRobot.NewUserLog, fmt.Sprintf(text, args...), atall)
}
