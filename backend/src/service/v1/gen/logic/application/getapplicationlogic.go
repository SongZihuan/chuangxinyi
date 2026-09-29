package application

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplicationLogic {
	return &GetApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApplicationLogic) GetApplication() (resp *types.ApplicationResp, err error) {
	res := make([]types.Application, 0, len(model.ApplicationList()))

	for _, a := range model.ApplicationList() {
		if a.Status == db.ApplicationStatusBanned {
			continue
		}

		web := action.GetWebsite(a.WebID)
		if web.ID == warp.UnknownWebsite {
			continue
		}

		res = append(res, a.GetApplicationType())
	}

	return &types.ApplicationResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.ApplicationData{
			Application: res,
		},
	}, nil
}
