package defray

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/back"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CreateBackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBackLogic {
	return &CreateBackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBackLogic) CreateBack(req *types.CreateBackReq) (resp *types.CreateBackResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	user, err := userModel.FindOneByUidWithoutDelete(l.ctx, req.UserID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.CreateBackResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	backID, err := back.NewBack(l.ctx, req.Get, req.Reason, req.Subject, user, req.CanWithdraw, web.ID)
	if err != nil {
		return &types.CreateBackResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateBackFail, errors.WarpQuick(err), "创建返现订单失败"),
		}, nil
	}

	return &types.CreateBackResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CreateBackData{
			TradeID: backID,
		},
	}, nil
}
