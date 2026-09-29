package admin_website

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAllPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllPermissionsLogic {
	return &GetAllPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllPermissionsLogic) GetAllPermissions() (resp *types.GetAllPermissionsResp, err error) {
	res := make([]types.LabelValueRecord, 0, len(model.WebsitePermissionList()))
	for _, p := range model.WebsitePermissionList() {
		if p.Status == db.WebsitePolicyStatusBanned {
			continue
		}
		res = append(res, types.LabelValueRecord{
			Label: p.Name,
			Value: p.Sign,
		})
	}

	return &types.GetAllPermissionsResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetAllPermissionsData{
			Permissions: res,
		},
	}, nil
}
