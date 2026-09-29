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

type GetWithdrawListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWithdrawListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWithdrawListLogic {
	return &GetWithdrawListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWithdrawListLogic) GetWithdrawList(req *types.GetWithdrawListReq) (resp *types.GetWithdrawListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	dList, err := withdrawModel.GetList(l.ctx, user.WalletId, req.Status, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := withdrawModel.GetCount(l.ctx, user.WalletId, req.Status, req.Src, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userMap := make(map[int64]types.UserEasy, len(dList))

	respList := make([]types.WithdrawRecord, 0, len(dList))
	for _, d := range dList {
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

		payAt := int64(0)
		if d.PayAt.Valid {
			payAt = d.PayAt.Time.Unix()
		}

		respList = append(respList, types.WithdrawRecord{
			User:              optUser,
			WithdrawID:        d.WithdrawId,
			WithdrawWay:       d.WithdrawWay,
			Name:              d.Name,
			AlipayLoginId:     d.AlipayLoginId.String,
			WechatpayNickName: d.WechatpayNickname.String,
			Cny:               d.Cny,
			Balance:           d.Balance.Int64,
			Status:            d.Status,
			WithdrawAt:        d.WithdrawAt.Unix(),
			PayAt:             payAt,
		})
	}

	return &types.GetWithdrawListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetWithdrawListData{
			Count:    count,
			Withdraw: respList,
		},
	}, nil
}
