package ui

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetFooterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFooterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFooterLogic {
	return &GetFooterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFooterLogic) GetFooter() (resp *types.FooterResp, err error) {
	footerModel := db.NewFooterModel(mysql.MySQLConn)
	footer, err := footerModel.FindTheNew(l.ctx)
	if errors.Is(err, db.ErrNotFound) {
		return &types.FooterResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.FooterData{
				ICP1:      "",
				ICP2:      "",
				Gongan:    "",
				Copyright: "",
			},
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.FooterResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.FooterData{
			ICP1:      footer.Icp1,
			ICP2:      footer.Icp2,
			Gongan:    footer.Gongan,
			Copyright: footer.Copyright,
		},
	}, nil
}
