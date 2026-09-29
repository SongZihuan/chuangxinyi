package admin_website_path

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

type GetPathListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPathListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPathListLogic {
	return &GetPathListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPathListLogic) GetPathList(req *types.GetWebPathListReq) (resp *types.AdminWebsitePathResp, err error) {
	getAll := len(req.Name) == 0
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	pathListResp := make([]types.WebsiteUrlPath, 0, len(model.WebsiteUrlPathMap()))
	for _, m := range model.WebsiteUrlPathMap() {
		if getAll || strings.Contains(m.Path, req.Name) || strings.Contains(m.Describe, req.Name) {
			pathListResp = append(pathListResp, m.GetWebsiteUrlPathType())
		}
	}

	sort.Slice(pathListResp, func(i, j int) bool {
		return pathListResp[i].ID > pathListResp[j].ID
	})

	count := int64(len(pathListResp))

	if start >= count {
		return &types.AdminWebsitePathResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.AdminWebsitePathData{
				Path:  []types.WebsiteUrlPath{},
				Count: count,
			},
		}, nil
	}

	if end > count {
		end = count
	}

	return &types.AdminWebsitePathResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminWebsitePathData{
			Path:  pathListResp[start:end],
			Count: count,
		},
	}, nil
}
