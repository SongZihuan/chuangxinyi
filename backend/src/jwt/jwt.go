package jwt

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/peers/wstype"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/center/userwstype"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/record"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/wuntsong-org/wterrors"
	"strconv"
	"strings"
	"time"
)

func InitJWT() errors.WTError {
	if config.BackendConfig.Jwt.Login.ExpiresSecond == 0 {
		return errors.Errorf("jwt login expires second must be given")
	}

	if config.BackendConfig.Jwt.Email.ExpiresSecond == 0 {
		return errors.Errorf("jwt email expires second must be given")
	}

	if config.BackendConfig.Jwt.Phone.ExpiresSecond == 0 {
		return errors.Errorf("jwt phone expires second must be given")
	}

	if config.BackendConfig.Jwt.SecondFA.ExpiresSecond == 0 {
		return errors.Errorf("jwt 2FA expires second must be given")
	}

	if config.BackendConfig.Jwt.User.ExpiresSecond == 0 {
		return errors.Errorf("jwt user expires second must be given")
	}

	if config.BackendConfig.Jwt.WeChat.ExpiresSecond == 0 {
		return errors.Errorf("jwt wechat expires second must be given")
	}

	if config.BackendConfig.Jwt.IDCard.ExpiresSecond == 0 {
		return errors.Errorf("jwt idcard expires second must be given")
	}

	if config.BackendConfig.Jwt.Company.ExpiresSecond == 0 {
		return errors.Errorf("jwt company expires second must be given")
	}

	if config.BackendConfig.Jwt.Face.ExpiresSecond == 0 {
		return errors.Errorf("jwt face expires second must be given")
	}

	if config.BackendConfig.Jwt.CheckSecondFA.ExpiresSecond == 0 {
		return errors.Errorf("jwt check 2fa expires second must be given")
	}

	if config.BackendConfig.Jwt.Delete.ExpiresSecond == 0 {
		return errors.Errorf("jwt delete expires second must be given")
	}

	if config.BackendConfig.Jwt.ExpireSecond == 0 {
		return errors.Errorf("jwt expire second must be given")
	}

	err := DeleteAllActive()
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

type DefrayTokenData struct {
	OwnerID            int64  `json:"ownerID"`
	TradeID            string `json:"tradeID"`
	Subject            string `json:"subject"`    // 标题
	Price              int64  `json:"price"`      // 价格
	Quantity           int64  `json:"quantity"`   // 数量
	UnitPrice          int64  `json:"unitPrice"`  // 单价
	Describe           string `json:"describe"`   // 描述
	SupplierID         int64  `json:"supplierID"` // 供应商ID
	ReturnURL          string `json:"returnURL"`
	TimeExpire         int64  `json:"timeExpire"` // 结束时间
	InvitePre          int64  `json:"invitePre"`
	DistributionLevel1 int64  `json:"distributionLevel1"`
	DistributionLevel2 int64  `json:"distributionLevel2"`
	DistributionLevel3 int64  `json:"distributionLevel3"`
	CanWithdraw        bool   `json:"canWithdraw"`
	MustSelfDefray     bool   `json:"mustSelfDefray"`
	ReturnDayLimit     int64  `json:"returnDayLimit"`
}

func CreateDefrayTokenToken(data DefrayTokenData) (string, errors.WTError) {
	// 预留多30s有效期
	token, _, err := createToken(data, config.BackendConfig.Jwt.Defray.Subject, config.BackendConfig.Jwt.Defray.Issuer, time.Unix(data.TimeExpire+30, 0), 0)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserDefrayToken(tokenString string) (DefrayTokenData, errors.WTError) {
	var data DefrayTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.Defray.Subject, config.BackendConfig.Jwt.Defray.Issuer)
	if err != nil {
		return DefrayTokenData{}, err
	}

	return data, nil
}

func DeleteDefrayToken(tokenString string) {
	deleteToken(tokenString)
}

type SecondFAPassTokenData struct {
	UserID  string `json:"userID"`
	UA      string `json:"UA"`
	GeoCode string `json:"geoCode"`
}

func CreateSecondFAPassToken(userID string, ua string, geoCode string, hour int64) (string, errors.WTError) {
	token, _, err := createToken(SecondFAPassTokenData{
		UserID:  userID,
		UA:      ua,
		GeoCode: geoCode,
	}, config.BackendConfig.Jwt.PassSecondFA.Subject, config.BackendConfig.Jwt.PassSecondFA.Issuer, time.Time{}, hour*60*60)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserSecondFAPassToken(tokenString string) (SecondFAPassTokenData, errors.WTError) {
	var data SecondFAPassTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.PassSecondFA.Subject, config.BackendConfig.Jwt.PassSecondFA.Issuer)
	if err != nil {
		return SecondFAPassTokenData{}, err
	}

	return data, nil
}

const (
	LoginWS   = "ws"
	LoginGet  = "get"
	LoginPost = "post"
)

func SetLogin(ctx context.Context, token string, t string) {
	_, _, err := ParserUserToken(ctx, token)
	if err != nil {
		return
	}

	_ = redis.Set(ctx, fmt.Sprintf("usertoken:active:%s:%s", t, token), "1", time.Minute*5)
}

func DelLogin(ctx context.Context, token string, t string) {
	_ = redis.Del(ctx, fmt.Sprintf("usertoken:active:%s:%s", t, token))
}

func IsLoginToken(ctx context.Context, token string) bool {
	_, _, errParser := ParserUserToken(ctx, token)
	if errParser != nil {
		return false
	}

	var res int64

	res, err := redis.Exists(context.Background(), fmt.Sprintf("usertoken:active:%s:%s", LoginGet, token)).Result()
	if err == nil && res == 1 {
		return true
	}

	res, err = redis.Exists(context.Background(), fmt.Sprintf("usertoken:active:%s:%s", LoginPost, token)).Result()
	if err == nil && res == 1 {
		return true
	}

	res, err = redis.Exists(context.Background(), fmt.Sprintf("usertoken:active:%s:%s", LoginWS, token)).Result()
	if err == nil && res == 1 {
		return true
	}

	return false
}

func DeleteAllActive() errors.WTError {
	keys, err := redis.Keys(context.Background(), "usertoken:active:*:*").Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, k := range keys {
		_ = redis.Del(context.Background(), k)
	}

	return nil
}

const (
	TypeDeleteUserToken  = 1
	TypeDeleteLoginToken = 2
)

type DeleteTokenData struct {
	Token string `json:"token"`
	Type  int64  `json:"type"`
}

func CreateDeleteToken(tk string, tp int64) (string, errors.WTError) {
	token, _, err := createToken(DeleteTokenData{Token: tk, Type: tp}, config.BackendConfig.Jwt.Delete.Subject, config.BackendConfig.Jwt.Delete.Issuer, time.Time{}, config.BackendConfig.Jwt.Delete.ExpiresSecond)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParserDeleteToken(tokenString string) (DeleteTokenData, errors.WTError) {
	var data DeleteTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.Delete.Subject, config.BackendConfig.Jwt.Delete.Issuer)
	if err != nil {
		return DeleteTokenData{}, err
	}

	return data, nil
}

type PhoneTokenData struct {
	Phone string `json:"phone"`
	WebID int64  `json:"webID"`
}

func CreatePhoneToken(phone string, webID int64) (string, errors.WTError) {
	token, _, err := createToken(PhoneTokenData{Phone: phone, WebID: webID}, config.BackendConfig.Jwt.Phone.Subject, config.BackendConfig.Jwt.Phone.Issuer, time.Time{}, config.BackendConfig.Jwt.Phone.ExpiresSecond)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserPhoneToken(tokenString string) (PhoneTokenData, errors.WTError) {
	var data PhoneTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.Phone.Subject, config.BackendConfig.Jwt.Phone.Issuer)
	if err != nil {
		return PhoneTokenData{}, err
	}

	return data, nil
}

type EmailTokenData struct {
	Email string `json:"email"`
	WebID int64  `json:"webID"`
}

func CreateEmailToken(email string, webID int64) (string, errors.WTError) {
	token, _, err := createToken(EmailTokenData{Email: email, WebID: webID}, config.BackendConfig.Jwt.Email.Subject, config.BackendConfig.Jwt.Email.Issuer, time.Time{}, config.BackendConfig.Jwt.Email.ExpiresSecond)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserEmailToken(tokenString string) (EmailTokenData, errors.WTError) {
	var data EmailTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.Email.Subject, config.BackendConfig.Jwt.Email.Issuer)
	if err != nil {
		return EmailTokenData{}, err
	}

	return data, nil
}

type Check2FATokenData struct {
	UserID string `json:"userID"`
	WebID  int64  `json:"webID"`
}

func CreateCheck2FAToken(userID string, webID int64) (string, errors.WTError) {
	token, _, err := createToken(Check2FATokenData{UserID: userID, WebID: webID}, config.BackendConfig.Jwt.CheckSecondFA.Subject, config.BackendConfig.Jwt.CheckSecondFA.Issuer, time.Time{}, config.BackendConfig.Jwt.CheckSecondFA.ExpiresSecond)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserCheck2FAToken(tokenString string) (Check2FATokenData, errors.WTError) {
	var data Check2FATokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.CheckSecondFA.Subject, config.BackendConfig.Jwt.CheckSecondFA.Issuer)
	if err != nil {
		return Check2FATokenData{}, err
	}

	return data, nil
}

type WechatTokenData struct {
	AccessToken string `json:"accessToken"`
	OpenID      string `json:"openID"`
	UnionID     string `json:"unionID"`
	IsFuwuhao   bool   `json:"isFuwuhao"`
	// 不支持外站
}

func CreateWeChatToken(accessToken, openID, unionID string, isFuwuhao bool) (string, errors.WTError) {
	token, _, err := createToken(WechatTokenData{
		AccessToken: accessToken,
		OpenID:      openID,
		UnionID:     unionID,
		IsFuwuhao:   isFuwuhao,
	}, config.BackendConfig.Jwt.WeChat.Subject, config.BackendConfig.Jwt.WeChat.Issuer, time.Time{}, config.BackendConfig.Jwt.WeChat.ExpiresSecond)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserWeChatToken(tokenString string) (WechatTokenData, errors.WTError) {
	var data WechatTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.WeChat.Subject, config.BackendConfig.Jwt.WeChat.Issuer)
	if err != nil {
		return WechatTokenData{}, err
	}

	return data, nil
}

type IDCardTokenData struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	IDCard     string `json:"idcard"`
	IDCardBack string `json:"idcardBack"`
	WebID      int64  `json:"webID"`
}

