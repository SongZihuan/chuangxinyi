package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetDiscountListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDiscountListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiscountListLogic {
	return &GetDiscountListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDiscountListLogic) GetDiscountList(req *types.AdminGetDiscountList) (resp *types.AdminGetDiscountListResp, err error) {
	discountModel := db.NewDiscountModel(mysql.MySQLConn)
	discountList, err := discountModel.GetList(l.ctx, req.Src, req.Page, req.PageSize, false)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := discountModel.GetCount(l.ctx, req.Src, false)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.AdminDiscount, 0, len(discountList))
	for _, d := range discountList {
		var quota map[string]interface{}
		err = utils.JsonUnmarshal([]byte(d.Quota), &quota)
		if err != nil {
			continue
		}

		respList = append(respList, types.AdminDiscount{
			ID:                d.Id,
			Name:              d.Name,
			Describe:          d.Describe,
			ShortDescribe:     d.ShortDescribe,
			Type:              d.Type,
			Quota:             quota,
			DayLimit:          d.DayLimit.Int64,
			MonthLimit:        d.MonthLimit.Int64,
			YearLimit:         d.YearLimit.Int64,
			Limit:             d.Limit.Int64,
			NeedVerify:        d.NeedVerify,
			NeedCompany:       d.NeedCompany,
			NeedUserOrigin:    d.NeedUserOrigin,
			NeedCompanyOrigin: d.NeedCompanyOrigin,
			NeedUserFace:      d.NeedUserFace,
			NeedCompanyFace:   d.NeedCompanyFace,
			Show:              d.Show,
			Remark:            d.Remark,
		})
	}

	return &types.AdminGetDiscountListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetDiscountListData{
			Count:    count,
			Discount: respList,
		},
	}, nil
}
