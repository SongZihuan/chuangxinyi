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

type GetWalletRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWalletRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletRecordListLogic {
	return &GetWalletRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWalletRecordListLogic) GetWalletRecordList(req *types.AdminGetWalletRecordListReq) (resp *types.AdminGetWalletRecordListResp, err error) {
	if req.ID == 0 && len(req.UID) == 0 {
		return &types.AdminGetWalletRecordListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	}

	recordModel := db.NewWalletRecordModel(mysql.MySQLConn)
	user, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetWalletRecordListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	recordList, err := recordModel.GetList(l.ctx, user.WalletId, req.Type, req.FundingID, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := recordModel.GetCount(l.ctx, user.WalletId, req.Type, req.FundingID, req.Src, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userMap := make(map[int64]types.UserEasy, len(recordList))

	respList := make([]types.AdminWalletRecord, 0, len(recordList))
	for _, p := range recordList {
		optUser, ok := userMap[p.UserId]
		if !ok {
			optUser, err = action.GetUserEasy(l.ctx, p.UserId, "")
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			userMap[p.UserId] = optUser
		}

		respList = append(respList, types.AdminWalletRecord{
			ID:       p.Id,
			UserID:   p.UserId,
			User:     optUser,
			WalletID: p.WalletId,

			Type:      p.Type,
			FundingId: p.FundingId,
			Reason:    p.Reason,

			Balance:     p.Balance,
			Cny:         p.Cny,
			NotBilled:   p.NotBilled,
			Billed:      p.Billed,
			HasBilled:   p.HasBilled,
			Withdraw:    p.Withdraw,
			NotWithdraw: p.NotWithdraw,
			HasWithdraw: p.HasWithdraw,

			BeforeBalance:     p.BeforeBalance,
			BeforeCny:         p.BeforeCny,
			BeforeNotBilled:   p.BeforeNotBilled,
			BeforeBilled:      p.BeforeBilled,
			BeforeHasBilled:   p.BeforeHasBilled,
			BeforeWithdraw:    p.BeforeWithdraw,
			BeforeNotWithdraw: p.BeforeNotWithdraw,
			BeforeHasWithdraw: p.BeforeHasWithdraw,

			Remark:   p.Remark,
			CreateAt: p.CreateAt.Unix(),
		})
	}

	return &types.AdminGetWalletRecordListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetWalletRecordListData{
			Count:  count,
			Record: respList,
		},
	}, nil
}
