package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type DeleteOneUserAllWebsiteTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOneUserAllWebsiteTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOneUserAllWebsiteTokenLogic {
	return &DeleteOneUserAllWebsiteTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOneUserAllWebsiteTokenLogic) DeleteOneUserAllWebsiteToken(req *types.AdminDeleteOneUserAllWebsiteTokenReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	}

	if req.WebID == 0 {
		err = jwt.DeleteAllWebsiteUserToken(l.ctx, srcUser.Uid)
		if err != nil {
			return nil, respmsg.RedisSystemError.WarpQuick(err)
		}
	} else {
		err = jwt.DeleteWebsiteUserToken(l.ctx, srcUser.Uid, req.WebID)
		if err != nil {
			return nil, respmsg.RedisSystemError.WarpQuick(err)
		}
	}

	audit.NewAdminAudit(user.Id, "管理员撤销用户（%s）的所有授权站点令牌", srcUser.Uid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
