package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UploadCompanyInfoByJsonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCompanyInfoByJsonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCompanyInfoByJsonLogic {
	return &UploadCompanyInfoByJsonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCompanyInfoByJsonLogic) UploadCompanyInfoByJson(req *types.UploadCompanyInfoByJson) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !utils.IsValidIDCard(req.LegalPersonIDCard) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadCompanyInfo, "错误的法人身份证号码"),
		}, nil
	}

	if !utils.IsValidChineseName(req.LegalPersonName) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadCompanyInfo, "错误的法人姓名"),
		}, nil
	}

	if !utils.IsValidCreditCode(req.CompanyID) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadCompanyInfo, "错误的企业统一社会信用代码"),
		}, nil
	}

	if !utils.IsValidChineseCompanyName(req.CompanyName) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadCompanyInfo, "错误的企业名称"),
		}, nil
	}

	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	_, err = idcardModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WithoutVerify, "需要使用人实名先"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	err = func() error {
		defer utils.Recover(logger.Logger, nil, "")

		resID, err := yundun.CheckIDCard(req.LegalPersonName, req.LegalPersonIDCard)
		if err != nil {
			return err
		}

		if !resID {
			// 身份证识别错误
			return errors.Errorf("bad id info")
		}

		resCompany, err := yundun.CheckCompany(req.CompanyName, req.CompanyID, req.LegalPersonName)
		if err != nil {
			return err
		}

		if !resCompany {
			// 企业营业执照
			return errors.Errorf("bad company info")
		}

		return nil
	}()
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadLicense, errors.WarpQuick(err), "企业信息检验失败"),
		}, nil
	}

	companyModel := db.NewCompanyModel(mysql.MySQLConn)
	_, err = companyModel.InsertWithDelete(context.Background(), &db.Company{
		UserId:            user.Id,
		LegalPersonName:   req.LegalPersonName,
		LegalPersonIdCard: req.LegalPersonIDCard,
		CompanyName:       req.CompanyName,
		CompanyId:         req.CompanyID,
	})
	if err != nil {
		return
	}

	audit.NewUserAudit(user.Id, "用户更新企业实名成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
