package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UploadCompanyFaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCompanyFaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCompanyFaceLogic {
	return &UploadCompanyFaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCompanyFaceLogic) UploadCompanyFace(req *types.UploadCompanyFaceReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	companyModel := db.NewCompanyModel(mysql.MySQLConn)
	company, err := companyModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotVerify, "用户需要先企业实名"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	faceData, err := jwt.ParserFaceToken(req.FaceToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if faceData.WebID != warp.UserCenterWebsite {
		return nil, respmsg.JWTError.New("bad website")
	}

	if company.LegalPersonName != faceData.Name || company.LegalPersonIdCard != faceData.ID {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NotCorrespond, "信息不匹配"),
		}, nil
	}

	company.FaceCheckId = sql.NullString{
		Valid:  true,
		String: faceData.CheckID,
	}

	err = companyModel.UpdateCh(l.ctx, company)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户企业实名信息（人脸识别）上传成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
