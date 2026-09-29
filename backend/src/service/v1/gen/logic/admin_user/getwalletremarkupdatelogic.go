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

type GetWalletRemarkUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWalletRemarkUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletRemarkUpdateLogic {
	return &GetWalletRemarkUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWalletRemarkUpdateLogic) GetWalletRemarkUpdate(req *types.AdminGetWalletRemarkUpdate) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	recordModel := db.NewWalletRecordModel(mysql.MySQLConn)
	record, err := recordModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WalletRecordNotFound, "钱包未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	record.Remark = req.Remark

	err = recordModel.Update(l.ctx, record)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新用户钱包记录（%d）备注成功", record.Id)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
