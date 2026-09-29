package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetFatherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFatherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFatherLogic {
	return &GetFatherLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFatherLogic) GetFather() (resp *types.GetFatherResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !user.FatherId.Valid {
		return &types.GetFatherResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetFatherData{
				HasFather: false,
			},
		}, nil
	}

	father, err := action.GetUserEasy(l.ctx, user.FatherId.Int64, "")
	if errors.Is(err, action.UserEasyNotFound) {
		return &types.GetFatherResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetFatherData{
				HasFather: false,
			},
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.GetFatherResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetFatherData{
			HasFather: true,
			Father:    father,
		},
	}, nil
}
