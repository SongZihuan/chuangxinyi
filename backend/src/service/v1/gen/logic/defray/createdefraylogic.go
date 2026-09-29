package defray

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/defray"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
)

type CreateDefrayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDefrayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDefrayLogic {
	return &CreateDefrayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDefrayLogic) CreateDefray(req *types.CreateDefrayReq) (resp *types.CreateDefrayResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	var owner *db.User = nil
	ownerID := int64(0)
	if len(req.OwnerID) != 0 {
		userModel := db.NewUserModel(mysql.MySQLConn)
		owner, err = userModel.FindOneByUidWithoutDelete(l.ctx, req.OwnerID)
		if errors.Is(err, db.ErrNotFound) {
			return &types.CreateDefrayResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		ownerID = owner.Id
	}

	if req.MustSelfDefray && ownerID == 0 {
		return &types.CreateDefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "必须支付的用户未找到"),
		}, nil
	}

	token, defrayID, err := defray.NewDefray(l.ctx, jwt.DefrayTokenData{
		MustSelfDefray:     req.MustSelfDefray,
		OwnerID:            ownerID,
		Subject:            req.Subject,
		Price:              req.Price,
		Quantity:           req.Quantity,
		UnitPrice:          req.UnitPrice,
		SupplierID:         web.ID,
		Describe:           req.Describe,
		ReturnURL:          req.ReturnURL,
		InvitePre:          req.InvitePre,
		DistributionLevel1: req.DistributionLevel1,
		DistributionLevel2: req.DistributionLevel2,
		DistributionLevel3: req.DistributionLevel3,
		CanWithdraw:        req.CanWithdraw,
		TimeExpire:         req.ExpireTime,
		ReturnDayLimit:     req.ReturnDayLimit,
	}, owner)
	if err != nil {
		return &types.CreateDefrayResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateDefrayFail, errors.WarpQuick(err), "创建支付订单失败"),
		}, nil
	}

	return &types.CreateDefrayResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CreateDefrayData{
			Token:   token,
			TradeID: defrayID,
		},
	}, nil
}
