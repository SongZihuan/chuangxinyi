package defray

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/auth"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/wuntsong-org/wterrors"
	"strings"
	"time"
)

const MsgSuccess = "MSG_DEFRAY_SUCCESS"
const MsgReturn = "MSG_DEFRAY_RETURN"

type NotifyMsg struct {
	MsgType string `json:"msgType"`
	Data    string `json:"data"`
}

type DefrayReturn struct {
	Status     string `json:"status"`
	DefrayID   string `json:"defrayID"`
	Subject    string `json:"subject"`
	PayerID    string `json:"payerID"`
	Reason     string `json:"reason"`
	DefrayTime int64  `json:"defrayTime"`
	ReturnTime int64  `json:"returnTime"`
}

type DefraySuccess struct {
	Status     string `json:"status"`
	DefrayID   string `json:"defrayID"`
	Subject    string `json:"subject"`
	PayerID    string `json:"payerID"`
	Token      string `json:"token"`
	DefrayTime int64  `json:"defrayTime"`
}

type RedisDefraySuccess struct {
	DefrayID string `json:"defrayID"`
	Token    string `json:"token"`
}

type RedisDefrayReturn struct {
	DefrayID string `json:"defrayID"`
}

var notifyTime = []time.Duration{
	1 * time.Second,
	30 * time.Second,
	30 * time.Second,
	30 * time.Second,
	90 * time.Second,
	90 * time.Second,
	90 * time.Second,
	5 * time.Minute,
	5 * time.Minute,
	5 * time.Minute,
	10 * time.Minute,
	15 * time.Minute,
	20 * time.Minute,
	30 * time.Minute,
	30 * time.Minute,
	30 * time.Minute,
	1 * time.Hour,
	1 * time.Hour,
}

func startAllNotify() errors.WTError {
	keys, err := redis.Keys(context.Background(), "defray:notify:*:*").Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, k := range keys {
		d, err := redis.Get(context.Background(), k).Result()
		dSplit := strings.Split(k, ":")
		if len(dSplit) != 4 {
			continue
		}

		if err != nil {
			continue
		}

		if dSplit[2] == "success" {
			var data RedisDefraySuccess
			err = utils.JsonUnmarshal([]byte(d), &data)
			if err != nil {
				continue
			}

			if dSplit[2] != data.DefrayID {
				continue
			}

			go NotifySuccess(data.DefrayID, data.Token)
		} else if dSplit[2] == "return" {
			var data RedisDefrayReturn
			err = utils.JsonUnmarshal([]byte(d), &data)
			if err != nil {
				continue
			}

			if dSplit[2] != data.DefrayID {
				continue
			}

			go NotifyReturn(data.DefrayID)
		}
	}

	return nil
}

func NotifySuccess(defrayID string, token string) {
	key := fmt.Sprintf("defray:notify:success:%s", defrayID)
	if !redis.AcquireLock(context.Background(), key, 2*time.Hour) {
		return
	}
	defer redis.ReleaseLock(key)

	for i, t := range notifyTime {
		defrayModel := db.NewDefrayModel(mysql.MySQLConn)
		userModel := db.NewUserModel(mysql.MySQLConn)

		defray, err := defrayModel.FindByDefrayID(context.Background(), defrayID)
		if errors.Is(err, db.ErrNotFound) {
			return // 退出循环
		} else if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			return // 退出循环
		}

		if defray.Status != db.DefraySuccess {
			return // 退出循环
		}

		if len(defray.ReturnUrl) == 0 || !utils.IsHttpOrHttps(defray.ReturnUrl) {
			return // 退出循环
		}

		user, err := userModel.FindOneByIDWithoutDelete(context.Background(), defray.UserId.Int64)
		if errors.Is(err, db.ErrNotFound) {
			return // 退出循环
		} else if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			return // 退出循环
		}

		redisDataByte, err := utils.JsonMarshal(RedisDefraySuccess{
			DefrayID: defrayID,
			Token:    token,
		})
		if err != nil {
			return // 退出循环
		}

		notifyDataKey := fmt.Sprintf("defray:notify:success:%s", defrayID)
		if i == 0 {
			_ = redis.Set(context.Background(), notifyDataKey, string(redisDataByte), time.Hour*24*10).Err() // 保留十天
		}

		defrayAt := int64(0)
		if defray.DefrayAt.Valid {
			defrayAt = defray.DefrayAt.Time.Unix()
		}

		dataByte, err := utils.JsonMarshal(DefraySuccess{
			Status:     "SUCCESS",
			DefrayID:   defray.DefrayId,
			Subject:    defray.Subject,
			PayerID:    user.Uid,
			Token:      token,
			DefrayTime: defrayAt,
		})
		if err != nil {
			return // 退出循环
		}

		msg := NotifyMsg{
			MsgType: MsgSuccess,
			Data:    string(dataByte),
		}

		res := func() bool {
			defer utils.Recover(logger.Logger, nil, "")

			resp := struct {
				auth.AuthResp
				Data struct {
					Status string `json:"status"`
				}
			}{}

			_, err := auth.SendRequests(msg, defray.ReturnUrl, config.BackendConfig.User.WebsiteUID, &resp)
			if err != nil {
				logger.Logger.Error("发送支付信息失败：%s", err)
				return false
			}

			if resp.Data.Status != "SUCCESS" {
				return false
			}

			return true
		}()
		if res {
			_ = redis.Del(context.Background(), notifyDataKey)
			return
		}

		if !redis.ExtendLock(context.Background(), key, 2*time.Hour) {
			return
		}
		time.Sleep(t)
	}

	return
}

