package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/back"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddBackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddBackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBackLogic {
	return &AddBackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddBackLogic) AddBack(req *types.AdminAddBackReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	_, err = back.NewBack(l.ctx, req.Get, req.Reason, req.Subject, srcUser, req.CanWithdraw, web.ID)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.ChangeBalanceFail, errors.WarpQuick(err), "添加返现订单失败"),
		}, nil
	}

	audit.NewAdminAudit(user.Id, "管理员新增用户返现：%s", srcUser.Uid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
