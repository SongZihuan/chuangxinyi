package defray

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type QueryDefrayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryDefrayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDefrayLogic {
	return &QueryDefrayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDefrayLogic) QueryDefray(req *types.QueryDefrayReq) (resp *types.QueryDefrayResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	defrayModel := db.NewDefrayModel(mysql.MySQLConn)
	d, err := defrayModel.FindByDefrayID(l.ctx, req.TradeID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.QueryDefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "支付订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if d.SupplierId != web.ID {
		return &types.QueryDefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "支付订单不属于该外站"),
		}, nil
	}

	if d.Status == db.DefrayWait {
		return &types.QueryDefrayResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryDefrayData{
				Status:   "WAIT",
				PayAt:    0,
				ReturnAt: 0,
			},
		}, nil
	} else if d.Status == db.DefraySuccess || d.Status == db.DefrayWaitReturn {
		userModel := db.NewUserModel(mysql.MySQLConn)
		payer, err := userModel.FindOneByIDWithoutDelete(l.ctx, d.UserId.Int64)
		if errors.Is(err, db.ErrNotFound) {
			return &types.QueryDefrayResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "付款人找不到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		if !d.DefrayAt.Valid {
			return &types.QueryDefrayResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDefrayAt, "付款时间错误"),
			}, nil
		}

		return &types.QueryDefrayResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryDefrayData{
				Status:   "SUCCESS",
				PayerID:  payer.Uid,
				PayAt:    d.DefrayAt.Time.Unix(),
				ReturnAt: 0,
			},
		}, nil
	} else if d.Status == db.DefrayReturn {
		userModel := db.NewUserModel(mysql.MySQLConn)
		payer, err := userModel.FindOneByIDWithoutDelete(l.ctx, d.UserId.Int64)
		if errors.Is(err, db.ErrNotFound) {
			return &types.QueryDefrayResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "付款人找不到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		if !d.DefrayAt.Valid {
			return &types.QueryDefrayResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDefrayAt, "付款时间错误"),
			}, nil
		}

		if !d.ReturnAt.Valid {
			return &types.QueryDefrayResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDefrayAt, "付款时间错误"),
			}, nil
		}

		return &types.QueryDefrayResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryDefrayData{
				Status:       "RETURN",
				PayerID:      payer.Uid,
				PayAt:        d.DefrayAt.Time.Unix(),
				ReturnAt:     d.ReturnAt.Time.Unix(),
				ReturnReason: d.ReturnReason.String,
			},
		}, nil
	} else {
		return &types.QueryDefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadDefrayStatus, "错误的支付状态"),
		}, nil
	}
}