func CreateIDCardToken(name, id, idcard string, idcardBack string, webID int64) (string, errors.WTError) {
	token, _, err := createToken(IDCardTokenData{
		Name:       name,
		ID:         id,
		IDCard:     idcard,
		IDCardBack: idcardBack,
		WebID:      webID,
	}, config.BackendConfig.Jwt.IDCard.Subject, config.BackendConfig.Jwt.IDCard.Issuer, time.Time{}, config.BackendConfig.Jwt.IDCard.ExpiresSecond)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserIDCardToken(tokenString string) (IDCardTokenData, errors.WTError) {
	var data IDCardTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.IDCard.Subject, config.BackendConfig.Jwt.IDCard.Issuer)
	if err != nil {
		return IDCardTokenData{}, err
	}

	return data, nil
}

type CompanyTokenData struct {
	Name            string `json:"name"`
	ID              string `json:"id"`
	LegalPersonName string `json:"legalPersonName"`
	LegalPersonID   string `json:"legalPersonID"`
	License         string `json:"license"`
	IDCard          string `json:"idcard"`
	IDCardBack      string `json:"idcardBack"`
	WebID           int64  `json:"webID"`
}

func CreateCompanyToken(name, id, legalPersonName, legalPersonID, license string, idcard string, idcardBack string, webID int64) (string, errors.WTError) {
	token, _, err := createToken(CompanyTokenData{
		Name:            name,
		ID:              id,
		LegalPersonName: legalPersonName,
		LegalPersonID:   legalPersonID,
		License:         license,
		IDCard:          idcard,
		IDCardBack:      idcardBack,
		WebID:           webID,
	}, config.BackendConfig.Jwt.Company.Subject, config.BackendConfig.Jwt.Company.Issuer, time.Time{}, config.BackendConfig.Jwt.Company.ExpiresSecond)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserCompanyToken(tokenString string) (CompanyTokenData, errors.WTError) {
	var data CompanyTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.Company.Subject, config.BackendConfig.Jwt.Company.Issuer)
	if err != nil {
		return CompanyTokenData{}, err
	}

	return data, nil
}

