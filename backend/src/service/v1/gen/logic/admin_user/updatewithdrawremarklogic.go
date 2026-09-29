package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateWithdrawRemarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWithdrawRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWithdrawRemarkLogic {
	return &UpdateWithdrawRemarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWithdrawRemarkLogic) UpdateWithdrawRemark(req *types.UpdateWithdrawRemarkReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	withdraw, err := withdrawModel.FindByWithdrawID(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WithdrawNotFound, "提现订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	withdraw.Remark = req.Remark

	err = withdrawModel.Update(l.ctx, withdraw)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新提现（%s）备注成功", withdraw.WithdrawAt)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
