package check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/password"
	"gitee.com/wuntsong-auth/backend/src/sender"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CheckPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPasswordLogic {
	return &CheckPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckPasswordLogic) CheckPassword(req *types.CheckPasswordReq) (resp *types.SuccessResp, err error) {
	if req.Type != UserToken {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTokenType, "错误的token类型"),
		}, nil
	}

	if !utils.IsSha256(req.PasswordHash) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPassword, "错误的密码哈希值"),
		}, nil
	}

	ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)
	passwordModel := db.NewPasswordModel(mysql.MySQLConn)
	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)

	user, err := utils2.FindUser(l.ctx, req.UserID, false)
	if errors.Is(err, utils2.UserNotFound) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.UserNotRegister, "用户未找到", "密码错误"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	ctrl, err := ctrlModel.FindByUserID(l.ctx, user.Id)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if !ctrl.AllowPassword {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.UserNotFound, "用户不允许密码登录", "密码错误"),
		}, nil
	}

	if config.BackendConfig.GetMode() == config.RunModeDevelop && req.PasswordHash == password.GetPasswordFirstHash("admin123") {
		// 直接通过
	} else {
		pw, err := passwordModel.FindByUserID(l.ctx, user.Id)
		if errors.Is(err, db.ErrNotFound) {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.PasswordError, "密码错误", "密码错误"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		if !pw.PasswordHash.Valid {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.PasswordError, "密码错误", "密码错误"),
			}, nil
		}

		passwordHash := password.GetPasswordSecondHash(req.PasswordHash, user.Uid)
		if passwordHash != pw.PasswordHash.String {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.PasswordError, "密码错误", "密码错误"),
			}, nil
		}
	}

	var tokenType string
	var token string
	var subType string

	secondfa, err := secondfaModel.FindByUserID(l.ctx, user.Id)
	if !errors.Is(err, db.ErrNotFound) && err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if ctrl.Allow2Fa && secondfa != nil && secondfa.Secret.Valid { // secondfa != nil 表示未找到
		tokenType = SecondFAToken
		token, err = jwt.CreateLogin2FAToken(user.Uid)
		if err != nil {
			return nil, respmsg.JWTError.WarpQuick(err)
		}

		audit.NewUserAudit(user.Id, "用户试图通过密码登录，需要2FA")
	} else {
		tokenType = UserToken
		sender.MessageSendLoginCenter(user.Id, l.ctx)
		sender.WxrobotSendLoginCenter(user.Id, l.ctx)
		sender.FuwuhaoSendLoginCenter(user.Id)

		if user.FatherId.Valid {
			token, err = jwt.CreateUserToken(l.ctx, user.Uid, user.Signin, user.TokenExpiration, jwt2.UserSonToken, "", 0)
			subType = jwt2.UserSonTokenString
		} else {
			token, err = jwt.CreateUserToken(l.ctx, user.Uid, user.Signin, user.TokenExpiration, jwt2.UserRootToken, "", 0)
			subType = jwt2.UserRootTokenString
		}
		if err != nil {
			return nil, respmsg.JWTError.WarpQuick(err)
		}

		audit.NewUserAudit(user.Id, "用户通过密码登录成功")
	}

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:    tokenType,
			Token:   token,
			SubType: subType,
		},
	}, nil
}
