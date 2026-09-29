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

type PhoneTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPhoneTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PhoneTokenLogic {
	return &PhoneTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PhoneTokenLogic) PhoneToken(req *types.CheckPhoneTokenReq) (resp *types.CheckPhoneTokenResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	phone, err := jwt.ParserPhoneToken(req.Token)
	if err != nil {
		return &types.CheckPhoneTokenResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析Token失败"),
			Data: types.CheckPhoneTokenData{
				IsOK: false,
			},
		}, nil
	} else if phone.WebID != web.ID {
		return &types.CheckPhoneTokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Success, "站点不匹配"),
			Data: types.CheckPhoneTokenData{
				IsOK: false,
			},
		}, nil
	}

	if phone.Phone != req.Phone {
		return &types.CheckPhoneTokenResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.CheckPhoneTokenData{
				IsOK: false,
			},
		}, nil
	}

	return &types.CheckPhoneTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CheckPhoneTokenData{
			IsOK: true,
		},
	}, nil
}
