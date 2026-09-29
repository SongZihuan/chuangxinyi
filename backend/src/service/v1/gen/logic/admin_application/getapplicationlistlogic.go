package admin_application

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"sort"
	"strings"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetApplicationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplicationListLogic {
	return &GetApplicationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApplicationListLogic) GetApplicationList(req *types.GetApplicationListReq) (resp *types.AdminApplicationResp, err error) {
	getAll := len(req.Name) == 0
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	applicationResp := make([]types.AdminApplication, 0, len(model.ApplicationList()))
	for _, m := range model.ApplicationList() {
		if getAll || strings.Contains(m.Name, req.Name) || strings.Contains(m.Describe, req.Name) {
			applicationResp = append(applicationResp, m.GetAdminApplicationType())
		}
	}

	sort.Slice(applicationResp, func(i, j int) bool {
		return applicationResp[i].Sort < applicationResp[j].Sort
	})

	count := int64(len(applicationResp))

	if start >= count {
		return &types.AdminApplicationResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.AdminApplicationData{
				Application: []types.AdminApplication{},
				Count:       count,
			},
		}, nil
	}

	if end > count {
		end = count
	}
	return &types.AdminApplicationResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminApplicationData{
			Application: applicationResp[start:end],
			Count:       count,
		},
	}, nil
}
