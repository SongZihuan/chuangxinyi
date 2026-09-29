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

type GetUserBackListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserBackListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserBackListLogic {
	return &GetUserBackListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserBackListLogic) GetUserBackList(req *types.AdminGetBackListReq) (resp *types.AdminGetBackListResp, err error) {
	var backList []db.Back
	var count int64

	backModel := db.NewBackModel(mysql.MySQLConn)
	if req.ID == 0 && len(req.UID) == 0 {
		backList, err = backModel.GetList(l.ctx, 0, req.SupplierID, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = backModel.GetCount(l.ctx, 0, req.SupplierID, req.Src, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		user, err := GetUser(l.ctx, req.ID, req.UID, true)
		if errors.Is(err, UserNotFound) {
			return &types.AdminGetBackListResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		backList, err = backModel.GetList(l.ctx, user.WalletId, req.SupplierID, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = backModel.GetCount(l.ctx, user.WalletId, req.SupplierID, req.Src, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	userMap := make(map[int64]types.UserEasy, len(backList))

	respList := make([]types.AdminBackRecord, 0, len(backList))
	for _, d := range backList {
		optUser, ok := userMap[d.UserId]
		if !ok {
			optUser, err = action.GetUserEasy(l.ctx, d.UserId, "")
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			userMap[d.UserId] = optUser
		}

		respList = append(respList, types.AdminBackRecord{
			WalletID:    d.WalletId,
			UserID:      d.UserId,
			User:        optUser,
			Subject:     d.Subject,
			BackID:      d.BackId,
			Get:         d.Get,
			Balance:     d.Balance,
			CanWithdraw: d.CanWithdraw,
			SupplierID:  d.SupplierId,
			Supplier:    d.Supplier,
			Remark:      d.Remark,
			CreateAt:    d.CreateAt.Unix(),
		})
	}

	return &types.AdminGetBackListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetBackListData{
			Count: count,
			Back:  respList,
		},
	}, nil
}
