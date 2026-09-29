package oss

import (
	"bytes"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	errors "github.com/wuntsong-org/wterrors"
)

const DefaultHeaderName = "default"

func UploadHeader(header []byte, uid string, filetype string) errors.WTError {
	if !utils.IsUID(uid) {
		return errors.Errorf("bad input")
	}

	filename := fmt.Sprintf("头像/%s", uid)
	err := HeaderBucket.PutObject(filename, bytes.NewReader(header), oss.ContentType(filetype))
	if err != nil {
		return errors.WarpQuick(err)
	}
	return nil
}

func UploadDefaultHeader(header []byte, filetype string) errors.WTError {
	filename := fmt.Sprintf("默认头像/%s", DefaultHeaderName)
	err := HeaderBucket.PutObject(filename, bytes.NewReader(header), oss.ContentType(filetype))
	if err != nil {
		return errors.WarpQuick(err)
	}
	return nil
}

func GetHeader(uid string, needProcess bool) (string, errors.WTError) {
	defaultName := fmt.Sprintf("默认头像/%s", DefaultHeaderName)
	fileName := fmt.Sprintf("头像/%s", uid)

	var url string
	var err error

	if len(uid) == 0 || uid == DefaultHeaderName {
		if needProcess {
			url, err = HeaderSignBucket.SignURL(defaultName, oss.HTTPGet, 30, oss.Process(fmt.Sprintf("style/%s", config.BackendConfig.Aliyun.Header.ImageStyle)))
		} else {
			url, err = HeaderSignBucket.SignURL(defaultName, oss.HTTPGet, 30)
		}
		if err != nil {
			return "", errors.WarpQuick(err)
		}
	} else if utils.IsUID(uid) {
		var isExists bool
		isExists, err = HeaderSignBucket.IsObjectExist(fileName)
		if err != nil {
			return "", errors.WarpQuick(err)
		}

		if !isExists {
			if needProcess {
				url, err = HeaderSignBucket.SignURL(defaultName, oss.HTTPGet, 30, oss.Process(fmt.Sprintf("style/%s", config.BackendConfig.Aliyun.Header.ImageStyle)))
			} else {
				url, err = HeaderSignBucket.SignURL(defaultName, oss.HTTPGet, 30)
			}
		} else {
			if needProcess {
				url, err = HeaderSignBucket.SignURL(fileName, oss.HTTPGet, 30, oss.Process(fmt.Sprintf("style/%s", config.BackendConfig.Aliyun.Header.ImageStyle)))
			} else {
				url, err = HeaderSignBucket.SignURL(fileName, oss.HTTPGet, 30)
			}
		}
		if err != nil {
			return "", errors.WarpQuick(err)
		}
	} else {
		return "", errors.Errorf("bad uid")
	}

	return url, nil
}
