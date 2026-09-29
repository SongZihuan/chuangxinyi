package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/selfpay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
)

type SelfpayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelfpayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelfpayLogic {
	return &SelfpayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelfpayLogic) Selfpay(req *types.NewPayReq) (resp *types.SelfPayResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	payID, err := selfpay.NewSelfPay(l.ctx, req.Cny, user, req.PayWay, req.CouponsID)
	if errors.Is(err, selfpay.BadCNY) {
		return &types.SelfPayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "错误的支付金额"),
		}, nil
	} else if errors.Is(err, selfpay.BadPayWay) {
		return &types.SelfPayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayWay, "错误的支付方式"),
		}, nil
	} else if err != nil {
		return &types.SelfPayResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.SelfPayFail, errors.WarpQuick(err), "自支付失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户申请线下人工充值")

	return &types.SelfPayResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SelfPayData{
			ID: payID,
		},
	}, nil
}
