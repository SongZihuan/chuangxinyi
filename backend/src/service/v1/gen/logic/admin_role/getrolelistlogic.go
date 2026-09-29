package admin_role

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

type GetRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleListLogic) GetRoleList(req *types.GetRoleListReq) (resp *types.RoleListResp, err error) {
	getAll := len(req.Name) == 0
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	roleListResp := make([]types.Role, 0, len(model.Roles()))
	for _, m := range model.Roles() {
		if getAll || strings.Contains(m.Name, req.Name) || strings.Contains(m.Describe, req.Name) {
			roleListResp = append(roleListResp, m.GetRole())
		}
	}

	sort.Slice(roleListResp, func(i, j int) bool {
		return roleListResp[i].ID > roleListResp[j].ID
	})

	count := int64(len(roleListResp))

	if start >= count {
		return &types.RoleListResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.RoleList{
				Role:  []types.Role{},
				Count: count,
			},
		}, nil
	}

	if end > count {
		end = count
	}

	return &types.RoleListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.RoleList{
			Role:  roleListResp[start:end],
			Count: count,
		},
	}, nil
}
