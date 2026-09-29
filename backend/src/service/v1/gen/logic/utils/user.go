package utils

import (
	"context"
	"database/sql"
	"encoding/base64"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	errors "github.com/wuntsong-org/wterrors"
)

var UserNotFound = errors.Errorf("UserNotFound")

func FindUser(ctx context.Context, userID string, findBanned bool) (*db.User, errors.WTError) {
	userModel := db.NewUserModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	emailModel := db.NewEmailModel(mysql.MySQLConn)
	usernameModel := db.NewUsernameModel(mysql.MySQLConn)

	user, err := userModel.FindOneByUidWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		user = nil
	} else if err != nil {
		return nil, errors.WarpQuick(err)
	}

	id := int64(0)
	if user == nil {
		phone, err := phoneModel.FindByPhone(ctx, userID)
		if errors.Is(err, db.ErrNotFound) {
		} else if err != nil {
			return nil, errors.WarpQuick(err)
		} else {
			id = phone.UserId
		}
	}

	if user == nil && id == 0 {
		email, err := emailModel.FindByEmail(ctx, userID)
		if errors.Is(err, db.ErrNotFound) {
		} else if err != nil {
			return nil, errors.WarpQuick(err)
		} else {
			id = email.UserId
		}
	}

	if user == nil && id == 0 {
		email, err := usernameModel.FindByUsernameWithBase64(ctx, userID)
		if errors.Is(err, db.ErrNotFound) {
		} else if err != nil {
			return nil, errors.WarpQuick(err)
		} else {
			id = email.UserId
		}
	}

	if user == nil && id == 0 {
		return nil, UserNotFound
	} else if user == nil && id != 0 {
		user, err = userModel.FindOneByIDWithoutDelete(ctx, id)
		if errors.Is(err, db.ErrNotFound) {
			return nil, UserNotFound
		} else if err != nil {
			return nil, errors.WarpQuick(err)
		}
	}

	if !findBanned && db.IsBanned(user) {
		return nil, UserNotFound
	}

	return user, nil
}

type InfoData struct {
	User     types.UserEasy     `json:"user"`
	Info     types.UserInfo     `json:"info"`
	Data     types.UserData     `json:"data"`
	Balance  types.UserBalance  `json:"balance"`
	Title    types.Title        `json:"title"`
	Address  types.Address      `json:"address"`
	Role     types.Role         `json:"role"`
	InfoEasy types.UserInfoEsay `json:"infoEasy"`
}

type InfoDataEasy struct {
	User     types.UserEasy     `json:"user"`
	Data     types.UserData     `json:"data"`
	InfoEasy types.UserInfoEsay `json:"infoEasy"`
}

type InfoDataWebsite struct {
	User     types.WebsiteUserEasy `json:"user"`
	Data     types.UserData        `json:"data"`
	InfoEasy types.UserInfoEsay    `json:"infoEasy"`
}

