package oss

import (
	"bytes"
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

func UploadWorkOrderFile(order *db.WorkOrder, communicate *db.WorkOrderCommunicate, fid string, file []byte, fast bool) errors.WTError {
	keyUUID, success := redis.GenerateUUIDMore(context.Background(), "order:file", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		orderFileModel := db.NewWorkOrderCommunicateFileModel(mysql.MySQLConn)
		_, err := orderFileModel.FindByKeyWithoutDelete(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}
		return false
	})
	if !success {
		return errors.Errorf("generate key fail")
	}

	key := fmt.Sprintf("工单文件/%s/%d/%s/%s", order.Uid, communicate.Id, keyUUID.String(), fid) // fig作为文件名

	orderFileModel := db.NewWorkOrderCommunicateFileModel(mysql.MySQLConn)
	_, err := orderFileModel.Insert(context.Background(), &db.WorkOrderCommunicateFile{
		OrderId:       order.Id,
		CommunicateId: communicate.Id,
		Fid:           fid,
		Key:           keyUUID.String(),
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	if fast {
		go func() {
			err = WorkOrderFileBucket.PutObject(key, bytes.NewReader(file))
			if err != nil {
				logger.Logger.Error("upload order file oss error: %s", err.Error())
			}
		}()
	} else {
		err = WorkOrderFileBucket.PutObject(key, bytes.NewReader(file))
		if err != nil {
			return errors.WarpQuick(err)
		}
	}
	return nil
}

func GetWorkOrderFile(order *db.WorkOrder, communicate *db.WorkOrderCommunicate, fid string, download bool) (string, errors.WTError) {
	ossModel := db.NewWorkOrderCommunicateFileModel(mysql.MySQLConn)
	f, mysqlErr := ossModel.FindByFidWithoutDelete(context.Background(), fid)
	if errors.Is(mysqlErr, db.ErrNotFound) {
		return "", errors.Errorf("file not found")
	}

	var contentDisposition string
	key := fmt.Sprintf("工单文件/%s/%d/%s/%s", order.Uid, communicate.Id, f.Key, fid)
	if download {
		contentDisposition = fmt.Sprintf("attachment;filename=\"%s\"", fid)
	} else {
		contentDisposition = "inline"
	}
	url, err := WorkOrderFileSignBucket.SignURL(key, oss.HTTPGet, 30, oss.ContentDisposition(contentDisposition))
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	return url, nil
}
