package center

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

func (l *GetWalletRecordListLogic) GetWalletRecordList(req *types.GetWalletRecordListReq) (resp *types.GetWalletRecordListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	recordModel := db.NewWalletRecordModel(mysql.MySQLConn)
	recordList, err := recordModel.GetList(l.ctx, user.WalletId, req.Type, req.FundingID, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := recordModel.GetCount(l.ctx, user.WalletId, req.Type, req.FundingID, req.Src, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userMap := make(map[int64]types.UserEasy, len(recordList))

	respList := make([]types.WalletRecord, 0, len(recordList))
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

		respList = append(respList, types.WalletRecord{
			ID:   p.Id,
			User: optUser,

			Type:      p.Type,
			FundingId: p.FundingId,
			Reason:    p.Reason,

			Balance:      p.Balance,
			WaitBalance:  p.WaitBalance,
			Cny:          p.Cny,
			NotBilled:    p.NotBilled,
			Billed:       p.Billed,
			HasBilled:    p.HasBilled,
			Withdraw:     p.Withdraw,
			WaitWithdraw: p.WaitWithdraw,
			NotWithdraw:  p.NotWithdraw,
			HasWithdraw:  p.HasWithdraw,

			BeforeBalance:      p.BeforeBalance,
			BeforeWaitBalance:  p.BeforeWaitBalance,
			BeforeCny:          p.BeforeCny,
			BeforeNotBilled:    p.BeforeNotBilled,
			BeforeBilled:       p.BeforeBilled,
			BeforeHasBilled:    p.BeforeHasBilled,
			BeforeWithdraw:     p.BeforeWithdraw,
			BeforeWaitWithdraw: p.BeforeWaitWithdraw,
			BeforeNotWithdraw:  p.BeforeNotWithdraw,
			BeforeHasWithdraw:  p.BeforeHasWithdraw,

			CreateAt: p.CreateAt.Unix(),
		})
	}

	return &types.GetWalletRecordListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetWalletRecordListData{
			Count:  count,
			Record: respList,
		},
	}, nil
}
