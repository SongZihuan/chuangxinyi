package admin_agreement

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

type CretaeAgreementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAgreementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CretaeAgreementLogic {
	return &CretaeAgreementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CretaeAgreementLogic) CreateAgreement(req *types.CreateAgreementReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	agreementModel := db.NewAgreementModel(mysql.MySQLConn)
	_, err = agreementModel.FindOneByAid(l.ctx, req.Aid)
	if !errors.Is(err, db.ErrNotFound) {
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.AgreementExists, "协议不存在"),
		}, nil
	}

	agreement := &db.Agreement{
		Aid:     req.Aid,
		Content: req.Content,
	}

	_, err = agreementModel.Insert(l.ctx, agreement)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员添加协议（%s）成功", req.Aid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
