package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/peers/wstype"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/center/userwstype"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
)

type SrcTitle struct {
	Name   string `json:"name"`   // 姓名，公司名
	TaxID  string `json:"taxID"`  // 税号，身份证号
	BankID string `json:"bankID"` // 银行卡号
	Bank   string `json:"bank"`   // 开户行
}

type SrcGetInfoData struct {
	User    SrcUserEasy `json:"user"`
	Info    SrcUserInfo `json:"info"`
	Data    SrcUserData `json:"data"`
	Title   SrcTitle    `json:"title"`
	Address SrcAddress  `json:"address"`
}

type SrcAddress struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email"`
	Province string   `json:"province"`
	City     string   `json:"city"`
	District string   `json:"district"`
	Address  string   `json:"address"`
	Area     []string `json:"area"`
}

type SrcUserEasy struct {
	UID            string `json:"id"`
	Phone          string `json:"phone"`
	UserName       string `json:"userName"`
	NickName       string `json:"nickname"`
	Header         string `json:"header"`
	Email          string `json:"email"`
	UserRealName   string `json:"userRealName"`
	CompanyName    string `json:"companyName"`
	WeChatNickName string `json:"wechatNickName"`
	WeChatHeader   string `json:"wechatHeader"`
	UnionID        string `json:"unionID"`
	Signin         bool   `json:"signin"`
	Status         string `json:"status"`
	InviteCount    int64  `json:"inviteCount"`
	TokenExpire    int64  `json:"tokenExpire"`
	CreateAt       int64  `json:"createAt"`
}

type SrcUserInfo struct {
	HasVerified     bool   `json:"hasVerified"`
	UserName        string `json:"userName"`
	UserIDCard      string `json:"userIDCard"`
	VerifiedPhone   string `json:"verifiedPhone"`
	IsCompany       bool   `json:"isCompany"`
	LegalPersonName string `json:"legalPersonName"`
	CompanyName     string `json:"companyName"`
}

type SrcUserSecret struct {
	LegalPersonIDCard string `json:"legalPersonIdCard"`
	CompanyID         string `json:"companyID"`
}

type SrcUserInfoEsay struct {
	HasVerified     bool   `json:"hasVerified"`
	UserName        string `json:"userName"`
	IsCompany       bool   `json:"isCompany"`
	LegalPersonName string `json:"legalPersonName"`
	CompanyName     string `json:"companyName"`
}

type SrcUserData struct {
	HasPassword             bool   `json:"hasPassword"`
	HasEmail                bool   `json:"hasEmail"`
	Has2FA                  bool   `json:"has2FA"`
	HasWeChat               bool   `json:"hasWeChat"`
	HasWxrobot              bool   `json:"hasWxrobot"`
	HasUnionID              bool   `json:"hasUnionId"`
	HasFuwuhao              bool   `json:"hasFuwuhao"`
	HasVerified             bool   `json:"hasVerified"`
	IsCompany               bool   `json:"isCompany"`
	HasUserOriginal         bool   `json:"hasUserOriginal"`
	HasUserFaceCheck        bool   `json:"hasUserFaceCheck"`
	HasCompanyOriginal      bool   `json:"hasCompanyOriginal"`
	HasLegalPersonFaceCheck bool   `json:"hasLegalPersonFaceCheck"`
	VerifiedPhone           string `json:"verifiedPhone"`
	AllowPhone              bool   `json:"allowPhone"`
	AllowEmail              bool   `json:"allowEmail"`
	AllowPassword           bool   `json:"allowPassword"`
	AllowWeChat             bool   `json:"allowWeChat"`
	AllowSecondFA           bool   `json:"allowSecondFA"`
}

