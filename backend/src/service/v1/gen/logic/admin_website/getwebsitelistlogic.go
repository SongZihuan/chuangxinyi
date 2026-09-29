package admin_website

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"strings"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetWebsiteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWebsiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWebsiteListLogic {
	return &GetWebsiteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWebsiteListLogic) GetWebsiteList(req *types.GetWebsiteListReq) (resp *types.GetWebsiteListResp, err error) {
	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	if web.ID == warp.UserCenterWebsite {
		getAll := len(req.Name) == 0
		start := (req.Page - 1) * req.PageSize
		end := start + req.PageSize

		websiteListResp := make([]types.Website, 0, len(model.WebsiteList()))
		for _, m := range model.WebsiteList() {
			if getAll || strings.Contains(m.Name, req.Name) || strings.Contains(m.Describe, req.Name) {
				websiteListResp = append(websiteListResp, m.GetWebsiteType())
			}
		}

		count := int64(len(websiteListResp))

		if start >= count {
			return &types.GetWebsiteListResp{
				Resp: respmsg.GetRespSuccess(l.ctx),
				Data: types.GetWebsiteListData{
					Website: []types.Website{},
					Count:   count,
				},
			}, nil
		}

		if end > count {
			end = count
		}

		return &types.GetWebsiteListResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetWebsiteListData{
				Website: websiteListResp[start:end],
				Count:   count,
			},
		}, nil
	} else {
		return &types.GetWebsiteListResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetWebsiteListData{
				Website: []types.Website{web.GetWebsiteType()},
				Count:   1,
			},
		}, nil
	}
}
