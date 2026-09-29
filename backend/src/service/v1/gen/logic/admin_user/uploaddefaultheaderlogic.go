package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UploadDefaultHeaderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadDefaultHeaderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadDefaultHeaderLogic {
	return &UploadDefaultHeaderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadDefaultHeaderLogic) UploadDefaultHeader(_ *types.UploadDefaultHeader, r *http.Request) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	headerHeader, ok := r.MultipartForm.File["header"]
	if !ok || len(headerHeader) == 0 {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "缺少header文件"),
		}, nil
	}

	if headerHeader[0].Size >= FileMaxSize {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.FileTooBig, "header文件太大: %d >= %d", headerHeader[0].Size, FileMaxSize),
		}, nil
	}

	headerFile, err := headerHeader[0].Open()
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "打开header文件错误"),
		}, nil
	}
	defer utils.Close(headerFile)

	headerFileByte, err := io.ReadAll(headerFile)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "读取header文件错误"),
		}, nil
	}

	headerFileType := utils.GetMediaType(headerFileByte)
	if !utils.IsImage(headerFileType) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPicType, "header图片类型未知"),
		}, nil
	}

	err = oss.UploadDefaultHeader(headerFileByte, headerFileType)
	if err != nil {
		return nil, respmsg.OSSError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员上传默认头像完成")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