type FaceTokenData struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	CertifyID string `json:"certifyID"`
	CheckID   string `json:"checkID"`
	WebID     int64  `json:"webID"`
}

func CreateFaceToken(name, id, certifyID string, checkID string, webID int64) (string, errors.WTError) {
	token, _, err := createToken(FaceTokenData{
		Name:      name,
		ID:        id,
		CertifyID: certifyID,
		CheckID:   checkID,
		WebID:     webID,
	}, config.BackendConfig.Jwt.Face.Subject, config.BackendConfig.Jwt.Face.Issuer, time.Time{}, config.BackendConfig.Jwt.Face.ExpiresSecond)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserFaceToken(tokenString string) (FaceTokenData, errors.WTError) {
	var data FaceTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.Face.Subject, config.BackendConfig.Jwt.Face.Issuer)
	if err != nil {
		return FaceTokenData{}, err
	}

	return data, nil
}

func DeleteLoginToken(ctx context.Context, userID string, token string) errors.WTError {
	loginData, err := ParserLoginToken(ctx, token)
	if err != nil {
		return nil
	} else if loginData.UserID != userID {
		return errors.Errorf("bad token for user")
	}

	_ = DeleteUserToken(ctx, userID, loginData.UserToken)
	deleteToken(token)

	RemoteIP, ok := ctx.Value("X-Real-IP").(string)
	if !ok || len(RemoteIP) == 0 {
		RemoteIP = "0.0.0.0"
	}

	Geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok || len(Geo) == 0 {
		Geo = "未知"
	}

	r := record.GetRecordIfExists(ctx)

	recordData := struct {
		RequestsID string `json:"requestsID"`
		IP         string `json:"IP"`
		Geo        string `json:"geo"`
	}{RequestsID: r.RequestsID, IP: RemoteIP, Geo: Geo}

	recordByte, err := utils.JsonMarshal(recordData)
	if err != nil {
		return errors.WarpQuick(err)
	}

	recordModel := db.NewTokenRecordModel(mysql.MySQLConn)
	_, mysqlErr := recordModel.Insert(ctx, &db.TokenRecord{
		TokenType: db.LoginToken,
		Token:     token,
		Type:      db.TokenDelete,
		Data:      string(recordByte),
	})
	if err != nil {
		return errors.WarpQuick(mysqlErr)
	}

	return nil
}

func DeleteAllWebsiteLoginToken(ctx context.Context, userID string, webID int64) errors.WTError {
	res, err := redis.Keys(ctx, fmt.Sprintf("logintoken:%s:%d:*", userID, webID)).Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, k := range res {
		keySplit := strings.Split(k, ":")
		if len(keySplit) != 4 {
			continue
		}

		err = DeleteLoginToken(ctx, userID, keySplit[3])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}
	return nil
}

func DeleteAllLoginToken(ctx context.Context, userID string) errors.WTError {
	res, err := redis.Keys(ctx, fmt.Sprintf("logintoken:%s:*:*", userID)).Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, k := range res {
		keySplit := strings.Split(k, ":")
		if len(keySplit) != 4 {
			continue
		}

		err = DeleteLoginToken(ctx, userID, keySplit[3])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}
	return nil
}

func DeleteAllLoginTokenForAllUser(ctx context.Context) errors.WTError {
	res, err := redis.Keys(ctx, "logintoken:*:*:*").Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, k := range res {
		keySplit := strings.Split(k, ":")
		if len(keySplit) != 4 {
			continue
		}

		err = DeleteLoginToken(ctx, keySplit[1], keySplit[3])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}
	return nil
}

func DeleteAllWebsiteLoginTokenForAllUser(ctx context.Context, webID int64) errors.WTError {
	res, err := redis.Keys(ctx, fmt.Sprintf("logintoken:*:%d:*", webID)).Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, k := range res {
		keySplit := strings.Split(k, ":")
		if len(keySplit) != 4 {
			continue
		}

		err = DeleteLoginToken(ctx, keySplit[1], keySplit[3])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}
	return nil
}

type LoginIPGeo struct {
	UserID    string `json:"userID"`
	Token     string `json:"token"`
	UserToken string `json:"userToken"`
	IP        string `json:"ip"`
	Geo       string `json:"geo"`
	WebID     int64  `json:"webID"`
}

