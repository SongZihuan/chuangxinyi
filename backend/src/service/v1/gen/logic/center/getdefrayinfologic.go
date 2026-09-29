package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/defray"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"time"

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

func (l *GetDefrayInfoLogic) GetDefrayInfo(req *types.GetDefrayInfoReq) (resp *types.GetDefrayInfoResp, err error) {
	data, expire, err := defray.ParserDefray(req.Token)
	if err != nil {
		return &types.GetDefrayInfoResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.DefrayNotFound, errors.WarpQuick(err), "消费Token解析失败"),
		}, nil
	}

	if time.Now().After(expire) {
		return &types.GetDefrayInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayTimeout, "订单超时"),
		}, nil
	}

	supplier := action.GetWebsite(data.SupplierID)
	if supplier.Status == db.WebsiteStatusBanned {
		return &types.GetDefrayInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "供应商已被封禁"),
		}, nil
	}

	var owner types.UserLessEasy
	hasOwner := false
	if data.OwnerID != 0 {
		hasOwner = true
		owner, err = action.GetUserLessEasy(l.ctx, data.OwnerID, "")
		if err != nil {
			return &types.GetDefrayInfoResp{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.DefrayNotFound, errors.WarpQuick(err), "获取owner失败"),
			}, nil
		}
	}

	return &types.GetDefrayInfoResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetDefrayInfoData{
			Info: types.DefrayData{
				MustSelfDefray: data.MustSelfDefray,
				HasOwner:       hasOwner,
				Owner:          owner,
				Subject:        data.Subject,
				Price:          data.Price,
				Quantity:       data.Quantity,
				UnitPrice:      data.UnitPrice,
				Supplier:       supplier.Name,
				Describe:       data.Describe,
			},
		},
	}, nil
}
