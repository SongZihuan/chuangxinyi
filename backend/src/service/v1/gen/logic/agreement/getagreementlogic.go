package agreement

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAgreementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgreementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgreementLogic {
	return &GetAgreementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgreementLogic) GetAgreement(req *types.GetAgreementReq, w http.ResponseWriter, r *http.Request) error {
	agreementModel := db.NewAgreementModel(mysql.MySQLConn)
	agreement, err := agreementModel.FindOneByAid(l.ctx, req.Aid)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("agreement not found")
	} else if err != nil {
		return respmsg.MySQLSystemError.WarpQuick(err)
	}

	_, _ = w.Write([]byte(agreement.Content))
	w.WriteHeader(http.StatusOK)
	return nil
}
