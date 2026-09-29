package oss

import (
	"bytes"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

func UploadLicense(license []byte, filetypeLicense string, idcard []byte, filetypeIDCard string, idcardback []byte, filetypeIDCardBack string, name string, id string, legalPerson string, legalPersonID string) (string, string, string, errors.WTError) {
	if !utils.IsValidChineseCompanyName(name) || !utils.IsValidCreditCode(id) || !utils.IsValidChineseName(legalPerson) || !utils.IsValidIDCard(legalPersonID) {
		return "", "", "", errors.Errorf("bad input data")
	}

	key := fmt.Sprintf("营业执照和法人身份证/%s/%s/%s/%s/%d", name, id, legalPerson, legalPersonID, time.Now().Unix())

	licenseFileName := fmt.Sprintf("%s/license.%s", key, utils.MediaTypeSuffixMap[filetypeLicense])
	err := IdentityBucket.PutObject(licenseFileName, bytes.NewReader(license))
	if err != nil {
		return "", "", "", errors.WarpQuick(err)
	}

	idcardFileName := fmt.Sprintf("%s/idcard.%s", key, utils.MediaTypeSuffixMap[filetypeIDCard])
	err = IdentityBucket.PutObject(idcardFileName, bytes.NewReader(idcard))
	if err != nil {
		return "", "", "", errors.WarpQuick(err)
	}

	idcardbackFileName := fmt.Sprintf("%s/idcardback.%s", key, utils.MediaTypeSuffixMap[filetypeIDCardBack])
	err = IdentityBucket.PutObject(idcardbackFileName, bytes.NewReader(idcardback))
	if err != nil {
		return "", "", "", errors.WarpQuick(err)
	}

	return licenseFileName, idcardFileName, idcardbackFileName, nil
}

func UploadIDCard(idcard []byte, filetypeIDCard string, idcardback []byte, filetypeIDCardBack string, name string, idnumber string) (string, string, errors.WTError) {
	if !utils.IsValidChineseName(name) || !utils.IsValidIDCard(idnumber) {
		return "", "", errors.Errorf("bad input")
	}

	key := fmt.Sprintf("使用人身份证/%s/%s/%d", name, idnumber, time.Now().Unix())

	idcardFileName := fmt.Sprintf("%s/idcard.%s", key, utils.MediaTypeSuffixMap[filetypeIDCard])
	err := IdentityBucket.PutObject(idcardFileName, bytes.NewReader(idcard))
	if err != nil {
		return "", "", errors.WarpQuick(err)
	}

	idcardbackFileName := fmt.Sprintf("%s/idcardback.%s", key, utils.MediaTypeSuffixMap[filetypeIDCardBack])
	err = IdentityBucket.PutObject(idcardbackFileName, bytes.NewReader(idcardback))
	if err != nil {
		return "", "", errors.WarpQuick(err)
	}

	return idcardFileName, idcardbackFileName, nil
}

func GetIdentity(key string) (string, errors.WTError) {
	url, err := IdentitySignBucket.SignURL(key, oss.HTTPGet, 30)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	return url, nil
}
