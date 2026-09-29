package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
)

type UploadCompanyInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCompanyInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCompanyInfoLogic {
	return &UploadCompanyInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCompanyInfoLogic) UploadCompanyInfo(req *types.UploadCompanyInfoReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	companyModel := db.NewCompanyModel(mysql.MySQLConn)
	company, err := companyModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotVerify, "需要先企业实名"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	companyData, err := jwt.ParserCompanyToken(req.CompanyToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if companyData.WebID != warp.UserCenterWebsite {
		return nil, respmsg.JWTError.New("bad website")
	}

	if company.LegalPersonName != companyData.LegalPersonName || company.LegalPersonIdCard != companyData.LegalPersonID {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NotCorrespond, "企业信息不匹配"),
		}, nil
	}

	if company.CompanyName != companyData.Name || company.CompanyId != companyData.ID {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NotCorrespond, "企业信息不匹配"),
		}, nil
	}

	company.LicenseKey = sql.NullString{
		Valid:  true,
		String: companyData.License,
	}
	company.IdcardKey = sql.NullString{
		Valid:  true,
		String: companyData.IDCard,
	}
	company.IdcardBackKey = sql.NullString{
		Valid:  true,
		String: companyData.IDCardBack,
	}

	err = companyModel.UpdateCh(l.ctx, company)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户企业实名信息（原件）上传成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
