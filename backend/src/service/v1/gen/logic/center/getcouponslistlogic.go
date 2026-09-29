package center

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

type GetCouponsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCouponsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCouponsListLogic {
	return &GetCouponsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCouponsListLogic) GetCouponsList(req *types.GetCouponsListReq) (resp *types.GetCouponsListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	couponsModel := db.NewCouponsModel(mysql.MySQLConn)
	couponsList, err := couponsModel.GetList(l.ctx, user.Id, req.Type, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := couponsModel.GetCount(l.ctx, user.Id, req.Type, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.Coupons, 0, len(couponsList))
	for _, c := range couponsList {
		var content map[string]interface{}
		err = utils.JsonUnmarshal([]byte(c.Content), &content)
		if err != nil {
			continue
		}

		respList = append(respList, types.Coupons{
			ID:      c.Id,
			Name:    c.Name,
			Type:    c.Type,
			Content: content,
		})
	}

	return &types.GetCouponsListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetCouponsListData{
			Count:   count,
			Coupons: respList,
		},
	}, nil
}
