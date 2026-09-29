package back

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"github.com/google/uuid"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	"github.com/wuntsong-org/wterrors"
	"time"
)

var BadCNY = errors.NewClass("bad cny")

func NewBack(ctx context.Context, get int64, reason string, subject string, user *db.User, canWithdraw bool, supplierID int64) (string, errors.WTError) {
	return NewBackWithSession(ctx, get, reason, subject, user, canWithdraw, supplierID, mysql.MySQLConn)
}

func NewBackWithSession(ctx context.Context, get int64, reason string, subject string, user *db.User, canWithdraw bool, supplierID int64, mysql sqlx.Session) (string, errors.WTError) {
	if get <= 0 {
		return "", BadCNY.New()
	}

	supplier := action.GetWebsite(supplierID)
	if supplier.Status == db.WebsiteStatusBanned {
		return "", errors.Errorf("bad supplier")
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		payModel := db.NewBackModelWithSession(mysql)
		_, err := payModel.FindByBackID(ctx, BackID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := BackID(OutTradeNoUUID.String())

	now := time.Now()

	back := &db.Back{
		WalletId:    user.WalletId,
		UserId:      user.Id,
		BackId:      OutTradeNo,
		Subject:     subject,
		Supplier:    supplier.Name,
		SupplierId:  supplier.ID,
		CanWithdraw: canWithdraw,
		Get:         get,
		CreateAt:    now,
	}

	_, err := balance.BackWithInsert(ctx, user, back, reason, mysql) // 自带Insert
	if err != nil {
		return "", err
	}

	sender.MessageSendBack(back.UserId, back.Get, reason)
	sender.WxrobotSendRecharge(back.UserId, back.Get, reason)
	audit.NewUserAudit(back.UserId, "获得优惠返现（%.2f）", float64(back.Get)/100.0)
	return OutTradeNo, nil
}
