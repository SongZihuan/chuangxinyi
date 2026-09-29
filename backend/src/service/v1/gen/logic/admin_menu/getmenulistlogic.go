package admin_menu

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

type GetMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuListLogic {
	return &GetMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuListLogic) GetMenuList(req *types.GetMenuListReq) (resp *types.AdminMenuResp, err error) {
	getAll := len(req.Name) == 0
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	menuListResp := make([]types.Menu, 0, len(model.Menus()))
	for _, m := range model.Menus() {
		if getAll || strings.Contains(m.Name, req.Name) || strings.Contains(m.Describe, req.Name) || strings.Contains(m.Title, req.Name) {
			menuListResp = append(menuListResp, m.GetMenuType())
		}
	}

	sort.Slice(menuListResp, func(i, j int) bool {
		return menuListResp[i].Sort < menuListResp[j].Sort
	})

	count := int64(len(menuListResp))

	if start >= count {
		return &types.AdminMenuResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.AdminMenuData{
				Menu:  []types.Menu{},
				Count: count,
			},
		}, nil
	}

	if end > count {
		end = count
	}

	return &types.AdminMenuResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminMenuData{
			Menu:  menuListResp[start:end],
			Count: count,
		},
	}, nil
}
