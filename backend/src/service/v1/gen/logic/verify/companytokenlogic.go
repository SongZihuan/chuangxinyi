package verify

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CompanyTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompanyTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompanyTokenLogic {
	return &CompanyTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompanyTokenLogic) CompanyToken(req *types.CheckCompanyTokenReq) (resp *types.CheckCompanyTokenResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	companyData, err := jwt.ParserCompanyToken(req.Token)
	if err != nil {
		return &types.CheckCompanyTokenResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析Token失败"),
			Data: types.CheckCompanyTokenData{
				IsOK: false,
			},
		}, nil
	} else if companyData.WebID != web.ID {
		return &types.CheckCompanyTokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Success, "站点不匹配"),
			Data: types.CheckCompanyTokenData{
				IsOK: false,
			},
		}, nil
	}

	if companyData.LegalPersonName != req.Name || companyData.LegalPersonID != req.ID || companyData.Name != req.CompanyName || companyData.ID != req.CompanyID {
		return &types.CheckCompanyTokenResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.CheckCompanyTokenData{
				IsOK: false,
			},
		}, nil
	}

	return &types.CheckCompanyTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CheckCompanyTokenData{
			IsOK: true,
		},
	}, nil
}
