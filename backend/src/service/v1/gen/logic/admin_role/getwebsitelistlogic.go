package admin_role

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

func (l *GetWebsiteListLogic) GetWebsiteList(req *types.PageReq) (resp *types.RoleGetWebsiteListResp, err error) {
	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	if web.ID != warp.UserCenterWebsite {
		return &types.RoleGetWebsiteListResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.RoleGetWebsiteListData{
				Website: []types.WebsiteEasy{web.GetWebsiteEasyType()},
				Count:   1,
			},
		}, nil
	}

	count := int64(len(model.Websites()))

	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= count {
		return &types.RoleGetWebsiteListResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.RoleGetWebsiteListData{
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
		res = append(res, w.GetWebsiteEasyType())
	}

	return &types.RoleGetWebsiteListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.RoleGetWebsiteListData{
			Website: res[start:end],
			Count:   count,
		},
	}, nil
}
