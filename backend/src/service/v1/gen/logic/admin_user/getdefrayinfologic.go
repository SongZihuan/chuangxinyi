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

type GetDefrayInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDefrayInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDefrayInfoLogic {
	return &GetDefrayInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDefrayInfoLogic) GetDefrayInfo(req *types.AdminGetDefrayInfoReq) (resp *types.AdminGetDefrayInfoResp, err error) {
	defrayModel := db.NewDefrayModel(mysql.MySQLConn)
	defray, err := defrayModel.FindByDefrayID(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.AdminGetDefrayInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "消费订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	defrayAt := int64(0)
	if defray.DefrayAt.Valid {
		defrayAt = defray.DefrayAt.Time.Unix()
	}

	returnAt := int64(0)
	if defray.ReturnAt.Valid {
		returnAt = defray.ReturnAt.Time.Unix()
	}

	optUser, err := action.GetUserEasy(l.ctx, defray.UserId.Int64, "")
	if errors.Is(err, action.UserEasyNotFound) {
		return &types.AdminGetDefrayInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	owner := types.UserEasy{}
	if defray.OwnerId.Valid {
		owner, err = action.GetUserEasy(l.ctx, defray.OwnerId.Int64, "")
		if errors.Is(err, action.UserEasyNotFound) {
			return &types.AdminGetDefrayInfoResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	return &types.AdminGetDefrayInfoResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetDefrayInfoData{
			Defray: types.AdminDefrayRecord{
				MustSelfDefray:     defray.MustSelfDefray,
				DefrayID:           defray.DefrayId,
				UserID:             defray.UserId.Int64,
				WalletID:           defray.WalletId.Int64,
				HasOwner:           defray.OwnerId.Valid,
				OwnerID:            defray.OwnerId.Int64,
				Owner:              owner,
				User:               optUser,
				Subject:            defray.Subject,
				Price:              defray.Price,
				RealPrice:          defray.RealPrice.Int64,
				UnitPrice:          defray.UnitPrice,
				Quantity:           defray.Quantity,
				Describe:           defray.Describe,
				SupplierID:         defray.SupplierId,
				Supplier:           defray.Supplier,
				Balance:            defray.Balance.Int64,
				InvitePre:          defray.InvitePre,
				DistributionLevel1: defray.DistributionLevel1,
				DistributionLevel2: defray.DistributionLevel2,
				DistributionLevel3: defray.DistributionLevel3,
				CanWithdraw:        defray.CanWithdraw,
				CouponsID:          defray.CouponsId.Int64,
				Remark:             defray.Remark,
				DefrayStatus:       defray.Status,
				CreateAt:           defray.CreateAt.Unix(),
				DefrayAt:           defrayAt,
				ReturnAt:           returnAt,
			},
		},
	}, nil
}
