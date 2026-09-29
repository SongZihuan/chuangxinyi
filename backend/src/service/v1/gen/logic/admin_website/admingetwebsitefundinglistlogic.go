package admin_website

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AdminGetWebsiteFundingListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetWebsiteFundingListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetWebsiteFundingListLogic {
	return &AdminGetWebsiteFundingListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetWebsiteFundingListLogic) AdminGetWebsiteFundingList(req *types.AdminGetWebsiteFundingListReq) (resp *types.AdminGetWebsiteFundingListResp, err error) {
	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	if web.ID != warp.UserCenterWebsite && web.ID != req.WebsiteID {
		return &types.AdminGetWebsiteFundingListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "外站不允许操作"),
		}, nil
	}

	srcWeb := action.GetWebsite(req.WebsiteID)
	if srcWeb.ID == warp.UnknownWebsite {
		return &types.AdminGetWebsiteFundingListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotFound, "外站未找到"),
		}, nil
	}

	websiteFundingModel := db.NewWebsiteFundingModel(mysql.MySQLConn)
	fundingList, err := websiteFundingModel.GetList(l.ctx, req.WebsiteID, req.Type, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := websiteFundingModel.GetCount(l.ctx, req.WebsiteID, req.Type, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	listResp := make([]types.WebsiteFunding, 0, len(fundingList))
	for _, f := range fundingList {
		listResp = append(listResp, types.WebsiteFunding{
			WebID:       srcWeb.ID,
			WebName:     srcWeb.Name,
			Type:        f.Type,
			FundingId:   f.FundingId,
			Profit:      f.Profit,
			Expenditure: f.Expenditure,
			Delta:       0 + f.Profit - f.Expenditure,
			Remark:      f.Remark,
			PayAt:       f.PayAt.Unix(),
		})
	}

	return &types.AdminGetWebsiteFundingListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetWebsiteFundingListData{
			Count:   count,
			Funding: listResp,
		},
	}, nil
}
