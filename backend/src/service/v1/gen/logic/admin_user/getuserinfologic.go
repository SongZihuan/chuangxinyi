package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.AdminGetUserReq) (resp *types.AdminGetUserInfoResp, err error) {
	user, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetUserInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	companyModel := db.NewCompanyModel(mysql.MySQLConn)

	idcard, err := idcardModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		idcard = &db.Idcard{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	company, err := companyModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		company = &db.Company{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.AdminGetUserInfoResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetUserInfoData{
			UserName:                 idcard.UserName,
			UserIDCard:               idcard.UserIdCard,
			UserIDCardUrl:            GetIdentityPic(idcard.IdcardKey.String),
			UserIDCardBackUrl:        GetIdentityPic(idcard.IdcardBackKey.String),
			CompanyName:              company.CompanyName,
			CompanyID:                company.CompanyId,
			LicenseUrl:               GetIdentityPic(company.LicenseKey.String),
			LegalPersonName:          company.LegalPersonName,
			LegalPersonIDCard:        company.LegalPersonIdCard,
			LegalPersonIDCardUrl:     GetIdentityPic(company.IdcardKey.String),
			LegalPersonIDCardBackUrl: GetIdentityPic(company.IdcardBackKey.String),
		},
	}, nil
}

func GetIdentityPic(key string) string {
	if len(key) == 0 {
		return ""
	}

	url, err := oss.GetIdentity(key)
	if err != nil {
		logger.Logger.Error("oss error: %s", err.Error())
	}

	return url
}
