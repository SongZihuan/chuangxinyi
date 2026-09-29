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

type FaceTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFaceTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FaceTokenLogic {
	return &FaceTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FaceTokenLogic) FaceToken(req *types.CheckFaceTokenReq) (resp *types.CheckFaceTokenResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	faceData, err := jwt.ParserFaceToken(req.Token)
	if err != nil {
		return &types.CheckFaceTokenResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析Token失败"),
			Data: types.CheckFaceTokenData{
				IsOK: false,
			},
		}, nil
	} else if faceData.WebID != web.ID {
		return &types.CheckFaceTokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Success, "站点不匹配"),
			Data: types.CheckFaceTokenData{
				IsOK: false,
			},
		}, nil
	}

	if faceData.Name != req.Name || faceData.ID != req.ID {
		return &types.CheckFaceTokenResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.CheckFaceTokenData{
				IsOK: false,
			},
		}, nil
	}

	return &types.CheckFaceTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CheckFaceTokenData{
			IsOK: true,
		},
	}, nil
}
