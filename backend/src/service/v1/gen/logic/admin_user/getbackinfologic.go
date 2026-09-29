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

type GetBackInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBackInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBackInfoLogic {
	return &GetBackInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBackInfoLogic) GetBackInfo(req *types.AdminGetBackInfoReq) (resp *types.AdminGetBackInfoResp, err error) {
	backModel := db.NewBackModel(mysql.MySQLConn)
	back, err := backModel.FindByBackID(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.AdminGetBackInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BackNotFound, "消费订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	optUser, err := action.GetUserEasy(l.ctx, back.UserId, "")
	if errors.Is(err, action.UserEasyNotFound) {
		return &types.AdminGetBackInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.AdminGetBackInfoResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetBackInfoData{
			Back: types.AdminBackRecord{
				WalletID:    back.WalletId,
				UserID:      back.UserId,
				User:        optUser,
				Subject:     back.Subject,
				BackID:      back.BackId,
				Get:         back.Get,
				Balance:     back.Balance,
				CanWithdraw: back.CanWithdraw,
				SupplierID:  back.SupplierId,
				Supplier:    back.Supplier,
				Remark:      back.Remark,
				CreateAt:    back.CreateAt.Unix(),
			},
		},
	}, nil
}
