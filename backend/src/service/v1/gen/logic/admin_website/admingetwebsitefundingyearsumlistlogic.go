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

type AdminGetWebsiteFundingYearSumListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetWebsiteFundingYearSumListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetWebsiteFundingYearSumListLogic {
	return &AdminGetWebsiteFundingYearSumListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetWebsiteFundingYearSumListLogic) AdminGetWebsiteFundingYearSumList(req *types.AdminGetWebsiteFundingListYearSumReq) (resp *types.AdminGetWebsiteFundingYearSumListResp, err error) {
	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	if web.ID != warp.UserCenterWebsite && web.ID != req.WebsiteID {
		return &types.AdminGetWebsiteFundingYearSumListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "外站不允许操作"),
		}, nil
	}

	srcWeb := action.GetWebsite(req.WebsiteID)
	if srcWeb.ID == warp.UnknownWebsite {
		return &types.AdminGetWebsiteFundingYearSumListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotFound, "外站未找到"),
		}, nil
	}

	websiteFundingModel := db.NewWebsiteFundingModel(mysql.MySQLConn)
	fundingList, err := websiteFundingModel.GetYearSum(l.ctx, req.WebsiteID, req.Year)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	listResp := make([]types.WebsiteFundingYearSum, 0, len(fundingList))
	for _, f := range fundingList {
		listResp = append(listResp, types.WebsiteFundingYearSum{
			WebID:       srcWeb.ID,
			WebName:     srcWeb.Name,
			Profit:      f.Profit,
			Expenditure: f.Expenditure,
			Delta:       f.Delta,
			Month:       f.Month,
			Day:         f.Day,
		})
	}

	return &types.AdminGetWebsiteFundingYearSumListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetWebsiteFundingYearSumListData{
			Funding: listResp,
		},
	}, nil
}
