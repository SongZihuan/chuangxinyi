package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateTitleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTitleLogic {
	return &UpdateTitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTitleLogic) UpdateTitle(req *types.UserUpdateTitleReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	titleModel := db.NewTitleModel(mysql.MySQLConn)
	_, err = titleModel.InsertWithDelete(context.Background(), &db.Title{
		UserId: user.Id,
		Name: sql.NullString{
			Valid:  len(req.Name) != 0,
			String: req.Name,
		},
		TaxId: sql.NullString{
			Valid:  len(req.TaxID) != 0,
			String: req.TaxID,
		},
		BankId: sql.NullString{
			Valid:  len(req.BankID) != 0,
			String: req.BankID,
		},
		Bank: sql.NullString{
			Valid:  len(req.Bank) != 0,
			String: req.Bank,
		},
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户更新发票抬头成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
