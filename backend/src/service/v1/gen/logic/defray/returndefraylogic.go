package defray

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/defray"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type ReturnDefrayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReturnDefrayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReturnDefrayLogic {
	return &ReturnDefrayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReturnDefrayLogic) ReturnDefray(req *types.ReturnDefrayReq) (resp *types.RespEmpty, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	defrayModel := db.NewDefrayModel(mysql.MySQLConn)
	d, err := defrayModel.FindByDefrayID(l.ctx, req.TradeID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "支付订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if d.SupplierId != web.ID {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "支付订单不属于该外站"),
		}, nil
	}

	if !req.Must && (!d.LastReturnAt.Valid || d.LastReturnAt.Time.Before(time.Now()) || d.HasDistribution) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.ReturnTooLate, "退款太迟"),
		}, nil
	}

	err = defray.ReturnWebsite(l.ctx, d, req.Reason, req.Must)
	if errors.Is(err, defray.DoubleReturn) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DoubleReturn, "双重退款"),
		}, nil
	} else if errors.Is(err, defray.DefrayNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "支付订单未找到"),
		}, nil
	} else if errors.Is(err, defray.InsufficientQuota) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.ReturnInsufficientQuota, "退款发票额度不足"),
		}, nil
	} else if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.DefrayReturnFail, errors.WarpQuick(err), "订单退款失败"),
		}, nil
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