func GetAllLoginTokenGeo(ctx context.Context, uuid string) ([]LoginIPGeo, errors.WTError) {
	res := redis.Keys(ctx, fmt.Sprintf("logintoken:%s:*:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return []LoginIPGeo{}, errors.WarpQuick(err)
	}

	ipgeo := make([]LoginIPGeo, 0, len(tokens))
	for _, t := range tokens {
		keySpilt := strings.Split(t, ":")
		if len(keySpilt) != 4 {
			continue
		}

		loginData, err := ParserLoginToken(ctx, keySpilt[3])
		if err != nil || loginData.UserID != uuid {
			continue
		}

		ipgeo = append(ipgeo, LoginIPGeo{
			UserID:    loginData.UserID,
			Token:     keySpilt[3],
			IP:        loginData.RemoteIP,
			Geo:       loginData.Geo,
			WebID:     loginData.WebID,
			UserToken: loginData.UserToken,
		})
	}

	return ipgeo, nil
}

type LoginTokenData struct {
	WebID     int64  `json:"webID"`
	UserID    string `json:"uuid"`
	UserToken string `json:"userToken"`
	RemoteIP  string `json:"remoteIP"`
	Geo       string `json:"geo"`
}

func CreateLoginToken(ctx context.Context, userID string, webID int64, userToken string) (string, errors.WTError) {
	var err error
	userData, expireTime, err := ParserUserToken(ctx, userToken)
	if err != nil {
		return "", errors.WarpQuick(err)
	} else if userData.UserID != userID || userData.WebsiteID != webID {
		return "", errors.Errorf("bad user token for login token")
	}

	RemoteIP, ok := ctx.Value("X-Real-IP").(string)
	if !ok || len(RemoteIP) == 0 {
		RemoteIP = "0.0.0.0"
	}

	Geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok || len(Geo) == 0 {
		Geo = "未知"
	}

	data := LoginTokenData{
		WebID:     webID,
		UserID:    userID,
		UserToken: userToken,
		RemoteIP:  RemoteIP,
		Geo:       Geo,
	}

	token, _, err := createToken(data, config.BackendConfig.Jwt.Login.Subject, config.BackendConfig.Jwt.Login.Issuer, expireTime, 0)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	key := fmt.Sprintf("logintoken:%s:%d:%s", userID, webID, token)
	err = redis.Set(ctx, fmt.Sprintf("logintoken:%s:%d:%s", userID, webID, token), "1", time.Minute*10).Err()
	if err == nil {
		_ = redis.ExpireAt(ctx, key, expireTime)
	}

	r := record.GetRecordIfExists(ctx)

	recordData := struct {
		RequestsID       string      `json:"requestsID"`
		IP               string      `json:"IP"`
		Geo              string      `json:"geo"`
		Data             interface{} `json:"data"`
		ExpireTime       int64       `json:"expireTime"`
		ExpireTimeFormat string      `json:"expireTimeFormat"`
	}{RequestsID: r.RequestsID, IP: RemoteIP, Geo: Geo, Data: data,
		ExpireTime:       expireTime.Unix(),
		ExpireTimeFormat: expireTime.Format("2006-01-02 15:04:05")}

	recordByte, err := utils.JsonMarshal(recordData)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	recordModel := db.NewTokenRecordModel(mysql.MySQLConn)
	_, mysqlErr := recordModel.Insert(ctx, &db.TokenRecord{
		TokenType: db.LoginToken,
		Token:     token,
		Type:      db.TokenCreate,
		Data:      string(recordByte),
	})
	if mysqlErr != nil {
		return "", errors.WarpQuick(mysqlErr)
	}

	recordDeleteData := struct {
		RequestsID string      `json:"requestsID"`
		Data       interface{} `json:"data"`
	}{RequestsID: r.RequestsID, Data: data}

	recordDeleteByte, err := utils.JsonMarshal(recordDeleteData)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	_, mysqlErr = recordModel.InsertWithCreate(ctx, &db.TokenRecord{
		TokenType: db.LoginToken,
		Token:     token,
		Type:      db.TokenDelete,
		Data:      string(recordDeleteByte),
		CreateAt:  expireTime,
	})
	if mysqlErr != nil {
		return "", errors.WarpQuick(mysqlErr)
	}

	return token, nil
}

func ParserLoginToken(ctx context.Context, tokenString string) (LoginTokenData, errors.WTError) {
	var data LoginTokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.Login.Subject, config.BackendConfig.Jwt.Login.Issuer)
	if err != nil {
		return LoginTokenData{}, err
	}

	userData, _, err := ParserUserToken(ctx, data.UserToken)
	if userData.UserID != data.UserID || userData.WebsiteID != data.WebID {
		return LoginTokenData{}, errors.Errorf("bad user token")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)

	user, mysqlErr := userModel.FindOneByUidWithoutDelete(ctx, data.UserID)
	if errors.Is(err, db.ErrNotFound) {
		return LoginTokenData{}, errors.Errorf("user has beed delete")
	} else if err != nil {
		return LoginTokenData{}, errors.WarpQuick(mysqlErr)
	}

	allow, allowErr := bannedModel.CheckAllow(ctx, user.Id, data.WebID, db.AllowLogin)
	if allowErr != nil {
		return LoginTokenData{}, errors.WarpQuick(allowErr)
	} else if !allow {
		return LoginTokenData{}, errors.Errorf("website not allow by user")
	}

	return data, nil
}

type Login2FATokenData struct {
	UserID string `json:"userID"`
}

func CreateLogin2FAToken(userID string) (string, errors.WTError) {
	token, _, err := createToken(Login2FATokenData{UserID: userID}, config.BackendConfig.Jwt.SecondFA.Subject, config.BackendConfig.Jwt.SecondFA.Issuer, time.Time{}, config.BackendConfig.Jwt.SecondFA.ExpiresSecond)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParserLogin2FAToken(tokenString string) (Login2FATokenData, errors.WTError) {
	var data Login2FATokenData
	_, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.SecondFA.Subject, config.BackendConfig.Jwt.SecondFA.Issuer)
	if err != nil {
		return Login2FATokenData{}, err
	}

	return data, nil
}

func LogoutToken(token string) {
	websocket.TokenConnMapMutex.Lock()
	defer websocket.TokenConnMapMutex.Unlock()

	lst, ok := websocket.TokenConnMap[token]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	msg := websocket.WSMessage{
		Code: userwstype.LogoutToken,
	}

	websocket.WritePeersMessage(wstype.PeersLogoutToken, struct {
		Token string `json:"token"`
	}{Token: token}, msg)

	for _, c := range lst {
		websocket.WriteMessage(c, msg)
	}
}