func NotifyReturn(defrayID string) {
	key := fmt.Sprintf("defray:notify:return:%s", defrayID)
	if !redis.AcquireLock(context.Background(), key, 2*time.Hour) {
		return
	}
	defer redis.ReleaseLock(key)

	for i, t := range notifyTime {
		defrayModel := db.NewDefrayModel(mysql.MySQLConn)
		userModel := db.NewUserModel(mysql.MySQLConn)

		defray, err := defrayModel.FindByDefrayID(context.Background(), defrayID)
		if errors.Is(err, db.ErrNotFound) {
			return // 退出循环
		} else if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			return // 退出循环
		}

		if defray.Status != db.DefrayWaitReturn {
			return // 退出循环
		}

		if len(defray.ReturnUrl) == 0 || !utils.IsHttpOrHttps(defray.ReturnUrl) {
			return // 退出循环
		}

		user, err := userModel.FindOneByIDWithoutDelete(context.Background(), defray.UserId.Int64)
		if errors.Is(err, db.ErrNotFound) {
			return // 退出循环
		} else if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			return // 退出循环
		}

		redisDataByte, err := utils.JsonMarshal(RedisDefrayReturn{
			DefrayID: defrayID,
		})
		if err != nil {
			return // 退出循环
		}

		notifyDataKey := fmt.Sprintf("defray:notify:return:%s", defrayID)
		if i == 0 {
			_ = redis.Set(context.Background(), notifyDataKey, string(redisDataByte), time.Hour*24*10).Err() // 保留十天
		}

		defrayAt := int64(0)
		if defray.DefrayAt.Valid {
			defrayAt = defray.DefrayAt.Time.Unix()
		}

		returnAt := int64(0)
		if defray.ReturnAt.Valid {
			defrayAt = defray.ReturnAt.Time.Unix()
		}

		dataByte, err := utils.JsonMarshal(DefrayReturn{
			Status:     "RETURN",
			DefrayID:   defray.DefrayId,
			Subject:    defray.Subject,
			PayerID:    user.Uid,
			Reason:     defray.ReturnReason.String,
			DefrayTime: defrayAt,
			ReturnTime: returnAt,
		})
		if err != nil {
			return // 退出循环
		}

		msg := NotifyMsg{
			MsgType: MsgReturn,
			Data:    string(dataByte),
		}

		res := func() bool {
			defer utils.Recover(logger.Logger, nil, "")

			resp := struct {
				auth.AuthResp
				Data struct {
					Status string `json:"status"`
				}
			}{}

			_, err := auth.SendRequests(msg, defray.ReturnUrl, config.BackendConfig.User.WebsiteUID, &resp)
			if err != nil {
				logger.Logger.Error("发送退款信息失败：%s", err)
				return false
			}

			if resp.Data.Status != "SUCCESS" {
				return false
			}

			return true
		}()
		if res {
			_ = redis.Del(context.Background(), notifyDataKey)
			return
		}

		if !redis.ExtendLock(context.Background(), key, 2*time.Hour) {
			return
		}
		time.Sleep(t)
	}
}
