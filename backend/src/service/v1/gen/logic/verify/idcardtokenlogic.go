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

type IDCardTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIDCardTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IDCardTokenLogic {
	return &IDCardTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IDCardTokenLogic) IDCardToken(req *types.CheckIDCardTokenReq) (resp *types.CheckIDCardTokenResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	idcardData, err := jwt.ParserIDCardToken(req.Token)
	if err != nil {
		return &types.CheckIDCardTokenResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析Token失败"),
			Data: types.CheckIDCardTokenData{
				IsOK: false,
			},
		}, nil
	} else if idcardData.WebID != web.ID {
		return &types.CheckIDCardTokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Success, "站点不匹配"),
			Data: types.CheckIDCardTokenData{
				IsOK: false,
			},
		}, nil
	}

	if idcardData.Name != req.Name || idcardData.ID != req.ID {
		return &types.CheckIDCardTokenResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.CheckIDCardTokenData{
				IsOK: false,
			},
		}, nil
	}

	return &types.CheckIDCardTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CheckIDCardTokenData{
			IsOK: true,
		},
	}, nil
}
