package admin_user

import (
	"context"
	"encoding/base64"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.AdminGetUserReq) (resp *types.AdminGetUserResp, err error) {
	user, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetUserResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	userNameModel := db.NewUsernameModel(mysql.MySQLConn)
	nickNameModel := db.NewNicknameModel(mysql.MySQLConn)
	headerModel := db.NewHeaderModel(mysql.MySQLConn)
	emailModel := db.NewEmailModel(mysql.MySQLConn)
	passwordModel := db.NewPasswordModel(mysql.MySQLConn)
	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	WxRobotModel := db.NewWxrobotModel(mysql.MySQLConn)
	AddressModel := db.NewAddressModel(mysql.MySQLConn)
	ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

	r := action.GetRole(user.RoleId, user.IsAdmin)

	phone, err := phoneModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		phone = &db.Phone{}
		if db.IsKeepInfoStatus(user.Status) {
			logger.Logger.Error("user not phone: %d, %s", user.Id, user.Uid)
		}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	email, err := emailModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		email = &db.Email{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	username, err := userNameModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		username = &db.Username{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userNameString, err := base64.StdEncoding.DecodeString(username.Username.String)
	if err != nil {
		userNameString = []byte("")
	}

	nickname, err := nickNameModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		nickname = &db.Nickname{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	header, err := headerModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		header = &db.Header{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	password, err := passwordModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		password = &db.Password{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	secondfa, err := secondfaModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		secondfa = &db.Secondfa{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	wechat, err := wechatModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		wechat = &db.Wechat{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	wxrobot, err := WxRobotModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		wxrobot = &db.Wxrobot{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	address, err := AddressModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		address = &db.Address{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	ctrl, err := ctrlModel.FindByUserID(context.Background(), user.Id)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.AdminGetUserResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetUserData{
			ID:              user.Id,
			UID:             user.Uid,
			Status:          db.UserStatusMap[user.Status],
			Signin:          user.Signin,
			Father:          user.FatherId.Int64,
			Invite:          user.InviteId.Int64,
			TokenExpiration: user.TokenExpiration,
			RoleID:          r.ID,
			RoleName:        r.Name,
			RoleSign:        r.Sign,
			IsAdmin:         user.IsAdmin,
			CreateAt:        user.CreateAt.Unix(),

			Phone: phone.Phone,
			Email: email.Email.String,

			Nickname: nickname.Nickname.String,
			Header:   header.Header.String,

			WxOpenID:      wechat.OpenId.String,
			WxUnionID:     wechat.UnionId.String,
			FuwuhaoOpenID: wechat.Fuwuhao.String,
			WxNickName:    wechat.Nickname.String,
			WxHeader:      wechat.Headimgurl.String,

			WxWebHook: wxrobot.Webhook.String,

			HasPassword: password.PasswordHash.Valid,
			Has2FA:      secondfa.Secret.Valid,

			UserName: string(userNameString),

			AddressName:     address.Name.String,
			AddressPhone:    address.Phone.String,
			AddressEmail:    address.Email.String,
			AddressProvince: address.Province.String,
			AddressCity:     address.City.String,
			AddressDistrict: address.District.String,
			AddressAddress:  address.Address.String,
			AddressArea:     address.GetAreaList(),

			AllowPhone:    ctrl.AllowPhone,
			AllowPassword: ctrl.AllowPassword,
			AllowEmail:    ctrl.AllowEmail,
			AllowWeChat:   ctrl.AllowWechat,
			AllowSecondFA: ctrl.Allow2Fa,
		},
	}, nil
}
