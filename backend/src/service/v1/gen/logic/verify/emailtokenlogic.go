package verify

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type EmailTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailTokenLogic {
	return &EmailTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailTokenLogic) EmailToken(req *types.CheckEmailTokenReq) (resp *types.CheckEmailTokenResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	email, err := jwt.ParserEmailToken(req.Token)
	if err != nil {
		return &types.CheckEmailTokenResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析Token失败"),
			Data: types.CheckEmailTokenData{
				IsOK: false,
			},
		}, nil
	} else if email.WebID != web.ID {
		return &types.CheckEmailTokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Success, "站点不匹配"),
			Data: types.CheckEmailTokenData{
				IsOK: false,
			},
		}, nil
	}

	if email.Email != req.Email {
		return &types.CheckEmailTokenResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.CheckEmailTokenData{
				IsOK: false,
			},
		}, nil
	}

	return &types.CheckEmailTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CheckEmailTokenData{
			IsOK: true,
		},
	}, nil
}
