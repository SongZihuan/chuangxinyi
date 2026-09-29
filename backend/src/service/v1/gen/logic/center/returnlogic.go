package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/defray"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type ReturnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReturnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReturnLogic {
	return &ReturnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReturnLogic) Return(req *types.ReturnReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	defrayModel := db.NewDefrayModel(mysql.MySQLConn)
	d, err := defrayModel.FindByDefrayID(l.ctx, req.TradeID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "消费订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if d.WalletId.Int64 != user.WalletId {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "消费订单并非此用户"),
		}, nil
	}

	if !d.LastReturnAt.Valid || d.LastReturnAt.Time.Before(time.Now()) || d.HasDistribution {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.ReturnTooLate, "退款太迟"),
		}, nil
	}

	err = defray.Return(l.ctx, d, req.Reason)
	if errors.Is(err, defray.DoubleReturn) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DoubleReturn, "二次退款"),
		}, nil
	} else if errors.Is(err, defray.DefrayNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "订单未找到"),
		}, nil
	} else if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.RefundFail, errors.WarpQuick(err), "退款失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户消费退款成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
