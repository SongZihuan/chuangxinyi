package register

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"github.com/google/uuid"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.SuccessResp, err error) {
	phone, err := jwt.ParserPhoneToken(req.PhoneToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if phone.WebID != warp.UserCenterWebsite {
		return nil, respmsg.JWTError.New("bad website")
	}

	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	var inviteUser *db.User
	if len(req.InviteID) == 0 {
		inviteUser = nil
	} else {
		inviteUser, err = utils2.FindUser(l.ctx, req.InviteID, false)
		if errors.Is(err, utils2.UserNotFound) {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InviteUserNotFound, "邀请用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	lock := fmt.Sprintf("phone:%s", phone)
	if !redis.AcquireLock(l.ctx, lock, time.Minute*2) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneHasBeenRegister, "上锁失败，手机号重复注册"),
		}, nil
	}
	defer redis.ReleaseLock(lock)

	_, err = phoneModel.FindByPhone(l.ctx, phone.Phone)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if err == nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneHasBeenRegister, "手机号已存在数据库，手机号重复注册"),
		}, nil
	}

	userIDByte, success := redis.GenerateUUIDMore(l.ctx, "user", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		userModel := db.NewUserModel(mysql.MySQLConn)
		_, err := userModel.FindOneByUidWithoutDelete(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}
		return false
	})
	if !success {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserRegisterAgain, "用户uuid生成失败"),
		}, nil
	}

	userID := userIDByte.String()

	// 不用GetRoleBySign，要精确找到角色
	role, ok := model.RolesSign()[config.BackendConfig.Admin.UserRole.RoleSign]
	if !ok {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SystemError, "用户默认角色未找到"),
		}, nil
	}

	var newUserID int64
	var newUser *db.User
	err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		userModel := db.NewUserModelWithSession(session)
		phoneModel := db.NewPhoneModelWithSession(session)
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
		newWalletID, err := res1.LastInsertId()
		if err != nil {
			return err
		}
		newUser = &db.User{
			Uid:             userID,
			SonLevel:        0,
			Status:          db.UserStatus_Normal,
			Signin:          false,
			IsAdmin:         false,
			RoleId:          role.ID,
			TokenExpiration: config.BackendConfig.Jwt.User.ExpiresSecond,
			WalletId:        newWalletID,
		}
		if inviteUser != nil {
			newUser.InviteId = sql.NullInt64{
				Valid: true,
				Int64: inviteUser.Id,
			}
		}
		res2, err := userModel.InsertCh(l.ctx, newUser)
		if err != nil {
			return err
		}
		newUserID, err = res2.LastInsertId()
		if err != nil {
			return err
		}
		// 手机最后插入
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

	token, err := jwt.CreateUserToken(l.ctx, userID, newUser.Signin, newUser.TokenExpiration, jwt2.UserRootToken, "", 0)

	sender.PhoneSendBind(phone.Phone)

	// 因为MessageSendLoginCenter不能用go，为保证顺序MessageSendRegister也不用
	sender.MessageSendRegister(newUserID, phone.Phone)
	sender.MessageSendLoginCenter(newUserID, l.ctx)

	audit.NewUserAudit(newUserID, "用户注册成功")

	logger.Logger.WXInfo("用户注册成功：%s", phone.Phone)
	_ = LogMsg(true, "%s收到新普通用户注册成功：%s", config.BackendConfig.User.ReadableName, phone.Phone)

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:    UserToken,
			Token:   token,
			SubType: jwt2.UserRootTokenString,
		},
	}, nil
}

func LogMsg(atall bool, text string, args ...any) errors.WTError {
	return logger.WxRobotSendNotRecord(config.BackendConfig.WXRobot.NewUserLog, fmt.Sprintf(text, args...), atall)
}
