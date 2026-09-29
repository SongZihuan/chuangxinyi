package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

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

func (l *GetWebsiteListLogic) GetWebsiteList(req *types.PageReq) (resp *types.Oauth2GetWebsiteListResp, err error) {
	count := int64(len(model.Websites()))

	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= count {
		return &types.Oauth2GetWebsiteListResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.Oauth2GetWebsiteListData{
				Website: []types.WebsiteEasy{},
				Count:   count,
			},
		}, nil
	}

	if end > count {
		end = count
	}

	res := make([]types.WebsiteEasy, 0, count)
	for _, w := range model.WebsiteList() {
		if w.ID == warp.UnknownWebsite {
			continue
		}
		res = append(res, w.GetWebsiteEasyType())
	}

	return &types.Oauth2GetWebsiteListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.Oauth2GetWebsiteListData{
			Website: res[start:end],
			Count:   count,
		},
	}, nil
}
