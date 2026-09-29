package check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CheckFaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckFaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckFaceLogic {
	return &CheckFaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckFaceLogic) CheckFace(req *types.CheckFaceReq) (resp *types.SuccessResp, err error) {
	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	name, id, status, checkID, err := alipay.QueryFaceCheck(l.ctx, req.CertifyID)
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.FaceCheckFail, errors.WarpQuick(err), "人脸验证查询失败"),
		}, nil
	} else if status != db.FaceCheckWait {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.FaceCheckWait, "人脸验证等待中"),
		}, nil
	} else if status != db.FaceCheckFail {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.FaceCheckFail, "人脸验证失败"),
		}, nil
	}

	token, err := jwt.CreateFaceToken(name, id, req.CertifyID, checkID, web.ID)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:  FaceToken,
			Token: token,
		},
	}, nil
}