func LogoutTokenByMsg(token string, msg websocket.WSMessage) {
	websocket.TokenConnMapMutex.Lock()
	defer websocket.TokenConnMapMutex.Unlock()

	lst, ok := websocket.TokenConnMap[token]
	if !ok {
		return
	}

	for _, c := range lst {
		websocket.WriteMessage(c, msg)
	}
}

func DeleteUserToken(ctx context.Context, uuid string, token string) errors.WTError {
	userData, _, err := ParserUserToken(ctx, token)
	if err != nil {
		return nil
	} else if userData.UserID != uuid {
		return errors.Errorf("bad token for user")
	}

	deleteToken(token)

	RemoteIP, ok := ctx.Value("X-Real-IP").(string)
	if !ok || len(RemoteIP) == 0 {
		RemoteIP = "0.0.0.0"
	}

	Geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok || len(Geo) == 0 {
		Geo = "未知"
	}

	r := record.GetRecordIfExists(ctx)

	recordData := struct {
		RequestsID string `json:"requestsID"`
		IP         string `json:"IP"`
		Geo        string `json:"geo"`
	}{RequestsID: r.RequestsID, IP: RemoteIP, Geo: Geo}

	recordByte, err := utils.JsonMarshal(recordData)
	if err != nil {
		return errors.WarpQuick(err)
	}

	recordModel := db.NewTokenRecordModel(mysql.MySQLConn)
	_, mysqlErr := recordModel.Insert(ctx, &db.TokenRecord{
		TokenType: db.UserToken,
		Token:     token,
		Type:      db.TokenDelete,
		Data:      string(recordByte),
	})
	if mysqlErr != nil {
		return errors.WarpQuick(mysqlErr)
	}

	LogoutToken(token)
	return nil
}

func DeleteAllUserWebsiteToken(ctx context.Context, uuid string) errors.WTError {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:website:*:%s:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, t := range tokens {
		keySplit := strings.Split(t, ":")
		if len(keySplit) != 3 {
			continue
		}

		err := DeleteUserToken(ctx, uuid, keySplit[2])
		if err != nil {
			continue
		}
	}

	return nil
}

