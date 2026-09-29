package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateHeaderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateHeaderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHeaderLogic {
	return &UpdateHeaderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateHeaderLogic) UpdateHeader(req *types.UpdateHeader, r *http.Request) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.IsDelete {
		headerModel := db.NewHeaderModel(mysql.MySQLConn)
		_, err = headerModel.InsertWithDelete(context.Background(), &db.Header{
			UserId: user.Id,
			Header: sql.NullString{
				Valid: false,
			},
		})

		return &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		}, nil
	}

	headerHeader, ok := r.MultipartForm.File["header"]
	if !ok || len(headerHeader) == 0 {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "缺少header文件"),
		}, nil
	}

	if headerHeader[0].Size >= FileMaxSize {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.FileTooBig, "头像文件太大: %d >= %d", headerHeader[0].Size, FileMaxSize),
		}, nil
	}

	headerFile, err := headerHeader[0].Open()
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "打开头像文件失败"),
		}, nil
	}
	defer utils.Close(headerFile)

	headerFileByte, err := io.ReadAll(headerFile)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "读取头像文件失败"),
		}, nil
	}

	headerFileType := utils.GetMediaType(headerFileByte)
	if !utils.IsImage(headerFileType) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPicType, "未知的头像图片类型"),
		}, nil
	}

	res, err := yundun.CheckHeaderPic(headerFileByte, headerFileType)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadHeader, errors.WarpQuick(err), "头像检测失败"),
		}, nil
	}

	if !res {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadHeader, "头像违规"),
		}, nil
	}

	headerModel := db.NewHeaderModel(mysql.MySQLConn)

	err = oss.UploadHeader(headerFileByte, user.Uid, headerFileType)
	if err != nil {
		return nil, respmsg.OSSError.WarpQuick(err)
	}

	_, err = headerModel.InsertWithDelete(context.Background(), &db.Header{
		UserId: user.Id,
		Header: sql.NullString{
			Valid:  true,
			String: user.Uid,
		},
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.MessageSend(user.Id, "头像上传", "头像已经上传成功啦！")
	audit.NewUserAudit(user.Id, "用户更新头像成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
