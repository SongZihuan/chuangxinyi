package admin_user

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

func (l *DeleteTokenLogic) DeleteToken(req *types.AdminDeleteUserTokenReq) (resp *types.RespEmpty, err error) {
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
	}

	err = jwt.DeleteUserToken(l.ctx, userData.UserID, tokenData.Token)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员删除用户令牌：%s", userData.UserID)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
