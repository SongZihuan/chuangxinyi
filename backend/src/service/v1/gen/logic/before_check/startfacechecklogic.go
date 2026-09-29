package before_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type StartFaceCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartFaceCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartFaceCheckLogic {
	return &StartFaceCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartFaceCheckLogic) StartFaceCheck(req *types.StartFackCheckReq) (resp *types.StartFackCheckResp, err error) {
	id, _, err := alipay.NewFaceCheck(l.ctx, req.Name, req.ID)
	if err != nil {
		return &types.StartFackCheckResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFaceIDCardOrName, errors.WarpQuick(err), "启动人脸身份认证失败"),
		}, nil
	}

	return &types.StartFackCheckResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.StartFackCheckData{
			CertifyID: id,
		},
	}, nil
}
