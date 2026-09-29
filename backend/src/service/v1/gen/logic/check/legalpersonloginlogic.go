package check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
)

type LegalPersonLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLegalPersonLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LegalPersonLoginLogic {
	return &LegalPersonLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LegalPersonLoginLogic) LegalPersonLogin(req *types.LegalPersonLoginReq) (resp *types.SuccessResp, err error) {
	newPhone, err := jwt.ParserPhoneToken(req.PhoneToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if newPhone.WebID != warp.UserCenterWebsite {
		return nil, respmsg.JWTError.New("bad website")
	}

	faceData, err := jwt.ParserFaceToken(req.FaceToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if faceData.WebID != warp.UserCenterWebsite {
		return nil, respmsg.JWTError.New("bad website")
	}

	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	companyModel := db.NewCompanyModel(mysql.MySQLConn)

	user, err := utils2.FindUser(l.ctx, req.ID, false)
	if errors.Is(err, utils2.UserNotFound) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadLoginInfo, "用户未找到", "身份信息不匹配"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if user.FatherId.Valid {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.OnlyRootUser, "仅限非子用户登录", "身份信息不匹配"),
		}, nil
	}

	res1, err := yundun.CheckPhone(faceData.Name, faceData.ID, newPhone.Phone)
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByErrorWithDebug(l.ctx, respmsg.BadLoginInfo, errors.WarpQuick(err), "用户检验失败"),
		}, nil
	}

	if !res1 {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadLoginInfo, "用户检验失败", "身份信息不匹配"),
		}, nil
	}

	_, err = idcardModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadLoginInfo, "用户未进行使用人实名", "身份信息不匹配"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	company, err := companyModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadLoginInfo, "用户未进行企业实名", "身份信息不匹配"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if company.CompanyId != req.CompanyID || company.CompanyName != req.CompanyName {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadLoginInfo, "企业实名信息不匹配", "身份信息不匹配"),
		}, nil
	}

	if company.LegalPersonName != req.LegalPersonName || company.LegalPersonIdCard != req.LegalPersonID {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadLoginInfo, "企业实名信息不匹配", "身份信息不匹配"),
		}, nil
	}

	if company.LegalPersonName != faceData.Name || company.LegalPersonIdCard != faceData.ID {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadLoginInfo, "企业实名信息不匹配", "身份信息不匹配"),
		}, nil
	}

	var token string
	var subType string
	token, err = jwt.CreateUserToken(l.ctx, user.Uid, user.Signin, user.TokenExpiration, jwt2.UserHighAuthorityRootToken, "", 0)
	subType = jwt2.UserHighAuthorityRootTokenString
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	sender.MessageSendLoginCenter(user.Id, l.ctx)
	sender.WxrobotSendLoginCenter(user.Id, l.ctx)
	sender.FuwuhaoSendLoginCenter(user.Id)
	audit.NewUserAudit(user.Id, "用户通过法人身份登录成功, 人脸识别代号为：%s", faceData.CheckID)

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:    UserToken,
			Token:   token,
			SubType: subType,
		},
	}, nil
}
