package action

import (
	"context"
	"encoding/base64"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	errors "github.com/wuntsong-org/wterrors"
)

var UserEasyNotFound = errors.NewClass("user easy not found")

func GetUserEasy(ctx context.Context, userID int64, userUID string) (types.UserEasy, errors.WTError) {
	var user *db.UserEasy
	var err error

	userModel := db.NewUserModel(mysql.MySQLConn)

	if userID != 0 {
		user, err = userModel.FindUserEasyByIDWithoutDelete(ctx, userID)
	} else if len(userUID) != 0 {
		user, err = userModel.FindUserEasyByUidWithoutDelete(ctx, userUID)
	} else {
		return types.UserEasy{}, UserEasyNotFound.New()
	}
	if errors.Is(err, db.ErrNotFound) {
		return types.UserEasy{}, UserEasyNotFound.New()
	} else if err != nil {
		return types.UserEasy{}, errors.WarpQuick(err)
	}

	return GetUserEasyOther(ctx, user)
}

func GetUserMoreEasy(ctx context.Context, userID int64, userUID string) (types.UserMoreEasy, errors.WTError) {
	var user *db.UserEasy
	var err error

	userModel := db.NewUserModel(mysql.MySQLConn)

	if userID != 0 {
		user, err = userModel.FindUserEasyByIDWithoutDelete(ctx, userID)
	} else if len(userUID) != 0 {
		user, err = userModel.FindUserEasyByUidWithoutDelete(ctx, userUID)
	} else {
		return types.UserMoreEasy{}, UserEasyNotFound.New()
	}
	if errors.Is(err, db.ErrNotFound) {
		return types.UserMoreEasy{}, UserEasyNotFound.New()
	} else if err != nil {
		return types.UserMoreEasy{}, errors.WarpQuick(err)
	}

	return GetUserMoreEasyOther(ctx, user)
}

func GetUserLessEasy(ctx context.Context, userID int64, userUID string) (types.UserLessEasy, errors.WTError) {
	var user *db.UserEasy

	userModel := db.NewUserModel(mysql.MySQLConn)

	var err error
	if userID != 0 {
		user, err = userModel.FindUserEasyByIDWithoutDelete(ctx, userID)
	} else if len(userUID) != 0 {
		user, err = userModel.FindUserEasyByUidWithoutDelete(ctx, userUID)
	} else {
		return types.UserLessEasy{}, UserEasyNotFound.New()
	}
	if errors.Is(err, db.ErrNotFound) {
		return types.UserLessEasy{}, UserEasyNotFound.New()
	} else if err != nil {
		return types.UserLessEasy{}, errors.WarpQuick(err)
	}

	return GetUserLessEasyOther(ctx, user)
}

func GetUserMoreEasyOther(ctx context.Context, user *db.UserEasy) (types.UserMoreEasy, errors.WTError) {
	r := GetRole(user.RoleID, user.IsAdmin)
	userName, err := base64.StdEncoding.DecodeString(user.UserName.String)
	if err != nil {
		userName = []byte("")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	inviteCount, err := userModel.GetUserInviteCount(ctx, user.ID)
	if err != nil {
		return types.UserMoreEasy{}, errors.WarpQuick(err)
	}

	return types.UserMoreEasy{
		UID:         user.UID,
		InviteCount: inviteCount,
		RoleID:      r.ID,
		RoleName:    r.Name,
		RoleSign:    r.Sign,
		UserName:    string(userName),
		Phone:       user.Phone.String,
		NickName:    user.NickName.String,
		Header:      user.Header.String,
		Email:       user.Email.String,
		Status:      db.UserStatusMap[user.Status],
		CreateAt:    user.CreateAt.Unix(),
	}, nil
}

func GetUserLessEasyOther(ctx context.Context, user *db.UserEasy) (types.UserLessEasy, errors.WTError) {
	userModel := db.NewUserModel(mysql.MySQLConn)
	inviteCount, err := userModel.GetUserInviteCount(ctx, user.ID)
	if err != nil {
		return types.UserLessEasy{}, errors.WarpQuick(err)
	}

	return types.UserLessEasy{
		UID:         user.UID,
		InviteCount: inviteCount,
		Phone:       user.Phone.String,
		Status:      db.UserStatusMap[user.Status],
		CreateAt:    user.CreateAt.Unix(),
	}, nil
}

func GetUserEasyOther(ctx context.Context, user *db.UserEasy) (types.UserEasy, errors.WTError) {
	r := GetRole(user.RoleID, user.IsAdmin)
	userName, err := base64.StdEncoding.DecodeString(user.UserName.String)
	if err != nil {
		userName = []byte("")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	inviteCount, err := userModel.GetUserInviteCount(ctx, user.ID)
	if err != nil {
		return types.UserEasy{}, errors.WarpQuick(err)
	}

	return types.UserEasy{
		TokenExpire:    user.TokenExpire,
		Signin:         user.SignIn,
		UID:            user.UID,
		InviteCount:    inviteCount,
		RoleID:         r.ID,
		RoleName:       r.Name,
		RoleSign:       r.Sign,
		UserName:       string(userName),
		Phone:          user.Phone.String,
		NickName:       user.NickName.String,
		Header:         user.Header.String,
		Email:          user.Email.String,
		UserRealName:   user.UserRealName.String,
		CompanyName:    user.CompanyName.String,
		WeChatNickName: user.WeChatNickName.String,
		WeChatHeader:   user.WeChatHeader.String,
		Status:         db.UserStatusMap[user.Status],
		CreateAt:       user.CreateAt.Unix(),
	}, nil
}

func GetUncleUserEasyOther(ctx context.Context, user *db.UncleUserEasy) (types.UserEasy, errors.WTError) {
	r := GetRole(user.RoleID, user.IsAdmin)
	userName, err := base64.StdEncoding.DecodeString(user.UserName.String)
	if err != nil {
		userName = []byte("")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	inviteCount, err := userModel.GetUserInviteCount(ctx, user.ID)
	if err != nil {
		return types.UserEasy{}, errors.WarpQuick(err)
	}

	return types.UserEasy{
		TokenExpire:    user.TokenExpire,
		Signin:         user.SignIn,
		UID:            user.UID,
		InviteCount:    inviteCount,
		RoleID:         r.ID,
		RoleName:       r.Name,
		RoleSign:       r.Sign,
		UserName:       string(userName),
		Phone:          user.Phone.String,
		NickName:       user.NickName.String,
		Header:         user.Header.String,
		Email:          user.Email.String,
		UserRealName:   user.UserRealName.String,
		CompanyName:    user.CompanyName.String,
		WeChatNickName: user.WeChatNickName.String,
		WeChatHeader:   user.WeChatHeader.String,
		Status:         db.UserStatusMap[user.Status],
		CreateAt:       user.CreateAt.Unix(),
	}, nil
}
