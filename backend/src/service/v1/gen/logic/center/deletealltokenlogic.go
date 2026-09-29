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

type DeleteAllTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAllTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAllTokenLogic {
	return &DeleteAllTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAllTokenLogic) DeleteAllToken() (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	t, ok := l.ctx.Value("X-Token").(string)
	if !ok || len(t) == 0 {
		t = ""
	}

	err = jwt.DeleteAllUserToken(l.ctx, user.Uid, t)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	err = jwt.DeleteAllUserSonToken(l.ctx, user.Uid)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户退出全部其他登录成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
