package admin_menu

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"sort"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuLogic) GetMenu() (resp *types.AdminMenuResp, err error) {
	res := make([]types.Menu, 0, len(model.Menus()))
	for _, m := range model.Menus() {
		res = append(res, m.GetMenuType())
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Sort < res[j].Sort
	})

	return &types.AdminMenuResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminMenuData{
			Menu: res,
		},
	}, nil
}
