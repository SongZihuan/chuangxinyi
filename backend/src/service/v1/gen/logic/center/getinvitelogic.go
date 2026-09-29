package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetInviteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInviteLogic {
	return &GetInviteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInviteLogic) GetInvite() (resp *types.GetInviteResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !user.InviteId.Valid {
		return &types.GetInviteResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetInviteData{
				HasInvite: false,
			},
		}, nil
	}

	invite, err := action.GetUserMoreEasy(l.ctx, user.InviteId.Int64, "")
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.GetInviteResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetInviteData{
			HasInvite: true,
			Invite:    invite,
		},
	}, nil
}
