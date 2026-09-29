package oss

import (
	"bytes"
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wuntsong-org/wterrors"
)

func UploadFile(fid string, file []byte, fileType string) errors.WTError {
	var err error
	suffix, ok := utils.MediaTypeSuffixMap[fileType]
	if !ok {
		suffix = "unk"
	}

	key := fmt.Sprintf("UI文件/%s.%s", fid, suffix)

	ossModel := db.NewOssFileModel(mysql.MySQLConn)
	_, err = ossModel.InsertWithDelete(context.Background(), &db.OssFile{
		Fid:       fid,
		Key:       key,
		MediaType: fileType,
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = FileBucket.PutObject(key, bytes.NewReader(file), oss.ContentType(fileType))
	if err != nil {
		return errors.WarpQuick(err)
	}
	return nil
}

func GetFile(fid string, download bool) (string, errors.WTError) {
	ossModel := db.NewOssFileModel(mysql.MySQLConn)
	f, err := ossModel.FindByFidWithoutDelete(context.Background(), fid)
	if errors.Is(err, db.ErrNotFound) {
		return "", errors.Errorf("file not found")
	}

	suffix, ok := utils.MediaTypeSuffixMap[f.MediaType]
	if !ok {
		suffix = "unk"
	}

	key := fmt.Sprintf("UI文件/%s.%s", fid, suffix)

	var contentDisposition string
	if download {
		contentDisposition = fmt.Sprintf("attachment;filename=\"%s%s\"", f.Fid, suffix)
	} else {
		contentDisposition = "inline"
	}

	url, err := FileSignBucket.SignURL(key, oss.HTTPGet, 30, oss.ContentDisposition(contentDisposition))
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	return url, nil
}

func DeleteFile(fid string) errors.WTError {
	ossModel := db.NewOssFileModel(mysql.MySQLConn)
	f, err := ossModel.FindByFidWithoutDelete(context.Background(), fid)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("file not found")
	}

	suffix, ok := utils.MediaTypeSuffixMap[f.MediaType]
	if !ok {
		suffix = "unk"
	}

	key := fmt.Sprintf("UI文件/%s.%s", fid, suffix)

	err = FileSignBucket.DeleteObject(key)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
