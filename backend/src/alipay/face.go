package alipay

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/SuperH-0630/gopay"
	"github.com/google/uuid"
	"github.com/wuntsong-org/wterrors"
	"time"
)

func NewFaceCheck(ctx context.Context, name string, id string) (string, string, errors.WTError) {
	if !config.BackendConfig.Alipay.UseFaceCheck {
		return "", "", errors.Errorf("not ok")
	}

	if !utils.IsValidChineseName(name) || !utils.IsValidIDCard(id) {
		return "", "", errors.Errorf("bad name or idcard")
	}

	outerOrderNoUUID, success := redis.GenerateUUIDMore(ctx, "face", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		faceCheckModel := db.NewFaceCheckModel(mysql.MySQLConn)
		_, err := faceCheckModel.FindOneByCheckID(ctx, AlipayFaceID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", "", errors.Errorf("generate outtradeno fail")
	}

	outerOrderNo := AlipayFaceID(outerOrderNoUUID.String())

	bm := make(gopay.BodyMap)
	bm.Set("outer_order_no", outerOrderNo)
	bm.Set("biz_code", "FACE")

	ip := make(gopay.BodyMap)
	ip.Set("identity_type", "CERT_INFO")
	ip.Set("cert_type", "IDENTITY_CARD")
	ip.Set("cert_type", "IDENTITY_CARD")
	ip.Set("cert_name", name)
	ip.Set("cert_no", id)
	bm.Set("identity_param", ip)

	mc := make(gopay.BodyMap)
	mc.Set("return_url", config.BackendConfig.Alipay.FaceReturnUrl)
	bm.Set("merchant_config", mc)

	res, err := AlipayClient.UserCertifyOpenInit(ctx, bm)
	if err != nil {
		return "", "", errors.WarpQuick(err)
	}

	certifyID := res.Response.CertifyId

	bm2 := make(gopay.BodyMap)
	bm2.Set("certify_id", certifyID)

	url, err := AlipayClient.UserCertifyOpenCertify(ctx, bm2)
	if err != nil {
		return "", "", errors.WarpQuick(err)
	}

	redis.SetCache(ctx, fmt.Sprintf("faceurl:%s", certifyID), url, time.Minute*15)

	faceCheckModel := db.NewFaceCheckModel(mysql.MySQLConn)
	_, err = faceCheckModel.Insert(ctx, &db.FaceCheck{
		CheckId:   outerOrderNo,
		CertifyId: certifyID,
		Name:      name,
		Idcard:    id,
		Status:    db.FaceCheckWait,
	})
	if err != nil {
		return "", "", errors.WarpQuick(err)
	}

	return certifyID, url, nil
}

func GetFaceUrl(ctx context.Context, certifyID string) (string, errors.WTError) {
	key := fmt.Sprintf("faceurl:%s", certifyID)
	res, ok := redis.GetCache(ctx, key)
	if !ok {
		return "", errors.Errorf("not found")
	}

	return res, nil
}

func QueryFaceCheck(ctx context.Context, certifyID string) (string, string, int64, string, errors.WTError) {
	faceCheckModel := db.NewFaceCheckModel(mysql.MySQLConn)
	face, err := faceCheckModel.FindOneByCertifyID(ctx, certifyID)
	if errors.Is(err, db.ErrNotFound) {
		return "", "", 0, "", errors.Errorf("face info not found")
	} else if err != nil {
		return "", "", 0, "", errors.WarpQuick(err)
	}

	if face.Status != db.FaceCheckWait {
		return face.Name, face.Idcard, face.Status, face.CheckId, nil
	}

	bm := make(gopay.BodyMap)
	bm.Set("certify_id", certifyID)
	res, err := AlipayClient.UserCertifyOpenQuery(ctx, bm)
	if err != nil {
		return "", "", 0, "", errors.WarpQuick(err)
	}

	if res.Response.Passed == "T" {
		face.Status = db.FaceCheckOK
	} else if res.Response.Passed == "F" {
		face.Status = db.FaceCheckFail
	}

	err = faceCheckModel.Update(ctx, face)
	if err != nil {
		return "", "", 0, "", errors.WarpQuick(err)
	}

	key := fmt.Sprintf("faceurl:%s", certifyID)
	redis.DelCache(ctx, key)

	return face.Name, face.Idcard, face.Status, face.CheckId, nil
}
