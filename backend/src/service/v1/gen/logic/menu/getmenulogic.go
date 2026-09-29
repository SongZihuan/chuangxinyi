package menu

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

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

func (l *GetMenuLogic) GetMenu() (resp *types.MenuResp, err error) {
	role, ok := l.ctx.Value("X-Token-Role").(warp.Role)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Role")
	}

	subType, ok := l.ctx.Value("X-Token-Type").(int)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Type")
	}

	subTypePermission, ok := jwt.UserSubTokenPermissionMap[subType]
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Type")
	}

	res := make([]types.RoleMenu, 0, len(model.Menus()))
	for _, m := range role.Menus {
		srcMenu, ok := (model.Menus())[m.ID]
		if !ok || !permission.CheckPermissionInt64(srcMenu.SubPolicyPermission, subTypePermission) {
			continue
		}
		res = append(res, m)
	}

	return &types.MenuResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.MenuData{
			Menu: res,
		},
	}, nil
}
