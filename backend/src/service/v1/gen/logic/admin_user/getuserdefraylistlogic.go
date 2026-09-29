package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserDefrayListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserDefrayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDefrayListLogic {
	return &GetUserDefrayListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserDefrayListLogic) GetUserDefrayList(req *types.AdminGetDefrayListReq) (resp *types.AdminGetDefrayListResp, err error) {
	var defrayList []db.Defray
	var count int64

	defrayModel := db.NewDefrayModel(mysql.MySQLConn)
	if req.ID == 0 && len(req.UID) == 0 {
		defrayList, err = defrayModel.GetList(l.ctx, 0, req.Status, req.Src, req.SupplierID, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = defrayModel.GetCount(l.ctx, 0, req.Status, req.Src, req.SupplierID, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		user, err := GetUser(l.ctx, req.ID, req.UID, true)
		if errors.Is(err, UserNotFound) {
			return &types.AdminGetDefrayListResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		defrayList, err = defrayModel.GetList(l.ctx, user.WalletId, req.Status, req.Src, req.SupplierID, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = defrayModel.GetCount(l.ctx, user.WalletId, req.Status, req.Src, req.SupplierID, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	userMap := make(map[int64]types.UserEasy, len(defrayList))

	respList := make([]types.AdminDefrayRecord, 0, len(defrayList))
	for _, d := range defrayList {
		optUser, ok := userMap[d.UserId.Int64]
		if !ok {
			optUser, err = action.GetUserEasy(l.ctx, d.UserId.Int64, "")
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

		respList = append(respList, types.AdminDefrayRecord{
			MustSelfDefray:     d.MustSelfDefray,
			DefrayID:           d.DefrayId,
			UserID:             d.UserId.Int64,
			WalletID:           d.WalletId.Int64,
			HasOwner:           d.OwnerId.Valid,
			OwnerID:            d.OwnerId.Int64,
			Owner:              owner,
			User:               optUser,
			Subject:            d.Subject,
			Price:              d.Price,
			RealPrice:          d.RealPrice.Int64,
			UnitPrice:          d.UnitPrice,
			Quantity:           d.Quantity,
			Describe:           d.Describe,
			Supplier:           d.Supplier,
			SupplierID:         d.SupplierId,
			Balance:            d.Balance.Int64,
			InvitePre:          d.InvitePre,
			DistributionLevel1: d.DistributionLevel1,
			DistributionLevel2: d.DistributionLevel2,
			DistributionLevel3: d.DistributionLevel3,
			CanWithdraw:        d.CanWithdraw,
			CouponsID:          d.CouponsId.Int64,
			DefrayStatus:       d.Status,
			Remark:             d.Remark,
			CreateAt:           d.CreateAt.Unix(),
			DefrayAt:           defrayAt,
			ReturnAt:           returnAt,
		})
	}

	return &types.AdminGetDefrayListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetDefrayListData{
			Count:  count,
			Defray: respList,
		},
	}, nil
}
