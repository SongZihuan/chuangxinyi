package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

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

func (l *DeleteAllTokenLogic) DeleteAllToken(req *types.DeleteAllOauth2TokenReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.WebID == 0 {
		err = jwt.DeleteAllLoginToken(l.ctx, user.Uid)
		if err != nil {
			return nil, respmsg.JWTError.WarpQuick(err)
		}

		audit.NewUserAudit(user.Id, "用户撤销所有对外访问权限")
	} else {
		web := action.GetWebsite(req.WebID)
		if web.ID == warp.UnknownWebsite || web.ID == warp.UserCenterWebsite {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotFound, "外站未找到"),
			}, nil
		}

		err = jwt.DeleteAllWebsiteLoginToken(l.ctx, user.Uid, req.WebID)
		if err != nil {
			return nil, respmsg.JWTError.WarpQuick(err)
		}

		audit.NewUserAudit(user.Id, "用户撤下站点（%s）所有访问权限", web.Name)
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