type SrcRoleMenu struct {
	ID             int64  `json:"id"`
	Sort           int64  `json:"sort"`
	FatherID       int64  `json:"parentID,optional"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Title          string `json:"title"`
	Icon           string `json:"icon"`
	Redirect       string `json:"redirect"`
	Superior       string `json:"menuSuperior"`
	Category       int64  `json:"menuCategory"`
	Component      string `json:"component"`
	ComponentAlias string `json:"componentAlias"`
	MetaLink       string `json:"metaIsLink"`
	Type           int64  `json:"menuType"`
	IsLink         bool   `json:"isLink"`
	IsHide         bool   `json:"isHide"`
	IsKeepalive    bool   `json:"isKeepAlive"`
	IsAffix        bool   `json:"isAffix"`
	IsIframe       bool   `json:"isIframe"`
	BtnPower       string `json:"btnPower"`
}

type SrcRole struct {
	RoleID               int64         `json:"roleID"`
	Describe             string        `json:"describe"` // 描述
	Name                 string        `json:"name"`
	Sign                 string        `json:"sign"`
	NotDelete            bool          `json:"notDelete"`
	NotChangeSign        bool          `json:"notChangeSign"`
	NotChangePermissions bool          `json:"notChangePermissions"`
	Status               int64         `json:"status"`
	CreateAt             int64         `json:"createAt"`
	Sort                 int64         `json:"sort"`
	Permissions          int64         `json:"permission"`
	Menus                []SrcRoleMenu `json:"menus"`
}

func UpdateUser(userID int64, mysql sqlx.SqlConn, ch chan websocket.WSMessage) {
	var lst []chan websocket.WSMessage

	websocket.UserConnMapMutex.Lock()
	defer websocket.UserConnMapMutex.Unlock()

	if ch == nil {
		lst = websocket.UserConnMap[userID]
	} else {
		lst = []chan websocket.WSMessage{ch}
	}

	userM := NewUserModel(mysql)
	phoneM := NewPhoneModel(mysql)
	userNameM := NewUsernameModel(mysql)
	nickNameM := NewNicknameModel(mysql)
	headerM := NewHeaderModel(mysql)
	emailM := NewEmailModel(mysql)
	idcardM := NewIdcardModel(mysql)
	companyM := NewCompanyModel(mysql)
	passwordM := NewPasswordModel(mysql)
	secondfaM := NewSecondfaModel(mysql)
	wechatM := NewWechatModel(mysql)
	titleM := NewTitleModel(mysql)
	addressM := NewAddressModel(mysql)
	wxrobotM := NewWxrobotModel(mysql)
	ctrlM := NewLoginControllerModel(mysql)

	user, err := userM.FindOneByIDWithoutDelete(context.Background(), userID)
	if errors.Is(err, ErrNotFound) || err != nil {
		return
	}

	phone, err := phoneM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		phone = &Phone{
			Phone: "",
		}
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	username, err := userNameM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		username = &Username{
			Username: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	userNameString, err := base64.StdEncoding.DecodeString(username.Username.String)
	if err != nil {
		userNameString = []byte("")
	}

	nickname, err := nickNameM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		nickname = &Nickname{
			Nickname: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	header, err := headerM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		header = &Header{
			Header: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	email, err := emailM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		email = &Email{
			Email: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	idcard, err := idcardM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		idcard = nil
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	company, err := companyM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		company = nil
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	password, err := passwordM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		password = nil
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	secondfa, err := secondfaM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		secondfa = nil
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	wechat, err := wechatM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		wechat = &Wechat{
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
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	wxrobot, err := wxrobotM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		wxrobot = nil
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	title, err := titleM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		title = nil
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	address, err := addressM.FindByUserID(context.Background(), userID)
	if errors.Is(err, ErrNotFound) {
		address = nil
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	inviteCount, err := userM.GetUserInviteCount(context.Background(), userID)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	ctrl, err := ctrlM.FindByUserID(context.Background(), userID)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return
	}

	userEasy := SrcUserEasy{
		UID:            user.Uid,
		InviteCount:    inviteCount,
		Signin:         user.Signin,
		TokenExpire:    user.TokenExpiration,
		Phone:          phone.Phone,
		UserName:       string(userNameString),
		NickName:       nickname.Nickname.String,
		Header:         header.Header.String,
		Email:          email.Email.String,
		UnionID:        wechat.UnionId.String,
		WeChatNickName: wechat.Nickname.String,
		WeChatHeader:   wechat.Headimgurl.String,
		Status:         UserStatusMap[user.Status],
		CreateAt:       user.CreateAt.Unix(),
	}

	userInfo := SrcUserInfo{}
	userSecret := SrcUserSecret{}
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

			userSecret.CompanyID = company.CompanyId
			userSecret.LegalPersonIDCard = company.LegalPersonIdCard
		} else {
			userEasy.CompanyName = ""
			userInfo.IsCompany = false
		}
	}

	userData := SrcUserData{
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

	userTitle := SrcTitle{}
	if title != nil {
		userTitle.Name = title.Name.String
		userTitle.TaxID = title.TaxId.String
		userTitle.BankID = title.BankId.String
		userTitle.Bank = title.Bank.String
	}

	userAddress := SrcAddress{}
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

	res := SrcGetInfoData{
		User:    userEasy,
		Info:    userInfo,
		Data:    userData,
		Address: userAddress,
		Title:   userTitle,
	}

	msg := websocket.WSMessage{
		Code:        userwstype.UpdateUserInfo,
		Data:        res,
		Secret:      userSecret,
		SecretToken: []int{jwt.UserRootToken, jwt.UserHighAuthorityRootToken, jwt.UserRootFatherToken},
	}

	if ch == nil {
		websocket.WritePeersMessage(wstype.PeersUpdateUserInfo, struct {
			UserID int64 `json:"userID"`
		}{UserID: userID}, msg)
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func UpdateUserByMsg(userID int64, msg websocket.WSMessage) {
	websocket.UserConnMapMutex.Lock()
	defer websocket.UserConnMapMutex.Unlock()

	lst, ok := websocket.UserConnMap[userID]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

type SrcAnnouncement struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	StartAt int64  `json:"startAt"`
	StopAt  int64  `json:"stopAt"`
	Sort    int64  `json:"sort"`
}

type SrcDeleteAnnouncement struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func UpdateAnnouncementByMsg(msg websocket.WSMessage) {
	websocket.AnnouncementConnListMutex.Lock()
	defer websocket.AnnouncementConnListMutex.Unlock()

	for _, i := range websocket.AnnouncementConnList {
		websocket.WriteMessage(i, msg)
	}
}

func UpdateAnnouncement(announcement *Announcement) {
	websocket.AnnouncementConnListMutex.Lock()
	defer websocket.AnnouncementConnListMutex.Unlock()

	msg := websocket.WSMessage{
		Code: userwstype.UpdateAnnouncement,
		Data: SrcAnnouncement{
			ID:      announcement.Id,
			Title:   announcement.Title,
			Content: announcement.Content,
			StartAt: announcement.StartAt.Unix(),
			StopAt:  announcement.StopAt.Unix(),
			Sort:    announcement.Sort,
		},
	}

	websocket.WritePeersMessage(wstype.PeersUpdateAnnouncement, struct{}{}, msg)

	for _, i := range websocket.AnnouncementConnList {
		websocket.WriteMessage(i, msg)
	}
}

func NewAnnouncement(announcement *Announcement) {
	websocket.AnnouncementConnListMutex.Lock()
	defer websocket.AnnouncementConnListMutex.Unlock()

	msg := websocket.WSMessage{
		Code: userwstype.NewAnnouncement,
		Data: SrcAnnouncement{
			ID:      announcement.Id,
			Title:   announcement.Title,
			Content: announcement.Content,
			StartAt: announcement.StartAt.Unix(),
			StopAt:  announcement.StopAt.Unix(),
			Sort:    announcement.Sort,
		},
	}

	websocket.WritePeersMessage(wstype.PeersUpdateAnnouncement, struct{}{}, msg)

	for _, i := range websocket.AnnouncementConnList {
		websocket.WriteMessage(i, msg)
	}
}

func DeleteAnnouncement(announcement *Announcement) {
	websocket.AnnouncementConnListMutex.Lock()
	defer websocket.AnnouncementConnListMutex.Unlock()

	msg := websocket.WSMessage{
		Code: userwstype.DeleteAnnouncement,
		Data: SrcDeleteAnnouncement{
			ID:    announcement.Id,
			Title: announcement.Title,
		},
	}

	websocket.WritePeersMessage(wstype.PeersUpdateAnnouncement, struct{}{}, msg)

	for _, i := range websocket.AnnouncementConnList {
		websocket.WriteMessage(i, msg)
	}
}

type SrcMessage struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Sender     string `json:"sender"`
	SenderLink string `json:"senderLink"`
	CreateAt   int64  `json:"createAt"`
	ReadAt     int64  `json:"readAt"`
	SenderID   int64  `json:"-"`
}

type SrcReadMessage struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	CreateAt int64  `json:"createAt"`
	ReadAt   int64  `json:"readAt"`
	SenderID int64  `json:"-"`
}

func UpdateMessageByMsg(userID int64, msg websocket.WSMessage) {
	websocket.MessageConnMapMutex.Lock()
	defer websocket.MessageConnMapMutex.Unlock()

	lst, ok := websocket.MessageConnMap[userID]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func UpdateMessage(msg *Message) {
	websocket.MessageConnMapMutex.Lock()
	defer websocket.MessageConnMapMutex.Unlock()

	lst, ok := websocket.MessageConnMap[msg.UserId]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	readAt := int64(0)
	if msg.ReadAt.Valid {
		readAt = msg.ReadAt.Time.Unix()
	}

	m := websocket.WSMessage{
		Code: userwstype.UpdateMessage,
		Data: SrcMessage{
			ID:         msg.Id,
			Title:      msg.Title,
			Content:    msg.Content,
			Sender:     msg.Sender,
			SenderLink: msg.SenderLink.String,
			CreateAt:   msg.CreateAt.Unix(),
			ReadAt:     readAt,
			SenderID:   msg.SenderId,
		},
		WebID: msg.SenderId,
	}

	websocket.WritePeersMessage(wstype.PeersUpdateMessage, struct {
		UserID int64 `json:"userID"`
	}{UserID: msg.UserId}, m)

	for _, i := range lst {
		websocket.WriteMessage(i, m)
	}
}

func ReadMessage(msg *Message) {
	if !msg.ReadAt.Valid {
		return
	}

	websocket.MessageConnMapMutex.Lock()
	defer websocket.MessageConnMapMutex.Unlock()

	lst, ok := websocket.MessageConnMap[msg.UserId]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	m := websocket.WSMessage{
		Code: userwstype.ReadMessage,
		Data: SrcReadMessage{
			ID:       msg.Id,
			Title:    msg.Title,
			CreateAt: msg.CreateAt.Unix(),
			ReadAt:   msg.ReadAt.Time.Unix(),
			SenderID: msg.SenderId,
		},
		WebID: msg.SenderId,
	}

	websocket.WritePeersMessage(wstype.PeersUpdateMessage, struct {
		UserID int64 `json:"userID"`
	}{UserID: msg.UserId}, m)

	for _, i := range lst {
		websocket.WriteMessage(i, m)
	}
}

type SrcWorkOrderCommunicate struct {
	ID       int64                         `json:"id"`
	OrderID  string                        `json:"orderID"`
	Content  string                        `json:"content"`
	From     int64                         `json:"from"`
	CreateAt int64                         `json:"createAt"`
	File     []SrcWorkOrderCommunicateFile `json:"file"`
	FromID   int64                         `json:"-"`
}

type SrcWorkOrderCommunicateFile struct {
	Fid string `json:"fid"`
}

type SrcWorkOrder struct {
	OrderID     string `json:"orderID"`
	Title       string `json:"title"`
	From        string `json:"from"`
	Status      int64  `json:"status"`
	CreateAt    int64  `json:"createAt"`
	LastReplyAt int64  `json:"lastReplyAt"`
	FinishAt    int64  `json:"finishAt"`
	FromID      int64  `json:"-"`
}

func NewWorkOrderCommunicate(wc *WorkOrderCommunicate, orderID string, fromID int64, mysql sqlx.SqlConn) {
	websocket.OrderConnMapMutex.Lock()
	defer websocket.OrderConnMapMutex.Unlock()

	lst, ok := websocket.OrderConnMap[orderID]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	workOrderCommunicateFileM := NewWorkOrderCommunicateFileModel(mysql)

	fileList, err := workOrderCommunicateFileM.GetList(context.Background(), wc.Id)
	if err != nil {
		logger.Logger.Error("mysql errors: %s", err.Error())
		fileList = []WorkOrderCommunicateFile{}
	}

	respFileList := make([]SrcWorkOrderCommunicateFile, 0, len(fileList))
	for _, f := range fileList {
		respFileList = append(respFileList, SrcWorkOrderCommunicateFile{
			Fid: f.Fid,
		})
	}

	msg := websocket.WSMessage{
		Code: userwstype.NewOrderReply,
		Data: SrcWorkOrderCommunicate{
			ID:       wc.Id,
			OrderID:  orderID,
			Content:  wc.Content,
			From:     wc.From,
			CreateAt: wc.CreateAt.Unix(),
			File:     respFileList,
			FromID:   fromID,
		},
		WebID: fromID,
	}

	websocket.WritePeersMessage(wstype.PeersUpdateOrder, struct {
		OrderID string `json:"orderID"`
	}{OrderID: orderID}, msg)

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func UpdateWorkOrder(w *WorkOrder) {
	websocket.OrderConnMapMutex.Lock()
	defer websocket.OrderConnMapMutex.Unlock()

	lst, ok := websocket.OrderConnMap[w.Uid]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	replyAt := int64(0)
	if w.LastReplyAt.Valid {
		replyAt = w.LastReplyAt.Time.Unix()
	}

	finishAt := int64(0)
	if w.FinishAt.Valid {
		finishAt = w.FinishAt.Time.Unix()
	}

	msg := websocket.WSMessage{
		Code: userwstype.UpdateOrder,
		Data: SrcWorkOrder{
			OrderID:     w.Uid,
			Title:       w.Title,
			From:        w.From,
			Status:      w.Status,
			CreateAt:    w.CreateAt.Unix(),
			FinishAt:    finishAt,
			LastReplyAt: replyAt,
			FromID:      w.FromId,
		},
		WebID: w.FromId,
	}

	websocket.WritePeersMessage(wstype.PeersUpdateOrder, struct {
		OrderID string `json:"orderID"`
	}{OrderID: w.Uid}, msg)

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func UpdateWorkOrderByMsg(orderID string, msg websocket.WSMessage) {
	websocket.OrderConnMapMutex.Lock()
	defer websocket.OrderConnMapMutex.Unlock()

	lst, ok := websocket.OrderConnMap[orderID]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}
