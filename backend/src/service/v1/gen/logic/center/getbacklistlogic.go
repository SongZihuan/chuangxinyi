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

type GetBackListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBackListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBackListLogic {
	return &GetBackListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBackListLogic) GetBackList(req *types.GetBackListReq) (resp *types.GetBackListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	if web.ID != warp.UserCenterWebsite && web.ID != req.SupplierID {
		return &types.GetBackListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "外站不允许操作"),
		}, nil
	}

	backModel := db.NewBackModel(mysql.MySQLConn)

	dList, err := backModel.GetList(l.ctx, user.WalletId, req.SupplierID, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := backModel.GetCount(l.ctx, user.WalletId, req.SupplierID, req.Src, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userMap := make(map[int64]types.UserEasy, len(dList))

	respList := make([]types.BackRecord, 0, len(dList))
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

		respList = append(respList, types.BackRecord{
			User:        optUser,
			Subject:     d.Subject,
			BackID:      d.BackId,
			Get:         d.Get,
			Balance:     d.Balance,
			CanWithdraw: d.CanWithdraw,
			Supplier:    d.Supplier,
			CreateAt:    d.CreateAt.Unix(),
		})
	}

	return &types.GetBackListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetBackListData{
			Count: count,
			Back:  respList,
		},
	}, nil
}
