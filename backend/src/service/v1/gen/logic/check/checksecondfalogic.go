package check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CheckSecondFALogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckSecondFALogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckSecondFALogic {
	return &CheckSecondFALogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckSecondFALogic) CheckSecondFA(req *types.CheckSecondFAReq, r *http.Request) (resp *types.SuccessResp, err error) {
	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	userData, _, err := jwt.ParserUserToken(l.ctx, req.Token)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	if web.ID != warp.UserCenterWebsite && web.ID != userData.WebsiteID {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadUserToken, "错误的用户Token，外站授权不匹配"),
		}, nil
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
	user, err := userModel.FindOneByUidWithoutDelete(l.ctx, userData.UserID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadUserToken, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if db.IsBanned(user) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadUserToken, "用户已封禁"),
		}, nil
	}

	secondfa, err := secondfaModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Bad2FACode, "用户未启用2FA"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if !secondfa.Secret.Valid {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Bad2FACode, "用户未启用2FA"),
		}, nil
	}

	if config.BackendConfig.GetMode() == config.RunModeDevelop && req.Code == "123456" {
		// 直接通过
	} else {
		if !utils.CheckTOTP(secondfa.Secret.String, req.Code) {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Bad2FACode, "2FA验证失败"),
			}, nil
		}
	}

	token, err := jwt.CreateCheck2FAToken(user.Uid, web.ID)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:  CheckSecondFAToken,
			Token: token,
		},
	}, nil
}
