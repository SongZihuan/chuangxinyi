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
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UploadUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadUserInfoLogic {
	return &UploadUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadUserInfoLogic) UploadUserInfo(req *types.UploadUserInfoReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	idcard, err := idcardModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotVerify, "用户需要先使用人实名"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	idcardData, err := jwt.ParserIDCardToken(req.IDCardToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if idcardData.WebID != warp.UserCenterWebsite {
		return nil, respmsg.JWTError.New("bad website")
	}

	if idcard.UserName != idcardData.Name || idcard.UserIdCard != idcardData.ID {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NotCorrespond, "使用人信息不匹配"),
		}, nil
	}

	idcard.IdcardKey = sql.NullString{
		Valid:  true,
		String: idcardData.IDCard,
	}
	idcard.IdcardBackKey = sql.NullString{
		Valid:  true,
		String: idcardData.IDCardBack,
	}

	err = idcardModel.UpdateCh(l.ctx, idcard)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户使用者实名信息（原件）上传成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
