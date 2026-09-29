package check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/ocr"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CheckUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserLogic {
	return &CheckUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckUserLogic) CheckUser(_ *types.CheckUser, r *http.Request) (resp *types.SuccessResp, err error) {
	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	err = r.ParseMultipartForm(10 << 20) // 限制最大10MB大小的表单数据
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadForm, errors.WarpQuick(err), "读取表单失败"),
		}, nil
	}

	idcardHeader, ok := r.MultipartForm.File["idcard"]
	if !ok || len(idcardHeader) == 0 {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "缺少idcard文件"),
		}, nil
	}

	if idcardHeader[0].Size >= FileMaxSize {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.FileTooBig, "idcard文件太大: %d >= %d", idcardHeader[0].Size, FileMaxSize),
		}, nil
	}

	idcardFile, err := idcardHeader[0].Open()
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "打开idcard文件错误"),
		}, nil
	}
	defer utils.Close(idcardFile)

	idcardFileByte, err := io.ReadAll(idcardFile)
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "读取idcard文件错误"),
		}, nil
	}

	idcardFileType := utils.GetMediaType(idcardFileByte)
	if !utils.IsImage(idcardFileType) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPicType, "idcard图片类型未知"),
		}, nil
	}

	idcardbackHeader, ok := r.MultipartForm.File["idcardback"]
	if !ok || len(idcardbackHeader) == 0 {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "缺少idcardback文件"),
		}, nil
	}

	if idcardbackHeader[0].Size >= FileMaxSize {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.FileTooBig, "idcardback文件太大: %d >= %d", idcardbackHeader[0].Size, FileMaxSize),
		}, nil
	}

	idcardbackFile, err := idcardbackHeader[0].Open()
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "打开idcardback文件错误"),
		}, nil
	}
	defer utils.Close(idcardbackFile)

	idcardbackFileByte, err := io.ReadAll(idcardbackFile)
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "读取idcardback文件错误"),
		}, nil
	}

	idcardbackFileType := utils.GetMediaType(idcardbackFileByte)
	if !utils.IsImage(idcardbackFileType) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPicType, "idcardback图片类型未知"),
		}, nil
	}

	id, err := func() (id ocr.IDCard, resErr errors.WTError) {
		defer utils.Recover(logger.Logger, &resErr, "")

		id, err = ocr.CheckIDCard(idcardFileByte)
		if err != nil {
			return ocr.IDCard{}, errors.WarpQuick(err)
		}

		err = ocr.CheckIDCardBack(idcardbackFileByte)
		if err != nil {
			return ocr.IDCard{}, errors.WarpQuick(err)
		}

		res1, err := yundun.CheckIDCard(id.Name, id.ID)
		if err != nil {
			return ocr.IDCard{}, errors.WarpQuick(err)
		}

		if !res1 {
			// 身份证识别错误
			return ocr.IDCard{}, errors.Errorf("bad idcard")
		}

		return id, nil
	}()
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadIDCard, errors.WarpQuick(err), "身份证识别错误"),
		}, nil
	}

	idcardFileName, idcardbackFileName, err := oss.UploadIDCard(idcardFileByte, idcardFileType, idcardbackFileByte, idcardbackFileType, id.Name, id.ID)
	if err != nil {
		return nil, respmsg.OSSError.WarpQuick(err)
	}

	token, err := jwt.CreateIDCardToken(id.Name, id.ID, idcardFileName, idcardbackFileName, web.ID)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:  IDCardToken,
			Token: token,
		},
	}, nil
}
