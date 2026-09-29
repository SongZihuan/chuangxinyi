package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type DeleteTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTokenLogic {
	return &DeleteTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTokenLogic) DeleteToken(req *types.DeleteOauth2TokenReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	tokenData, err := jwt.ParserDeleteToken(req.Token)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if tokenData.Type != jwt.TypeDeleteLoginToken {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDeleteToken, "删除token不是用于删除登录token的"),
		}, nil
	}

	loginData, err := jwt.ParserLoginToken(l.ctx, tokenData.Token)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		}, nil
	} else if loginData.UserID != user.Uid {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDeleteToken, "登录token不属于用户"),
		}, nil
	}

	web := (model.Websites())[loginData.WebID]
	if web.ID == warp.UnknownWebsite {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDeleteToken, "登录token不属于用户"),
		}, nil
	}

	err = jwt.DeleteLoginToken(l.ctx, user.Uid, tokenData.Token)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户撤下站点（%s）的一个访问权限", web.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
