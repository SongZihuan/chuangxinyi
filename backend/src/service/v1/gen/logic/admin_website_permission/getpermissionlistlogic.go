package admin_website_permission

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"sort"
	"strings"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetPermissionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPermissionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionListLogic {
	return &GetPermissionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPermissionListLogic) GetPermissionList(req *types.GetWebPermissionListReq) (resp *types.AdminWebsitePermissionResp, err error) {
	getAll := len(req.Name) == 0
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	permissionListResp := make([]types.WebsitePolicy, 0, len(model.WebsitePermissionList()))
	for _, m := range model.WebsitePermissionList() {
		if getAll || strings.Contains(m.Name, req.Name) || strings.Contains(m.Describe, req.Name) {
			permissionListResp = append(permissionListResp, m.GetWebsitePolicyType())
		}
	}

	sort.Slice(permissionListResp, func(i, j int) bool {
		return permissionListResp[i].Sort < permissionListResp[j].Sort
	})

	count := int64(len(permissionListResp))

	if start >= count {
		return &types.AdminWebsitePermissionResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.AdminWebsitePermissionData{
				Permission: []types.WebsitePolicy{},
				Count:      count,
			},
		}, nil
	}

	if end > count {
		end = count
	}
	return &types.AdminWebsitePermissionResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminWebsitePermissionData{
			Permission: permissionListResp[start:end],
			Count:      count,
		},
	}, nil
}
