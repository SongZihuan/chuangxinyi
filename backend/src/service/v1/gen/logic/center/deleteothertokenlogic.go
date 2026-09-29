package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type DeleteOtherTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOtherTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOtherTokenLogic {
	return &DeleteOtherTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOtherTokenLogic) DeleteOtherToken(req *types.DeleteOtherToken) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	tokenData, err := jwt.ParserDeleteToken(req.Token)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if tokenData.Type != jwt.TypeDeleteUserToken {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDeleteToken, "删除token不是用于删除用户token的"),
		}, nil
	}

	userData, _, err := jwt.ParserUserToken(l.ctx, tokenData.Token)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		}, nil
	} else if userData.UserID != user.Uid {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDeleteToken, "删除的token不属于该用户"),
		}, nil
	}

	err = jwt.DeleteUserToken(l.ctx, user.Uid, tokenData.Token)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户撤销其中一个登录令牌成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
