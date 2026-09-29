package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"

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

func (l *GetCouponsListLogic) GetCouponsList(req *types.AdminGetCouponsListReq) (resp *types.AdminGetCouponsListResp, err error) {
	var couponsList []db.Coupons
	var count int64

	couponsModel := db.NewCouponsModel(mysql.MySQLConn)
	if req.ID == 0 && len(req.UID) == 0 {
		couponsList, err = couponsModel.GetList(l.ctx, 0, req.Type, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = couponsModel.GetCount(l.ctx, 0, req.Type, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		user, err := GetUser(l.ctx, req.ID, req.UID, true)
		if errors.Is(err, UserNotFound) {
			return &types.AdminGetCouponsListResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		couponsList, err = couponsModel.GetList(l.ctx, user.Id, req.Type, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = couponsModel.GetCount(l.ctx, user.Id, req.Type, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.AdminCoupons, 0, len(couponsList))
	for _, c := range couponsList {
		var content map[string]interface{}
		err = utils.JsonUnmarshal([]byte(c.Content), &content)
		if err != nil {
			continue
		}

		respList = append(respList, types.AdminCoupons{
			ID:      c.Id,
			UserID:  c.UserId,
			Name:    c.Name,
			Type:    c.Type,
			Content: content,
		})
	}

	return &types.AdminGetCouponsListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetCouponsListData{
			Count:   count,
			Coupons: respList,
		},
	}, nil
}
