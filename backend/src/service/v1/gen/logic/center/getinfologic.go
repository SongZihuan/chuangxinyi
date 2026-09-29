package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInfoLogic {
	return &GetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInfoLogic) GetInfo() (resp *types.GetInfoResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	tokenType, ok := l.ctx.Value("X-Token-Type").(int)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Type")
	}

	data, err := utils.GetUserInfo(l.ctx, user, tokenType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.GetInfoResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetInfoData{
			User:    data.User,
			Info:    data.Info,
			Data:    data.Data,
			Balance: data.Balance,
			Address: data.Address,
			Title:   data.Title,
			Role:    data.Role,
		},
	}, nil
}
