package admin_agreement

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAgreementListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgreementListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgreementListLogic {
	return &GetAgreementListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgreementListLogic) GetAgreementList(req *types.PageReq) (resp *types.GetAgreementListResp, err error) {
	agreementModel := db.NewAgreementModel(mysql.MySQLConn)
	agreementList, err := agreementModel.GetList(l.ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := agreementModel.GetCount(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]string, 0, len(agreementList))
	for _, a := range agreementList {
		respList = append(respList, a.Aid)
	}

	return &types.GetAgreementListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetAgreementListData{
			Count:     count,
			Agreement: respList,
		},
	}, nil
}
