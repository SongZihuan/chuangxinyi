package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/selfpay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type ProcessWithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessWithdrawLogic {
	return &ProcessWithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessWithdrawLogic) ProcessWithdraw(req *types.AdminProcessWithdrawReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	userModel := db.NewUserModel(mysql.MySQLConn)

	withdraw, err := withdrawModel.FindByWithdrawID(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WithdrawNotFound, "提现订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if withdraw.Status != db.WithdrawWait {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WithdrawNotFound, "提现订单已提现"),
		}, nil
	}

	if withdraw.WithdrawWay == alipay.WithdrawAlipay || withdraw.WithdrawWay == wechatpay.WithdrawWechatpay {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayWay, "错误的提现方式"),
		}, nil
	}

	if req.Success {
		err = selfpay.Withdraw(l.ctx, withdraw)
	} else {
		srcUser, err := userModel.FindOneByIDWithoutDelete(l.ctx, withdraw.UserId)
		if errors.Is(err, db.ErrNotFound) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		err = selfpay.WithdrawFail(l.ctx, srcUser, withdraw)
	}
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.PayFail, errors.WarpQuick(err), "处理提现失败"),
		}, nil
	}

	audit.NewAdminAudit(user.Id, "管理员处理提现（%s）", withdraw.WithdrawId)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