func DeleteAllUserToken(ctx context.Context, uuid string, exceptToken string) errors.WTError {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:%s:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, t := range tokens {
		keySplit := strings.Split(t, ":")
		if len(keySplit) != 3 {
			continue
		}

		if len(exceptToken) != 0 && keySplit[2] == exceptToken {
			continue
		}

		err := DeleteUserToken(ctx, uuid, keySplit[2])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return nil
}

func DeleteAllFatherUserToken(ctx context.Context, uuid string) errors.WTError {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:father:*:%s:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, t := range tokens {
		keySplit := strings.Split(t, ":")
		if len(keySplit) != 5 {
			continue
		}

		err := DeleteUserToken(ctx, uuid, keySplit[4])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return nil
}

func DeleteAllOneFatherUserToken(ctx context.Context, fatherID string, uuid string) errors.WTError {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:father:%s:%s:*", fatherID, uuid))
	tokens, err := res.Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, t := range tokens {
		keySplit := strings.Split(t, ":")
		if len(keySplit) != 5 {
			continue
		}

		err := DeleteUserToken(ctx, uuid, keySplit[4])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return nil
}

func DeleteAllUserSonToken(ctx context.Context, fatherID string) errors.WTError {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:father:%s:*:*", fatherID))
	tokens, err := res.Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, t := range tokens {
		keySplit := strings.Split(t, ":")
		if len(keySplit) != 5 {
			continue
		}

		err := DeleteUserToken(ctx, keySplit[3], keySplit[4])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return nil
}

func DeleteWebsiteUserToken(ctx context.Context, uuid string, webID int64) errors.WTError {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:website:%d:%s:*", webID, uuid))
	tokens, err := res.Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, t := range tokens {
		keySplit := strings.Split(t, ":")
		if len(keySplit) != 5 {
			continue
		}

		err := DeleteUserToken(ctx, uuid, keySplit[4])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return nil
}

func DeleteAllWebsiteUserToken(ctx context.Context, uuid string) errors.WTError {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:website:*:%s:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, t := range tokens {
		keySplit := strings.Split(t, ":")
		if len(keySplit) != 5 {
			continue
		}

		err := DeleteUserToken(ctx, uuid, keySplit[4])
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return nil
}

const (
	TokenTypeUser    = 1
	TokenTypeFather  = 2
	TokenTypeWebsite = 3
)

type UserIPGeo struct {
	UserID  string `json:"userID"`
	SubType string `json:"subType"`
	Token   string `json:"token"`
	IP      string `json:"ip"`
	Geo     string `json:"geo"`
	NowIP   string `json:"nowIP"`
	NowGeo  string `json:"nowGeo"`
}

type FatherUserIPGeo struct {
	UserID      string `json:"userID"`
	SubType     string `json:"subType"`
	Token       string `json:"token"`
	IP          string `json:"ip"`
	Geo         string `json:"geo"`
	Father      string `json:"father"`
	FatherToken string `json:"fatherToken"`
	NowIP       string `json:"nowIP"`
	NowGeo      string `json:"nowGeo"`
}

type SonUserIPGeo struct {
	UserID  string `json:"userID"`
	SubType string `json:"subType"`
	Token   string `json:"token"`
	IP      string `json:"ip"`
	Geo     string `json:"geo"`
	NowIP   string `json:"nowIP"`
	NowGeo  string `json:"nowGeo"`
}

type WebsiteUserIPGeo struct {
	UserID  string `json:"userID"`
	SubType string `json:"subType"`
	Token   string `json:"token"`
	IP      string `json:"ip"`
	Geo     string `json:"geo"`
	WebID   int64  `json:"webID"`
	NowIP   string `json:"nowIP"`
	NowGeo  string `json:"nowGeo"`
}

func GetUserTokenGeo(ctx context.Context, uuid string, token string) (UserIPGeo, errors.WTError) {
	userData, _, err := ParserUserToken(ctx, token)
	if err != nil {
		return UserIPGeo{}, errors.Errorf("token not found")
	}

	if userData.UserID != uuid {
		return UserIPGeo{}, errors.Errorf("token not found")
	}

	nowIP, nowGeo := GetUserTokenNowGeo(ctx, token)

	return UserIPGeo{
		UserID:  userData.UserID,
		SubType: jwt.SubTypeMap[userData.SubType],
		Token:   token,
		IP:      userData.IP,
		Geo:     userData.Geo,
		NowGeo:  nowGeo,
		NowIP:   nowIP,
	}, nil
}

func GetUserFatherTokenGeo(ctx context.Context, uuid string, token string) (FatherUserIPGeo, errors.WTError) {
	userData, _, err := ParserUserToken(ctx, token)
	if err != nil {
		return FatherUserIPGeo{}, errors.Errorf("token not found")
	}

	if userData.UserID != uuid {
		return FatherUserIPGeo{}, errors.Errorf("token not found")
	}

	nowIP, nowGeo := GetUserTokenNowGeo(ctx, token)

	return FatherUserIPGeo{
		UserID:      userData.UserID,
		SubType:     jwt.SubTypeMap[userData.SubType],
		Token:       token,
		IP:          userData.IP,
		Geo:         userData.Geo,
		Father:      userData.FatherID,
		FatherToken: userData.FatherToken,
		NowGeo:      nowGeo,
		NowIP:       nowIP,
	}, nil
}

func GetAllUserTokenGeo(ctx context.Context, uuid string) ([]UserIPGeo, errors.WTError) {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:%s:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return []UserIPGeo{}, errors.WarpQuick(err)
	}

	ipgeo := make([]UserIPGeo, 0, len(tokens))
	for _, t := range tokens {
		keySpilt := strings.Split(t, ":")
		if len(keySpilt) != 3 {
			continue
		}

		userData, _, err := ParserUserToken(ctx, keySpilt[2])
		if err != nil || userData.UserID != uuid {
			continue
		}

		nowIP, nowGeo := GetUserTokenNowGeo(ctx, keySpilt[2])

		ipgeo = append(ipgeo, UserIPGeo{
			UserID:  userData.UserID,
			SubType: jwt.SubTypeMap[userData.SubType],
			Token:   keySpilt[2],
			IP:      userData.IP,
			Geo:     userData.Geo,
			NowGeo:  nowGeo,
			NowIP:   nowIP,
		})

	}

	return ipgeo, nil
}

func GetAllUserFatherTokenGeo(ctx context.Context, uuid string) ([]FatherUserIPGeo, errors.WTError) {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:father:*:%s:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return []FatherUserIPGeo{}, errors.WarpQuick(err)
	}

	ipgeo := make([]FatherUserIPGeo, 0, len(tokens))
	for _, t := range tokens {
		keySpilt := strings.Split(t, ":")
		if len(keySpilt) != 5 {
			continue
		}

		userData, _, err := ParserUserToken(ctx, keySpilt[4])
		if err != nil || userData.UserID != uuid || userData.FatherID != keySpilt[2] {
			continue
		}

		nowIP, nowGeo := GetUserTokenNowGeo(ctx, keySpilt[4])

		ipgeo = append(ipgeo, FatherUserIPGeo{
			UserID:      userData.UserID,
			SubType:     jwt.SubTypeMap[userData.SubType],
			Token:       keySpilt[4],
			IP:          userData.IP,
			Geo:         userData.Geo,
			Father:      userData.FatherID,
			FatherToken: userData.FatherToken,
			NowIP:       nowIP,
			NowGeo:      nowGeo,
		})

	}

	return ipgeo, nil
}

func GetAllUserSonTokenGeo(ctx context.Context, uuid string) ([]SonUserIPGeo, errors.WTError) {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:father:%s:*:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return []SonUserIPGeo{}, errors.WarpQuick(err)
	}

	ipgeo := make([]SonUserIPGeo, 0, len(tokens))
	for _, t := range tokens {
		keySpilt := strings.Split(t, ":")
		if len(keySpilt) != 5 {
			continue
		}

		userData, _, err := ParserUserToken(ctx, keySpilt[4])
		if err != nil || userData.FatherID != uuid {
			continue
		}

		nowIP, nowGeo := GetUserTokenNowGeo(ctx, keySpilt[4])

		ipgeo = append(ipgeo, SonUserIPGeo{
			UserID:  userData.UserID,
			SubType: jwt.SubTypeMap[userData.SubType],
			Token:   keySpilt[4],
			IP:      userData.IP,
			Geo:     userData.Geo,
			NowIP:   nowIP,
			NowGeo:  nowGeo,
		})

	}

	return ipgeo, nil
}

func GetAllWebsiteTokenGeo(ctx context.Context, uuid string) ([]WebsiteUserIPGeo, errors.WTError) {
	res := redis.Keys(ctx, fmt.Sprintf("usertoken:website:*:%s:*", uuid))
	tokens, err := res.Result()
	if err != nil {
		return []WebsiteUserIPGeo{}, errors.WarpQuick(err)
	}

	ipgeo := make([]WebsiteUserIPGeo, 0, len(tokens))
	for _, t := range tokens {
		keySpilt := strings.Split(t, ":")
		if len(keySpilt) != 5 {
			continue
		}

		keyWebID, err := strconv.ParseInt(keySpilt[2], 10, 64)
		if err != nil {
			continue
		}

		userData, _, err := ParserUserToken(ctx, keySpilt[4])
		if err != nil || userData.UserID != uuid || userData.WebsiteID != keyWebID {
			continue
		}

		nowIP, nowGeo := GetUserTokenNowGeo(ctx, keySpilt[4])

		ipgeo = append(ipgeo, WebsiteUserIPGeo{
			UserID:  userData.UserID,
			SubType: jwt.SubTypeMap[userData.SubType],
			Token:   keySpilt[4],
			IP:      userData.IP,
			Geo:     userData.Geo,
			WebID:   userData.WebsiteID,
			NowGeo:  nowGeo,
			NowIP:   nowIP,
		})
	}

	return ipgeo, nil
}

func IsSubType(t int) bool {
	for i, _ := range jwt.SubTypeMap {
		if i == t {
			return true
		}
	}
	return false
}

type UserTokenData struct {
	SubType     int    `json:"subType"`
	FatherID    string `json:"fatherID"`
	UserID      string `json:"userID"`
	WebsiteID   int64  `json:"websiteID"`
	FatherToken string `json:"fatherToken"`
	Geo         string `json:"geo"`
	IP          string `json:"ip"`
}

func CreateUserToken(ctx context.Context, uuid string, siginin bool, expirationSec int64, tokenType int, fatherToken string, websiteID int64) (string, errors.WTError) {
	if expirationSec == 0 || expirationSec > config.BackendConfig.Jwt.User.ExpiresSecond {
		expirationSec = config.BackendConfig.Jwt.User.ExpiresSecond
	}

	if !IsSubType(tokenType) {
		return "", errors.Errorf("bad token type")
	}

	if tokenType != jwt.UserWebsiteToken {
		websiteID = 0
	}

	var fatherID string
	if tokenType != jwt.UserFatherToken && tokenType != jwt.UserRootFatherToken && tokenType != jwt.UserUncleToken {
		fatherID = ""
		fatherToken = ""
	} else {
		fatherData, _, err := ParserUserToken(ctx, fatherToken)
		if err != nil {
			return "", err
		}

		fatherID = fatherData.UserID
	}

	RemoteIP, ok := ctx.Value("X-Real-IP").(string)
	if !ok || len(RemoteIP) == 0 {
		RemoteIP = "0.0.0.0"
	}

	Geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok || len(Geo) == 0 {
		Geo = "未知"
	}

	data := UserTokenData{
		SubType:     tokenType,
		FatherID:    fatherID,
		UserID:      uuid,
		WebsiteID:   websiteID,
		FatherToken: fatherToken,
		IP:          RemoteIP,
		Geo:         Geo,
	}

	token, clamis, err := createToken(data, config.BackendConfig.Jwt.User.Subject, config.BackendConfig.Jwt.User.Issuer, time.Time{}, expirationSec)
	if err != nil {
		return "", err
	}

	if tokenType == jwt.UserRootToken || tokenType == jwt.UserSonToken || tokenType == jwt.UserHighAuthorityRootToken {
		if siginin {
			err := DeleteAllUserToken(ctx, uuid, "")
			if err != nil {
				logger.Logger.Error("redis error: %s", err.Error())
			}
		}

		_ = redis.Set(ctx, fmt.Sprintf("usertoken:%s:%s", uuid, token), "1", time.Second*time.Duration(expirationSec))
	} else if tokenType == jwt.UserWebsiteToken {
		_ = redis.Set(ctx, fmt.Sprintf("usertoken:website:%d:%s:%s", websiteID, uuid, token), "1", time.Second*time.Duration(expirationSec))
	} else {
		if len(fatherID) == 0 {
			return "", errors.Errorf("bad father id")
		}
		_ = redis.Set(ctx, fmt.Sprintf("usertoken:father:%s:%s:%s", fatherID, uuid, token), "1", time.Second*time.Duration(expirationSec))
	}

	r := record.GetRecordIfExists(ctx)

	recordData := struct {
		RequestsID       string      `json:"requestsID"`
		IP               string      `json:"IP"`
		Geo              string      `json:"geo"`
		Data             interface{} `json:"data"`
		ExpirationSec    int64       `json:"expirationSec"`
		ExpireTime       int64       `json:"expireTime"`
		ExpireTimeFormat string      `json:"expireTimeFormat"`
	}{RequestsID: r.RequestsID, IP: RemoteIP, Geo: Geo, Data: data, ExpirationSec: expirationSec,
		ExpireTime:       clamis.ExpiresAt,
		ExpireTimeFormat: time.Unix(clamis.ExpiresAt, 0).Format("2006-01-02 15:04:05")}

	recordByte, err := utils.JsonMarshal(recordData)
	if err != nil {
		return "", err
	}

	recordModel := db.NewTokenRecordModel(mysql.MySQLConn)
	_, mysqlErr := recordModel.Insert(ctx, &db.TokenRecord{
		TokenType: db.UserToken,
		Token:     token,
		Type:      db.TokenCreate,
		Data:      string(recordByte),
	})
	if err != nil {
		return "", errors.WarpQuick(mysqlErr)
	}

	recordDeleteData := struct {
		RequestsID string      `json:"requestsID"`
		Data       interface{} `json:"data"`
	}{RequestsID: r.RequestsID, Data: data}

	recordDeleteByte, err := utils.JsonMarshal(recordDeleteData)
	if err != nil {
		return "", err
	}

	_, mysqlErr = recordModel.InsertWithCreate(ctx, &db.TokenRecord{
		TokenType: db.UserToken,
		Token:     token,
		Type:      db.TokenDelete,
		Data:      string(recordDeleteByte),
		CreateAt:  time.Unix(clamis.ExpiresAt, 0),
	})
	if err != nil {
		return "", errors.WarpQuick(mysqlErr)
	}

	err = UpdateUserTokenGeo(ctx, token, false)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserTokenNowGeo(ctx context.Context, token string) (IP string, Geo string) {
	key := fmt.Sprintf("usertoken:geo:%s", token)
	old, err := redis.Get(ctx, key).Result()

	if err == nil {
		tmp := strings.Split(old, ";")
		if len(tmp) == 2 {
			return tmp[0], tmp[1]
		}
	}

	return "0.0.0.0", "未知"
}

func UpdateUserTokenGeo(ctx context.Context, token string, updateMysql bool) errors.WTError {
	key := fmt.Sprintf("usertoken:geo:%s", token)
	old, err := redis.Get(ctx, key).Result()

	RemoteIP, ok := ctx.Value("X-Real-IP").(string)
	if !ok || len(RemoteIP) == 0 {
		RemoteIP = "0.0.0.0"
	}

	Geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok || len(Geo) == 0 {
		Geo = "未知"
	}

	var oldIP string
	var oldGeo string

	if err == nil {
		tmp := strings.Split(old, ";")
		if len(tmp) == 2 {
			oldIP = tmp[0]
			oldGeo = tmp[1]
		}
	}

	if oldIP == RemoteIP && oldGeo == Geo {
		return nil
	}

	err = redis.Set(ctx, key, fmt.Sprintf("%s;%s", RemoteIP, Geo), time.Second*time.Duration(config.BackendConfig.Jwt.User.ExpiresSecond)).Err()
	if err != nil {
		return errors.WarpQuick(err)
	}

	if updateMysql {
		r := record.GetRecordIfExists(ctx)

		recordData := struct {
			RequestsID string `json:"requestsID"`
			OldIP      string `json:"oldIP"`
			OldGeo     string `json:"oldGeo"`
			NewIP      string `json:"newIP"`
			NewGeo     string `json:"newGeo"`
		}{RequestsID: r.RequestsID, OldIP: oldIP, OldGeo: oldGeo, NewIP: RemoteIP, NewGeo: Geo}

		recordByte, err := utils.JsonMarshal(recordData)
		if err != nil {
			return errors.WarpQuick(err)
		}

		recordModel := db.NewTokenRecordModel(mysql.MySQLConn)
		_, mysqlErr := recordModel.Insert(ctx, &db.TokenRecord{
			TokenType: db.UserToken,
			Token:     token,
			Type:      db.TokenGeoIPUpdate,
			Data:      string(recordByte),
		})
		if err != nil {
			return errors.WarpQuick(mysqlErr)
		}
	}

	return nil
}

func ParserUserToken(ctx context.Context, tokenString string) (UserTokenData, time.Time, errors.WTError) {
	return parserUserToken(ctx, tokenString, 0)
}

func parserUserToken(ctx context.Context, tokenString string, limit int64) (UserTokenData, time.Time, errors.WTError) {
	if limit > 1000 {
		return UserTokenData{}, time.Time{}, errors.Errorf("bad token limit")
	}

	var data UserTokenData
	clamis, err := parserToken(tokenString, &data, config.BackendConfig.Jwt.User.Subject, config.BackendConfig.Jwt.User.Issuer)
	if err != nil {
		return UserTokenData{}, time.Time{}, err
	}

	if len(data.FatherToken) != 0 {
		fatherData, _, err := parserUserToken(ctx, data.FatherToken, limit+1)
		if err != nil {
			return UserTokenData{}, time.Time{}, err
		}
		if fatherData.UserID != data.FatherID {
			return UserTokenData{}, time.Time{}, errors.Errorf("bad father")
		}
	}

	if !IsSubType(data.SubType) {
		return UserTokenData{}, time.Time{}, errors.Errorf("bad subtype")
	}

	return data, time.Unix(clamis.ExpiresAt, 0), nil
}

type claims struct {
	Issuer       string `json:"iss"`
	Subject      string `json:"sub"`
	ExpiresAt    int64  `json:"exp"`
	NotBefore    int64  `json:"nbf"`
	IssuedNanoAt int64  `json:"iat"`
	Data         string `json:"dat"`
}

func createToken(data interface{}, subject string, issuer string, expireTime time.Time, expireSecond int64) (string, claims, errors.WTError) {
	now := time.Now()

	var exp time.Time
	if expireSecond == 0 {
		exp = expireTime
	} else {
		exp = now.Add(time.Second * time.Duration(expireSecond))
	}

	dataByte, err := utils.JsonMarshal(data)
	if err != nil {
		return "", claims{}, err
	}

	c := claims{
		Issuer:       issuer,
		Subject:      subject,
		Data:         string(dataByte),
		ExpiresAt:    exp.Unix(),
		NotBefore:    now.Unix(),
		IssuedNanoAt: now.UnixNano(),
	}

	d, err := utils.JsonMarshal(c)
	if err != nil {
		return "", claims{}, err
	}

	token := ""
	count := 0

	for {
		if count >= 10 {
			return "", claims{}, errors.Errorf("can not generate n")
		} else {
			count += 1
		}

		n, err := utils.GenerateUniqueNumber(18)
		if err != nil {
			continue
		}

		token = utils.HashSHA256WithBase62(fmt.Sprintf("%s\n%s", string(d), n))

		keyData := fmt.Sprintf("token:data:%s", token)
		res, redisErr := redis.SetNX(context.Background(), keyData, string(d), time.Second*time.Duration(config.BackendConfig.Jwt.ExpireSecond)).Result()
		if redisErr != nil {
			return "", claims{}, errors.WarpQuick(redisErr)
		}

		if !res {
			continue
		}

		break
	}

	return token, c, nil
}

func parserToken(t string, data interface{}, subject string, issuer string) (claims, errors.WTError) {
	key := fmt.Sprintf("token:data:%s", t)
	res, err := redis.Get(context.Background(), key).Result()
	if err != nil {
		return claims{}, errors.Errorf("token not found")
	}

	var d claims
	err = utils.JsonUnmarshal([]byte(res), &d)
	if err != nil {
		return claims{}, errors.WarpQuick(err)
	}

	now := time.Now()

	if d.Subject != subject || d.Issuer != issuer {
		return claims{}, errors.Errorf("token not found")
	}

	if time.Unix(d.NotBefore, 0).After(now) {
		return claims{}, errors.Errorf("token is not active")
	}

	if time.Unix(d.ExpiresAt, 0).Before(now) {
		return claims{}, errors.Errorf("token is expires")
	}

	err = utils.JsonUnmarshal([]byte(d.Data), data)
	if err != nil {
		return claims{}, errors.WarpQuick(err)
	}

	return d, nil
}

func deleteToken(t string) {
	key := fmt.Sprintf("token:data:%s", t)
	_ = redis.Del(context.Background(), key)
}
