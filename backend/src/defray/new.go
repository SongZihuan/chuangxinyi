package defray

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/jwt"
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

var Insufficient = errors.NewClass("insufficient") // 余额不足

func NewDefray(ctx context.Context, data jwt.DefrayTokenData, owner *db.User) (string, string, errors.WTError) {
	if len(data.ReturnURL) == 0 {
		return "", "", errors.Errorf("bad return url")
	}

	if owner == nil {
		owner = &db.User{
			Id: 0,
		}
	}

	if owner.Id != data.OwnerID {
		return "", "", errors.Errorf("bad owner id")
	}

	if data.OwnerID == 0 && data.MustSelfDefray {
		return "", "", errors.Errorf("bad payer")
	}

	if data.Price < 0 {
		return "", "", errors.Errorf("bad price")
	}

	if data.UnitPrice < 0 {
		return "", "", errors.Errorf("bad unit price")
	}

	if data.Quantity <= 0 {
		return "", "", errors.Errorf("bad quantity")
	}

	if data.UnitPrice*data.Quantity < data.Price {
		return "", "", errors.Errorf("bad price or unit price")
	}

	supplier := action.GetWebsite(data.SupplierID)
	if supplier.Status == db.WebsiteStatusBanned {
		return "", "", errors.Errorf("bad supplier")
	}

	DefrayIDUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		defrayModel := db.NewDefrayModel(mysql.MySQLConn)
		_, err := defrayModel.FindByDefrayID(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", "", errors.Errorf("generate outtradeno fail")
	}

	DefrayID := DefrayIDUUID.String()
	data.TradeID = DefrayID

	now := time.Now()

	defrayModel := db.NewDefrayModel(mysql.MySQLConn)
	if data.Price == 0 && data.MustSelfDefray {
		d := &db.Defray{
			OwnerId: sql.NullInt64{
				Valid: data.OwnerID != 0,
				Int64: data.OwnerID,
			},
			DefrayId:           DefrayID,
			Subject:            data.Subject,
			Price:              data.Price,
			UnitPrice:          data.UnitPrice,
			Quantity:           data.Quantity,
			Describe:           data.Describe,
			SupplierId:         supplier.ID,
			Supplier:           supplier.Name,
			ReturnUrl:          data.ReturnURL,
			InvitePre:          data.InvitePre,
			DistributionLevel1: data.DistributionLevel1,
			DistributionLevel2: data.DistributionLevel2,
			DistributionLevel3: data.DistributionLevel3,
			CanWithdraw:        data.CanWithdraw,
			MustSelfDefray:     data.MustSelfDefray,
			ReturnDayLimit:     data.ReturnDayLimit,
			Status:             db.DefrayWait, // balance需要DefrayWait
			UserId: sql.NullInt64{
				Valid: true,
				Int64: owner.Id,
			},
			WalletId: sql.NullInt64{
				Valid: true,
				Int64: owner.WalletId,
			},
			DefrayAt: sql.NullTime{
				Valid: true,
				Time:  now,
			},
			LastReturnAt: sql.NullTime{
				Valid: true,
				Time:  now.Add(time.Duration(data.ReturnDayLimit) * time.Hour * 24),
			},
			RealPrice: sql.NullInt64{
				Valid: true,
				Int64: 0,
			},
		}

		err := mysql.MySQLConn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
			_, _, err := balance.Defray(ctx, owner, d, session)
			if errors.Is(err, balance.Insufficient) {
				return Insufficient.New()
			} else if err != nil {
				return err
			}

			err = waitOrDistribution(d, session)
			if err != nil {
				return err
			}

			defrayModel := db.NewDefrayModelWithSession(session)
			_, mysqlErr := defrayModel.Insert(ctx, d)
			if err != nil {
				return mysqlErr
			}

			return nil
		})
		if err != nil {
			return "", "", errors.WarpQuick(err)
		}

		return "", DefrayID, nil
	} else {
		_, err := defrayModel.Insert(ctx, &db.Defray{
			OwnerId: sql.NullInt64{
				Valid: data.OwnerID != 0,
				Int64: data.OwnerID,
			},
			DefrayId:           DefrayID,
			Subject:            data.Subject,
			Price:              data.Price,
			UnitPrice:          data.UnitPrice,
			Quantity:           data.Quantity,
			Describe:           data.Describe,
			SupplierId:         supplier.ID,
			Supplier:           supplier.Name,
			ReturnUrl:          data.ReturnURL,
			InvitePre:          data.InvitePre,
			DistributionLevel1: data.DistributionLevel1,
			DistributionLevel2: data.DistributionLevel2,
			DistributionLevel3: data.DistributionLevel3,
			CanWithdraw:        data.CanWithdraw,
			MustSelfDefray:     data.MustSelfDefray,
			ReturnDayLimit:     data.ReturnDayLimit,
			Status:             db.DefrayWait,
		})
		if err != nil {
			return "", "", errors.WarpQuick(err)
		}

		token, err := jwt.CreateDefrayTokenToken(data)
		if err != nil {
			return "", "", errors.WarpQuick(err)
		}

		return token, DefrayID, nil
	}
}