func GetUserInfo(ctx context.Context, user *db.User, subType int) (InfoData, errors.WTError) {
	userModel := db.NewUserModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	userNameModel := db.NewUsernameModel(mysql.MySQLConn)
	nickNameModel := db.NewNicknameModel(mysql.MySQLConn)
	headerModel := db.NewHeaderModel(mysql.MySQLConn)
	emailModel := db.NewEmailModel(mysql.MySQLConn)
	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	companyModel := db.NewCompanyModel(mysql.MySQLConn)
	passwordModel := db.NewPasswordModel(mysql.MySQLConn)
	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	titleModel := db.NewTitleModel(mysql.MySQLConn)
	addressModel := db.NewAddressModel(mysql.MySQLConn)
	wxrobotModel := db.NewWxrobotModel(mysql.MySQLConn)
	ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

	role := action.GetRole(user.RoleId, user.IsAdmin)
	userRole := role.GetRole()
	userRole.Menus = action.ClearRoleMenu(userRole.Menus, role, subType)
	userRole.UrlPaths = action.ClearRoleUrlPath(userRole.UrlPaths, role, subType)

	phone, err := phoneModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		phone = &db.Phone{
			Phone: "",
		}
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	username, err := userNameModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		username = &db.Username{
			Username: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	userNameString, err := base64.StdEncoding.DecodeString(username.Username.String)
	if err != nil {
		userNameString = []byte("")
	}

	nickname, err := nickNameModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		nickname = &db.Nickname{
			Nickname: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	header, err := headerModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		header = &db.Header{
			Header: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	email, err := emailModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		email = &db.Email{
			Email: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	idcard, err := idcardModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		idcard = nil
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	company, err := companyModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		company = nil
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	password, err := passwordModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		password = nil
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	secondfa, err := secondfaModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		secondfa = nil
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	wechat, err := wechatModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		wechat = &db.Wechat{
			UnionId: sql.NullString{
				String: "",
			},
			Nickname: sql.NullString{
				String: "",
			},
			Headimgurl: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	wxrobot, err := wxrobotModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		wxrobot = nil
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	title, err := titleModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		title = nil
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	address, err := addressModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		address = nil
	} else if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	inviteCount, err := userModel.GetUserInviteCount(ctx, user.Id)
	if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	ctrl, err := ctrlModel.FindByUserID(ctx, user.Id)
	if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	userEasy := types.UserEasy{
		UID:            user.Uid,
		InviteCount:    inviteCount,
		RoleID:         user.RoleId,
		RoleName:       userRole.Name,
		RoleSign:       userRole.Sign,
		Signin:         user.Signin,
		TokenExpire:    user.TokenExpiration,
		Phone:          phone.Phone,
		UserName:       string(userNameString),
		NickName:       nickname.Nickname.String,
		Header:         header.Header.String,
		Email:          email.Email.String,
		WeChatNickName: wechat.Nickname.String,
		WeChatHeader:   wechat.Headimgurl.String,
		Status:         db.UserStatusMap[user.Status],
		CreateAt:       user.CreateAt.Unix(),
	}

	userInfo := types.UserInfo{}
	if idcard == nil {
		userEasy.UserRealName = ""

		userInfo.HasVerified = false
		userInfo.IsCompany = false
	} else {
		userEasy.UserRealName = idcard.UserName

		userInfo.HasVerified = true
		userInfo.UserName = idcard.UserName
		userInfo.UserIDCard = idcard.UserIdCard

		if idcard.Phone.Valid {
			userInfo.VerifiedPhone = idcard.Phone.String
		} else {
			userInfo.VerifiedPhone = ""
		}

		if company != nil {
			userEasy.CompanyName = company.CompanyName

			userInfo.IsCompany = true
			userInfo.CompanyName = company.CompanyName
			userInfo.LegalPersonName = company.LegalPersonName

			if subType == jwt.UserRootToken || subType == jwt.UserHighAuthorityRootToken || subType == jwt.UserRootFatherToken {
				userInfo.CompanyID = company.CompanyId
				userInfo.LegalPersonIDCard = company.LegalPersonIdCard
			}
		} else {
			userEasy.CompanyName = ""

			userInfo.IsCompany = false
		}
	}

	userData := types.UserData{
		Has2FA:                  secondfa != nil && secondfa.Secret.Valid,
		HasEmail:                email.Email.Valid,
		HasPassword:             password != nil && password.PasswordHash.Valid,
		HasWeChat:               wechat.OpenId.Valid,
		HasFuwuhao:              wechat.Fuwuhao.Valid,
		HasUnionID:              wechat.UnionId.Valid,
		HasWxrobot:              wxrobot != nil && wxrobot.Webhook.Valid,
		HasVerified:             idcard != nil,
		IsCompany:               idcard != nil && company != nil,
		VerifiedPhone:           "", // 下面设置
		HasUserOriginal:         idcard != nil && idcard.IdcardKey.Valid,
		HasUserFaceCheck:        idcard != nil && idcard.FaceCheckId.Valid,
		HasCompanyOriginal:      idcard != nil && company != nil && company.IdcardKey.Valid,
		HasLegalPersonFaceCheck: idcard != nil && company != nil && company.FaceCheckId.Valid,
		AllowPhone:              ctrl.AllowPhone,
		AllowPassword:           ctrl.AllowPassword,
		AllowEmail:              ctrl.AllowEmail,
		AllowWeChat:             ctrl.AllowWechat,
		AllowSecondFA:           ctrl.Allow2Fa,
	}

	if idcard != nil && idcard.Phone.Valid {
		userData.VerifiedPhone = idcard.Phone.String
	}

	b, err := balance.QueryBalance(ctx, user.Id)
	if err != nil {
		return InfoData{}, errors.WarpQuick(err)
	}

	userBalance := types.UserBalance{
		Balance:      b.Balance,
		WaitBalance:  b.WaitBalance,
		Cny:          b.Cny,
		NotBilled:    b.NotBilled,
		Billed:       b.Billed,
		HasBilled:    b.HasBilled,
		HasWithdraw:  b.HasWithdraw,
		NotWithdraw:  b.NotWithdraw,
		Withdraw:     b.Withdraw,
		WaitWithdraw: b.WaitWithdraw,
		WalletID:     user.WalletId,
	}

	userTitle := types.Title{}
	if title != nil {
		userTitle.Name = title.Name.String
		userTitle.TaxID = title.TaxId.String
		userTitle.BankID = title.BankId.String
		userTitle.Bank = title.Bank.String
	}

	userAddress := types.Address{}
	if address != nil {
		userAddress.Name = address.Name.String
		userAddress.Phone = address.Phone.String
		userAddress.Email = address.Email.String
		userAddress.Province = address.Province.String
		userAddress.City = address.City.String
		userAddress.District = address.District.String
		userAddress.Address = address.Address.String
		userAddress.Area = address.GetAreaList()
	}

	userInfoEasy := types.UserInfoEsay{}
	if idcard != nil {
		userInfoEasy.HasVerified = true
		userInfoEasy.UserName = idcard.UserName

		if company != nil {
			userInfoEasy.IsCompany = true
			userInfoEasy.LegalPersonName = company.LegalPersonName
			userInfoEasy.CompanyName = company.CompanyName
		} else {
			userInfoEasy.IsCompany = false
		}
	} else {
		userInfoEasy.HasVerified = false
	}

	return InfoData{
		User:     userEasy,
		Info:     userInfo,
		Data:     userData,
		Balance:  userBalance,
		Address:  userAddress,
		Title:    userTitle,
		Role:     userRole,
		InfoEasy: userInfoEasy,
	}, nil
}

func GetUserInfoEasy(ctx context.Context, user *db.User) (InfoDataEasy, errors.WTError) {
	userModel := db.NewUserModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	userNameModel := db.NewUsernameModel(mysql.MySQLConn)
	nickNameModel := db.NewNicknameModel(mysql.MySQLConn)
	headerModel := db.NewHeaderModel(mysql.MySQLConn)
	emailModel := db.NewEmailModel(mysql.MySQLConn)
	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	companyModel := db.NewCompanyModel(mysql.MySQLConn)
	passwordModel := db.NewPasswordModel(mysql.MySQLConn)
	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	wxrobotModel := db.NewWxrobotModel(mysql.MySQLConn)
	ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

	userRole := action.GetRole(user.RoleId, user.IsAdmin)

	phone, err := phoneModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		phone = &db.Phone{
			Phone: "",
		}
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	username, err := userNameModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		username = &db.Username{
			Username: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	userNameString, err := base64.StdEncoding.DecodeString(username.Username.String)
	if err != nil {
		userNameString = []byte("")
	}

	nickname, err := nickNameModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		nickname = &db.Nickname{
			Nickname: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	header, err := headerModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		header = &db.Header{
			Header: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	email, err := emailModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		email = &db.Email{
			Email: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	idcard, err := idcardModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		idcard = nil
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	company, err := companyModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		company = nil
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	password, err := passwordModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		password = nil
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	secondfa, err := secondfaModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		secondfa = nil
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	wechat, err := wechatModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		wechat = &db.Wechat{
			UnionId: sql.NullString{
				String: "",
			},
			Nickname: sql.NullString{
				String: "",
			},
			Headimgurl: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	wxrobot, err := wxrobotModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		wxrobot = nil
	} else if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	inviteCount, err := userModel.GetUserInviteCount(ctx, user.Id)
	if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	ctrl, err := ctrlModel.FindByUserID(ctx, user.Id)
	if err != nil {
		return InfoDataEasy{}, errors.WarpQuick(err)
	}

	userEasy := types.UserEasy{
		UID:            user.Uid,
		InviteCount:    inviteCount,
		RoleID:         user.RoleId,
		RoleName:       userRole.Name,
		RoleSign:       userRole.Sign,
		Signin:         user.Signin,
		TokenExpire:    user.TokenExpiration,
		Phone:          phone.Phone,
		UserName:       string(userNameString),
		NickName:       nickname.Nickname.String,
		Header:         header.Header.String,
		Email:          email.Email.String,
		WeChatNickName: wechat.Nickname.String,
		WeChatHeader:   wechat.Headimgurl.String,
		Status:         db.UserStatusMap[user.Status],
		CreateAt:       user.CreateAt.Unix(),
	}

	if idcard == nil {
		userEasy.UserRealName = ""
	} else {
		userEasy.UserRealName = idcard.UserName
		if company != nil {
			userEasy.CompanyName = company.CompanyName
		} else {
			userEasy.CompanyName = ""
		}
	}

	userData := types.UserData{
		Has2FA:                  secondfa != nil && secondfa.Secret.Valid,
		HasEmail:                email.Email.Valid,
		HasPassword:             password != nil && password.PasswordHash.Valid,
		HasWeChat:               wechat.OpenId.Valid,
		HasFuwuhao:              wechat.Fuwuhao.Valid,
		HasUnionID:              wechat.UnionId.Valid,
		HasWxrobot:              wxrobot != nil && wxrobot.Webhook.Valid,
		HasVerified:             idcard != nil,
		IsCompany:               idcard != nil && company != nil,
		VerifiedPhone:           "", // 下面设置
		HasUserOriginal:         idcard != nil && idcard.IdcardKey.Valid,
		HasUserFaceCheck:        idcard != nil && idcard.FaceCheckId.Valid,
		HasCompanyOriginal:      idcard != nil && company != nil && company.IdcardKey.Valid,
		HasLegalPersonFaceCheck: idcard != nil && company != nil && company.FaceCheckId.Valid,
		AllowPhone:              ctrl.AllowPhone,
		AllowPassword:           ctrl.AllowPassword,
		AllowEmail:              ctrl.AllowEmail,
		AllowWeChat:             ctrl.AllowWechat,
		AllowSecondFA:           ctrl.Allow2Fa,
	}

	if idcard != nil && idcard.Phone.Valid {
		userData.VerifiedPhone = idcard.Phone.String
	}

	userInfoEasy := types.UserInfoEsay{}
	if idcard != nil {
		userInfoEasy.HasVerified = true
		userInfoEasy.UserName = idcard.UserName

		if company != nil {
			userInfoEasy.IsCompany = true
			userInfoEasy.LegalPersonName = company.LegalPersonName
			userInfoEasy.CompanyName = company.CompanyName
		} else {
			userInfoEasy.IsCompany = false
		}
	} else {
		userInfoEasy.HasVerified = false
	}

	return InfoDataEasy{
		User:     userEasy,
		Data:     userData,
		InfoEasy: userInfoEasy,
	}, nil
}

func GetUserInfoWebsite(ctx context.Context, user *db.User) (InfoDataWebsite, errors.WTError) {
	userModel := db.NewUserModel(mysql.MySQLConn)
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	userNameModel := db.NewUsernameModel(mysql.MySQLConn)
	nickNameModel := db.NewNicknameModel(mysql.MySQLConn)
	headerModel := db.NewHeaderModel(mysql.MySQLConn)
	emailModel := db.NewEmailModel(mysql.MySQLConn)
	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	companyModel := db.NewCompanyModel(mysql.MySQLConn)
	passwordModel := db.NewPasswordModel(mysql.MySQLConn)
	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	wxrobotModel := db.NewWxrobotModel(mysql.MySQLConn)
	ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

	userRole := action.GetRole(user.RoleId, user.IsAdmin)

	phone, err := phoneModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		phone = &db.Phone{
			Phone: "",
		}
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	username, err := userNameModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		username = &db.Username{
			Username: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	userNameString, err := base64.StdEncoding.DecodeString(username.Username.String)
	if err != nil {
		userNameString = []byte("")
	}

	nickname, err := nickNameModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		nickname = &db.Nickname{
			Nickname: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	header, err := headerModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		header = &db.Header{
			Header: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	email, err := emailModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		email = &db.Email{
			Email: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	idcard, err := idcardModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		idcard = nil
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	company, err := companyModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		company = nil
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	password, err := passwordModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		password = nil
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	secondfa, err := secondfaModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		secondfa = nil
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	wechat, err := wechatModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		wechat = &db.Wechat{
			UnionId: sql.NullString{
				String: "",
			},
			Nickname: sql.NullString{
				String: "",
			},
			Headimgurl: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	wxrobot, err := wxrobotModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		wxrobot = nil
	} else if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	inviteCount, err := userModel.GetUserInviteCount(ctx, user.Id)
	if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	ctrl, err := ctrlModel.FindByUserID(ctx, user.Id)
	if err != nil {
		return InfoDataWebsite{}, errors.WarpQuick(err)
	}

	userEasy := types.WebsiteUserEasy{
		UID:            user.Uid,
		InviteCount:    inviteCount,
		RoleID:         user.RoleId,
		RoleName:       userRole.Name,
		RoleSign:       userRole.Sign,
		Signin:         user.Signin,
		TokenExpire:    user.TokenExpiration,
		Phone:          phone.Phone,
		UserName:       string(userNameString),
		NickName:       nickname.Nickname.String,
		Header:         header.Header.String,
		Email:          email.Email.String,
		WeChatNickName: wechat.Nickname.String,
		WeChatHeader:   wechat.Headimgurl.String,
		UnionID:        wechat.UnionId.String,
		Status:         db.UserStatusMap[user.Status],
		CreateAt:       user.CreateAt.Unix(),
	}

	if idcard == nil {
		userEasy.UserRealName = ""
	} else {
		userEasy.UserRealName = idcard.UserName
		if company != nil {
			userEasy.CompanyName = company.CompanyName
		} else {
			userEasy.CompanyName = ""
		}
	}

	userData := types.UserData{
		Has2FA:                  secondfa != nil && secondfa.Secret.Valid,
		HasEmail:                email.Email.Valid,
		HasPassword:             password != nil && password.PasswordHash.Valid,
		HasWeChat:               wechat.OpenId.Valid,
		HasFuwuhao:              wechat.Fuwuhao.Valid,
		HasUnionID:              wechat.UnionId.Valid,
		HasWxrobot:              wxrobot != nil && wxrobot.Webhook.Valid,
		HasVerified:             idcard != nil,
		IsCompany:               idcard != nil && company != nil,
		VerifiedPhone:           "", // 下面设置
		HasUserOriginal:         idcard != nil && idcard.IdcardKey.Valid,
		HasUserFaceCheck:        idcard != nil && idcard.FaceCheckId.Valid,
		HasCompanyOriginal:      idcard != nil && company != nil && company.IdcardKey.Valid,
		HasLegalPersonFaceCheck: idcard != nil && company != nil && company.FaceCheckId.Valid,
		AllowPhone:              ctrl.AllowPhone,
		AllowPassword:           ctrl.AllowPassword,
		AllowEmail:              ctrl.AllowEmail,
		AllowWeChat:             ctrl.AllowWechat,
		AllowSecondFA:           ctrl.Allow2Fa,
	}

	if idcard != nil && idcard.Phone.Valid {
		userData.VerifiedPhone = idcard.Phone.String
	}

	userInfoEasy := types.UserInfoEsay{}
	if idcard != nil {
		userInfoEasy.HasVerified = true
		userInfoEasy.UserName = idcard.UserName

		if company != nil {
			userInfoEasy.IsCompany = true
			userInfoEasy.LegalPersonName = company.LegalPersonName
			userInfoEasy.CompanyName = company.CompanyName
		} else {
			userInfoEasy.IsCompany = false
		}
	} else {
		userInfoEasy.HasVerified = false
	}

	return InfoDataWebsite{
		User:     userEasy,
		Data:     userData,
		InfoEasy: userInfoEasy,
	}, nil
}
