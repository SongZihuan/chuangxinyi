package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetDefrayOwnerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDefrayOwnerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDefrayOwnerListLogic {
	return &GetDefrayOwnerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDefrayOwnerListLogic) GetDefrayOwnerList(req *types.GetDefrayOwnerListReq) (resp *types.GetDefrayOwnerListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	if web.ID != warp.UserCenterWebsite && web.ID != req.SupplierID {
		return &types.GetDefrayOwnerListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "外站不允许操作"),
		}, nil
	}

	defrayModel := db.NewDefrayModel(mysql.MySQLConn)
	dList, err := defrayModel.GetListWithOwnerID(l.ctx, user.Id, req.Status, req.Src, req.SupplierID, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := defrayModel.GetCountWithOwnerID(l.ctx, user.Id, req.Status, req.Src, req.SupplierID, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userMap := make(map[int64]types.UserLessEasy, len(dList))

	respList := make([]types.OwnerDefrayRecord, 0, len(dList))
	for _, d := range dList {
		optUser, ok := userMap[d.UserId.Int64]
		if !ok {
			optUser, err = action.GetUserLessEasy(l.ctx, d.UserId.Int64, "")
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			userMap[d.UserId.Int64] = optUser
		}

		owner := types.UserEasy{}
		if d.OwnerId.Valid {
			owner, err = action.GetUserEasy(l.ctx, d.OwnerId.Int64, "")
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		}

		defrayAt := int64(0)
		if d.DefrayAt.Valid {
			defrayAt = d.DefrayAt.Time.Unix()
		}

		returnAt := int64(0)
		if d.ReturnAt.Valid {
			returnAt = d.ReturnAt.Time.Unix()
		}

		lastReturnAt := int64(0)
		if d.LastReturnAt.Valid {
			lastReturnAt = d.LastReturnAt.Time.Unix()
		}

		respList = append(respList, types.OwnerDefrayRecord{
			MustSelfDefray:     d.MustSelfDefray,
			User:               optUser,
			Owner:              owner,
			DefrayID:           d.DefrayId,
			Subject:            d.Subject,
			Price:              d.Price,
			RealPrice:          d.RealPrice.Int64,
			UnitPrice:          d.UnitPrice,
			Quantity:           d.Quantity,
			Describe:           d.Describe,
			Supplier:           d.Supplier,
			Balance:            d.Balance.Int64,
			InvitePre:          d.InvitePre,
			DistributionLevel1: d.DistributionLevel1,
			DistributionLevel2: d.DistributionLevel2,
			DistributionLevel3: d.DistributionLevel3,
			CanWithdraw:        d.CanWithdraw,
			HasCoupons:         d.CouponsId.Valid,
			ReturnDayLimit:     d.ReturnDayLimit,
			DefrayStatus:       d.Status,
			CreateAt:           d.CreateAt.Unix(),
			DefrayAt:           defrayAt,
			ReturnAt:           returnAt,
			LastReturnAt:       lastReturnAt,
		})
	}

	return &types.GetDefrayOwnerListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetDefrayOwnerListData{
			Count:  count,
			Defray: respList,
		},
	}, nil
}
