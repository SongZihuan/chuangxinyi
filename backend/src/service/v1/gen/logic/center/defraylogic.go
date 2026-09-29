package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/defray"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type DefrayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDefrayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DefrayLogic {
	return &DefrayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DefrayLogic) Defray(req *types.DefrayReq) (resp *types.DefrayResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	data, expire, err := defray.ParserDefray(req.Token)
	if err != nil {
		return &types.DefrayResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.DefrayNotFound, errors.WarpQuick(err), "解析支付Token失败"),
		}, nil
	}

	web := action.GetWebsite(data.SupplierID)
	if web.Status == db.WebsiteStatusBanned {
		return &types.DefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "供应商被封禁"),
		}, nil
	}

	if web.ID != warp.UserCenterWebsite {
		bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
		allow, err := bannedModel.CheckAllow(l.ctx, user.Id, web.ID, db.AllowLogin)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		} else if !allow {
			return &types.DefrayResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "用户禁止向该网站支付"),
			}, nil
		}
	}

	if time.Now().After(expire) {
		return &types.DefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayTimeout, "支付超时"),
		}, nil
	}

	d, dt, err := defray.Pay(l.ctx, data.TradeID, user, req.CouponsID, req.Token)
	switch true {
	case err == nil:
		// 啥也不做
	case errors.Is(err, defray.UserNotFount):
		return &types.DefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	case errors.Is(err, defray.DoubleDefray):
		return &types.DefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DoublePayment, "双重支付"),
		}, nil
	case errors.Is(err, defray.DefrayNotFound):
		return &types.DefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayNotFound, "支付订单未找到"),
		}, nil
	case errors.Is(err, defray.Insufficient):
		return &types.DefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DefrayInsufficient, "额度不足"),
			Data: types.DefrayRespData{
				Cny: dt,
			},
		}, nil
	case errors.Is(err, defray.MustSelfDefray):
		return &types.DefrayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.MustSelfDefray, "不允许代付"),
		}, nil
	default:
		return &types.DefrayResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.PayFail, errors.WarpQuick(err), "支付失败"),
		}, nil
	}

	jwt.DeleteDefrayToken(req.Token)

	sender.PhoneSendChange(user.Id, "余额（支付）")
	sender.EmailSendChange(user.Id, "余额（支付）")
	sender.MessageSendPay(user.Id, d.Price, d.Subject)
	sender.WxrobotSendPay(user.Id, d.Price, d.Subject)
	sender.FuwuhaoSendDefray(d)
	audit.NewUserAudit(user.Id, "用户支付订单：%s", d.Subject)

	return &types.DefrayResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
