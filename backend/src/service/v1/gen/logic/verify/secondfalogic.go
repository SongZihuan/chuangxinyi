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

type SecondFALogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSecondFALogic(ctx context.Context, svcCtx *svc.ServiceContext) *SecondFALogic {
	return &SecondFALogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SecondFALogic) SecondFA(req *types.SecondFACheckReq) (resp *types.CheckSecondFATokenResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	uid, err := jwt.ParserCheck2FAToken(req.Token)
	if err != nil {
		return &types.CheckSecondFATokenResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析Token失败"),
			Data: types.CheckSecondFATokenData{
				IsOK: false,
			},
		}, nil
	} else if uid.WebID != web.ID {
		return &types.CheckSecondFATokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Success, "站点不匹配"),
			Data: types.CheckSecondFATokenData{
				IsOK: false,
			},
		}, nil
	}

	if uid.UserID != req.UserID {
		return &types.CheckSecondFATokenResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.CheckSecondFATokenData{
				IsOK: false,
			},
		}, nil
	}

	return &types.CheckSecondFATokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CheckSecondFATokenData{
			IsOK: true,
		},
	}, nil
}
