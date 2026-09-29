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

type CheckCompanyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckCompanyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckCompanyLogic {
	return &CheckCompanyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckCompanyLogic) CheckCompany(_ *types.CheckCompany, r *http.Request) (resp *types.SuccessResp, err error) {
	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	err = r.ParseMultipartForm(10 << 20) // 限制最大10MB大小的表单数据
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadForm, errors.WarpQuick(err), "解析表单错误"),
		}, nil
	}

	idcardHeader, ok := r.MultipartForm.File["idcard"]
	if !ok || len(idcardHeader) == 0 {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "没有idcard文件"),
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
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPicType, "未知的idcard图片类型"),
		}, nil
	}

	idcardbackHeader, ok := r.MultipartForm.File["idcardback"]
	if !ok || len(idcardbackHeader) == 0 {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "没有idcardback文件"),
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
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPicType, "未知的idcardback图片类型"),
		}, nil
	}

	licenseHeader, ok := r.MultipartForm.File["license"]
	if !ok || len(licenseHeader) == 0 {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "没有license文件"),
		}, nil
	}

	if licenseHeader[0].Size >= FileMaxSize {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.FileTooBig, "license文件太大: %d >= %d", licenseHeader[0].Size, FileMaxSize),
		}, nil
	}

	licenseFile, err := licenseHeader[0].Open()
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "打开license文件错误"),
		}, nil
	}
	defer utils.Close(licenseFile)

	licenseFileByte, err := io.ReadAll(licenseFile)
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "读取license文件错误"),
		}, nil
	}

	licenseFileType := utils.GetMediaType(licenseFileByte)
	if !utils.IsImage(licenseFileType) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPicType, "未知的license图片类型"),
		}, nil
	}

	id, company, err := func() (id ocr.IDCard, company ocr.Company, resErr errors.WTError) {
		defer utils.Recover(logger.Logger, &resErr, "")

		id, err = ocr.CheckIDCard(idcardFileByte)
		if err != nil {
			return ocr.IDCard{}, ocr.Company{}, errors.WarpQuick(err)
		}

		err = ocr.CheckIDCardBack(idcardbackFileByte)
		if err != nil {
			return ocr.IDCard{}, ocr.Company{}, errors.WarpQuick(err)
		}

		resID, err := yundun.CheckIDCard(id.Name, id.ID)
		if err != nil {
			return ocr.IDCard{}, ocr.Company{}, errors.WarpQuick(err)
		}

		if !resID {
			return ocr.IDCard{}, ocr.Company{}, errors.Errorf("bad idcard")
		}

		company, err = ocr.CheckCompany(licenseFileByte)
		if err != nil {
			return ocr.IDCard{}, ocr.Company{}, errors.WarpQuick(err)
		}

		if company.LegalPerson != id.Name {
			// 法人不对应
			return ocr.IDCard{}, ocr.Company{}, errors.Errorf("bad license")
		}

		resCompany, err := yundun.CheckCompany(company.Name, company.ID, company.LegalPerson)
		if err != nil {
			return ocr.IDCard{}, ocr.Company{}, err
		}

		if !resCompany {
			return ocr.IDCard{}, ocr.Company{}, errors.Errorf("bad license")
		}

		return id, company, nil
	}()
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadLicense, errors.WarpQuick(err), "企业信息鉴别失败"),
		}, nil
	}

	licenseFileName, idcardFileName, idcardBackFileName, err := oss.UploadLicense(licenseFileByte, licenseFileType, idcardFileByte, idcardFileType, idcardbackFileByte, idcardbackFileType, company.Name, company.ID, id.Name, id.ID)
	if err != nil {
		return nil, respmsg.OSSError.WarpQuick(err)
	}

	token, err := jwt.CreateCompanyToken(company.Name, company.ID, id.Name, id.ID, licenseFileName, idcardFileName, idcardBackFileName, web.ID)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:  CompanyToken,
			Token: token,
		},
	}, nil
}
