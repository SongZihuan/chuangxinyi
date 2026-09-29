package admin_path

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAllSubPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllSubPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllSubPermissionsLogic {
	return &GetAllSubPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllSubPermissionsLogic) GetAllSubPermissions() (resp *types.GetAllSubPermissionsResp, err error) {
	res := make([]types.LabelValueRecord, 0, len(jwt.UserSubTokenStringList))
	for _, p := range jwt.UserSubTokenStringList {
		label, ok := jwt.UserSubTokenStringChineseMap[p]
		if !ok {
			continue
		}

		res = append(res, types.LabelValueRecord{
			Label: label,
			Value: p,
		})
	}

	return &types.GetAllSubPermissionsResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetAllSubPermissionsData{
			Permissions: res,
		},
	}, nil
}
