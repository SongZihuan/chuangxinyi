package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDomainLogic {
	return &GetDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDomainLogic) GetDomain(req *types.GetDoaminReq) (resp *types.GetDomainResp, err error) {
	w := action.GetWebsiteByUID(req.DomainUID)
	if w.Status == db.WebsiteStatusBanned || w.ID == warp.UserCenterWebsite {
		return &types.GetDomainResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DomainNotFound, "外站未找到"),
		}, nil
	}

	return &types.GetDomainResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: w.GetGetDomainDataType(),
	}, nil
}