func ParserDefray(defrayToken string) (dd jwt.DefrayTokenData, t time.Time, resErr errors.WTError) {
	data, err := jwt.ParserDefrayToken(defrayToken)
	if err != nil {
		return jwt.DefrayTokenData{}, time.Time{}, err
	}

	return data, time.Unix(data.TimeExpire, 0), nil
}

type AdminDefrayData struct {
	OwnerID            int64  `json:"ownerID"`
	Subject            string `json:"subject"`  // 标题
	Price              int64  `json:"price"`    // 价格
	Describe           string `json:"describe"` // 描述
	InvitePre          int64  `json:"invitePre"`
	DistributionLevel1 int64  `json:"distributionLevel1"`
	DistributionLevel2 int64  `json:"distributionLevel2"`
	DistributionLevel3 int64  `json:"distributionLevel3"`
	CanWithdraw        bool   `json:"canWithdraw"`
	SupplierID         int64  `json:"supplier"`
	ReturnDayLimit     int64  `json:"returnDayLimit"`
}

func NewAdminDefray(ctx context.Context, data AdminDefrayData, user *db.User) errors.WTError {
	if data.Price < 0 {
		return errors.Errorf("bad price")
	}

	DefrayIDUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		defrayModel := db.NewDefrayModel(mysql.MySQLConn)
		_, err := defrayModel.FindByDefrayID(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return errors.Errorf("generate outtradeno fail")
	}

	supplier := action.GetWebsite(data.SupplierID)
	if supplier.Status == db.WebsiteStatusBanned {
		return errors.Errorf("unknown supplier")
	}

	now := time.Now()

	d := &db.Defray{
		OwnerId: sql.NullInt64{
			Valid: data.OwnerID != 0,
			Int64: data.OwnerID,
		},
		DefrayId:           DefrayIDUUID.String(),
		Subject:            data.Subject,
		Price:              data.Price,
		UnitPrice:          data.Price,
		Quantity:           1,
		Describe:           data.Describe,
		SupplierId:         supplier.ID,
		Supplier:           supplier.Name,
		ReturnUrl:          "",
		InvitePre:          data.InvitePre,
		DistributionLevel1: data.DistributionLevel1,
		DistributionLevel2: data.DistributionLevel2,
		DistributionLevel3: data.DistributionLevel3,
		CanWithdraw:        data.CanWithdraw,
		ReturnDayLimit:     data.ReturnDayLimit,
		Status:             db.DefrayWait, // balance需要DefrayWait
		UserId: sql.NullInt64{
			Valid: true,
			Int64: user.Id,
		},
		WalletId: sql.NullInt64{
			Valid: true,
			Int64: user.WalletId,
		},
		DefrayAt: sql.NullTime{
			Valid: true,
			Time:  now,
		},
		LastReturnAt: sql.NullTime{
			Valid: true,
			Time:  now.Add(time.Duration(data.ReturnDayLimit) * time.Hour * 24),
		},
		RealPrice: sql.NullInt64{
			Valid: true,
			Int64: data.Price,
		},
		MustSelfDefray: false,
	}

	var realPrice int64
	if d.Price > 0 {
		if d.InvitePre > 0 && user.InviteId.Valid { // 有邀请人要打折
			realPrice = int64((float64(d.InvitePre) / 100) * float64(realPrice))
			if realPrice < 0 {
				realPrice = 0
			}
		}
	} else if d.Price == 0 {
		realPrice = 0
	} else {
		return errors.Errorf("bad price")
	}

	err := mysql.MySQLConn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		defrayModel := db.NewDefrayModelWithSession(mysql.MySQLConn)

		_, _, err := balance.Defray(ctx, user, d, session)
		if errors.Is(err, balance.Insufficient) {
			return Insufficient.New()
		} else if err != nil {
			return errors.WarpQuick(err)
		}

		err = waitOrDistribution(d, session)
		if err != nil {
			return err
		}

		_, mysqlErr := defrayModel.Insert(ctx, d)
		if mysqlErr != nil {
			return errors.WarpQuick(err)
		}

		if d.Price > 0 {
			err = waitOrDistribution(d, session)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	sender.PhoneSendChange(d.UserId.Int64, "余额（管理员录入订单消费）")
	sender.EmailSendChange(d.UserId.Int64, "余额（管理员录入订单消费）")
	sender.MessageSendPayAdmin(d.UserId.Int64, d.Price, d.Subject)
	sender.WxrobotSendPayAdmin(d.UserId.Int64, d.Price, d.Subject)
	sender.FuwuhaoSendDefray(d)
	audit.NewUserAudit(d.UserId.Int64, "管理员录入订单消费成功（%.2f）", float64(d.Price)/100.00)

	return nil
}
