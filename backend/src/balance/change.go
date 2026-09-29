package balance

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/peers/wstype"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/center/userwstype"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	"github.com/wuntsong-org/wterrors"
	"time"
)

var Insufficient = errors.NewClass("insufficient") // 余额不足
var InsufficientQuota = errors.NewClass("insufficient quota")

type UserBalance struct {
	Balance      int64 `json:"balance"`
	WaitBalance  int64 `json:"waitBalance"`
	Cny          int64 `json:"cny"`
	NotBilled    int64 `json:"notBilled"` // 可能为负数，表示倒欠的发票
	Billed       int64 `json:"billed"`
	HasBilled    int64 `json:"hasBilled"`
	WalletID     int64 `json:"walletID"`
	Withdraw     int64 `json:"withdraw"`
	WaitWithdraw int64 `json:"wait_withdraw"`
	NotWithdraw  int64 `json:"not_withdraw"`
	HasWithdraw  int64 `json:"has_withdraw"`
}

func NewWalletRecord(wallet *db.Wallet, user *db.User, t int64, id string, reason string) db.WalletRecord {
	return db.WalletRecord{
		WalletId:  wallet.Id,
		UserId:    user.Id,
		Type:      t,
		FundingId: id,
		Reason:    reason,

		Balance:      wallet.Balance,
		Cny:          wallet.Cny,
		WaitBalance:  wallet.WaitBalance,
		NotBilled:    wallet.NotBilled,
		Billed:       wallet.Billed,
		HasBilled:    wallet.HasBilled,
		Withdraw:     wallet.Withdraw,
		WaitWithdraw: wallet.WaitWithdraw,
		NotWithdraw:  wallet.NotWithdraw,
		HasWithdraw:  wallet.HasWithdraw,

		BeforeBalance:      wallet.Balance,
		BeforeCny:          wallet.Cny,
		BeforeWaitBalance:  wallet.WaitBalance,
		BeforeNotBilled:    wallet.NotBilled,
		BeforeBilled:       wallet.Billed,
		BeforeHasBilled:    wallet.HasBilled,
		BeforeWithdraw:     wallet.Withdraw,
		BeforeWaitWithdraw: wallet.WaitWithdraw,
		BeforeNotWithdraw:  wallet.NotWithdraw,
		BeforeHasWithdraw:  wallet.HasWithdraw,
	}
}

func UpdateWalletByMsg(walletID int64, msg websocket.WSMessage) {
	websocket.WalletConnMapMutex.Lock()
	defer websocket.WalletConnMapMutex.Unlock()

	lst, ok := websocket.WalletConnMap[walletID]
	if !ok {
		return
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func UpdateWallet(wallet *db.Wallet, ch chan websocket.WSMessage) {
	var lst []chan websocket.WSMessage

	websocket.WalletConnMapMutex.Lock()
	defer websocket.WalletConnMapMutex.Unlock()

	if ch == nil {
		var ok bool
		lst, ok = websocket.WalletConnMap[wallet.Id]
		if !ok {
			lst = []chan websocket.WSMessage{}
		}

	} else {
		lst = []chan websocket.WSMessage{ch}
	}

	msg := websocket.WSMessage{
		Code: userwstype.UpdateWalletInfo,
		Data: UserBalance{
			Balance:      wallet.Balance,
			WaitBalance:  wallet.WaitBalance,
			Cny:          wallet.Cny,
			NotBilled:    wallet.NotBilled,
			Billed:       wallet.Billed,
			HasBilled:    wallet.HasBilled,
			Withdraw:     wallet.Withdraw,
			WaitWithdraw: wallet.WaitWithdraw,
			NotWithdraw:  wallet.NotWithdraw,
			HasWithdraw:  wallet.HasWithdraw,
			WalletID:     wallet.Id,
		},
	}

	if ch != nil {
		websocket.WritePeersMessage(wstype.PeersUpdateWalletInfo, struct {
			WalletID int64 `json:"walletID"`
		}{WalletID: wallet.Id}, msg)
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func Pay(ctx context.Context, user *db.User, pay *db.Pay) (int64, errors.WTError) {
	if pay.TradeStatus != db.PayWait {
		return 0, errors.Errorf("double pay")
	}

	walletModel := db.NewWalletModel(mysql.MySQLConn)

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if pay.WalletId != wallet.Id {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletPay, pay.PayId, "用户充值")

	websiteFunding := &db.WebsiteFunding{
		WebId:       warp.UserCenterWebsite,
		Type:        db.WebsiteFundingPay,
		FundingId:   pay.PayId,
		Profit:      pay.Cny,
		Expenditure: 0,
		Year:        int64(pay.PayAt.Time.Year()),
		Month:       int64(pay.PayAt.Time.Month()),
		Day:         int64(pay.PayAt.Time.Day()),
		PayAt:       pay.PayAt.Time,
	}

	wallet.Balance += pay.Get
	wallet.Cny += pay.Cny

	walletRecord.Balance = wallet.Balance
	walletRecord.Cny = wallet.Cny

	pay.TradeStatus = db.PaySuccess
	pay.Balance = sql.NullInt64{
		Valid: true,
		Int64: wallet.Balance,
	}

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)
		payModel := db.NewPayModelWithSession(session)
		websiteFundingModel := db.NewWebsiteFundingModelWithSession(session)

		_, err := websiteFundingModel.Insert(ctx, websiteFunding)
		if err != nil {
			return errors.WarpQuick(err)
		}

		_, err = walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = payModel.Update(ctx, pay)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, nil
}

func PayRefund(ctx context.Context, user *db.User, pay *db.Pay) (int64, errors.WTError) {
	if pay.TradeStatus != db.PaySuccess && pay.TradeStatus != db.PayWait && pay.TradeStatus != db.PayCloseRefund {
		return 0, errors.Errorf("double refund")
	}

	walletModel := db.NewWalletModel(mysql.MySQLConn)
	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if pay.WalletId != wallet.Id {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletPay, pay.PayId, "用户充值退款")

	pay.TradeStatus = db.PayWaitRefund
	pay.RefundAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	websiteFunding := &db.WebsiteFunding{
		WebId:       warp.UserCenterWebsite,
		Type:        db.WebsiteFundingPayRefund,
		FundingId:   pay.PayId,
		Profit:      0,
		Expenditure: pay.Cny,
		Year:        int64(pay.RefundAt.Time.Year()),
		Month:       int64(pay.RefundAt.Time.Month()),
		Day:         int64(pay.RefundAt.Time.Day()),
		PayAt:       pay.RefundAt.Time,
	}

	wallet.Balance -= pay.Get
	wallet.Cny -= pay.Cny

	walletRecord.Balance = wallet.Balance
	walletRecord.Cny = wallet.Cny

	if wallet.Balance < 0 || wallet.Cny < 0 {
		return 0, Insufficient.New()
	}

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)
		payModel := db.NewPayModelWithSession(session)
		websiteFundingModel := db.NewWebsiteFundingModelWithSession(session)

		_, err = websiteFundingModel.Insert(ctx, websiteFunding)
		if err != nil {
			return errors.WarpQuick(err)
		}

		_, err := walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = payModel.Update(ctx, pay)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, nil
}

func PayRefundFail(ctx context.Context, user *db.User, pay *db.Pay, status int64) (int64, errors.WTError) {
	if pay.TradeStatus != db.PayWaitRefund && pay.TradeStatus != db.PaySuccessRefund && pay.TradeStatus != db.PaySuccessRefundInside && pay.TradeStatus != db.PayCloseRefund {
		return 0, errors.Errorf("double refund")
	}

	walletModel := db.NewWalletModel(mysql.MySQLConn)
	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if pay.WalletId != wallet.Id {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletPay, pay.PayId, "用户充值退款失败")

	now := time.Now()

	websiteFunding := &db.WebsiteFunding{
		WebId:       warp.UserCenterWebsite,
		Type:        db.WebsiteFundingPayRefundFail,
		FundingId:   pay.PayId,
		Profit:      pay.Cny,
		Expenditure: 0,
		Year:        int64(now.Year()),
		Month:       int64(now.Month()),
		Day:         int64(now.Day()),
		PayAt:       now,
	}

	wallet.Balance += pay.Get
	wallet.Cny += pay.Cny

	walletRecord.Balance = wallet.Balance
	walletRecord.Cny = wallet.Cny

	pay.TradeStatus = status

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)
		payModel := db.NewPayModelWithSession(session)
		websiteFundingModel := db.NewWebsiteFundingModelWithSession(session)

		_, err = websiteFundingModel.Insert(ctx, websiteFunding)
		if err != nil {
			return errors.WarpQuick(err)
		}

		_, err := walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = payModel.Update(ctx, pay)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	err = walletModel.Update(ctx, wallet)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, nil
}

func Defray(ctx context.Context, user *db.User, defray *db.Defray, session sqlx.Session) (int64, int64, errors.WTError) {
	if defray.Status != db.DefrayWait {
		return 0, 0, errors.Errorf("double defray")
	}

	if !defray.UserId.Valid {
		return 0, 0, errors.Errorf("not user id")
	}

	walletModel := db.NewWalletModelWithSession(session)
	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, 0, errors.WarpQuick(err)
	} else if wallet.Id != defray.WalletId.Int64 {
		return 0, 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletDefray, defray.DefrayId, "用户消费")

	websiteFunding := &db.WebsiteFunding{
		WebId:       defray.SupplierId,
		Type:        db.WebsiteFundingDefray,
		FundingId:   defray.DefrayId,
		Profit:      defray.RealPrice.Int64,
		Expenditure: 0,
		Year:        int64(defray.DefrayAt.Time.Year()),
		Month:       int64(defray.DefrayAt.Time.Month()),
		Day:         int64(defray.DefrayAt.Time.Day()),
		PayAt:       defray.DefrayAt.Time,
	}

	wallet.Balance -= defray.RealPrice.Int64
	wallet.NotBilled += defray.RealPrice.Int64
	wallet.Billed += defray.RealPrice.Int64

	walletRecord.Balance = wallet.Balance
	walletRecord.NotBilled = wallet.NotBilled
	walletRecord.Billed = wallet.Billed

	if wallet.Balance < 0 {
		return 0, -wallet.Balance, Insufficient.New()
	} // 返回资金不足的额度

	defray.Status = db.DefraySuccess
	defray.Balance = sql.NullInt64{
		Valid: true,
		Int64: wallet.Balance,
	}

	walletRecordModel := db.NewWalletRecordModelWithSession(session)
	defrayModel := db.NewDefrayModelWithSession(session)
	websiteFundingModel := db.NewWebsiteFundingModelWithSession(session)

	_, err = websiteFundingModel.Insert(ctx, websiteFunding)
	if err != nil {
		return 0, 0, errors.WarpQuick(err)
	}

	_, err = walletRecordModel.Insert(ctx, &walletRecord)
	if err != nil {
		return 0, 0, errors.WarpQuick(err)
	}

	err = defrayModel.Update(ctx, defray)
	if err != nil {
		return 0, 0, errors.WarpQuick(err)
	}

	err = walletModel.Update(ctx, wallet)
	if err != nil {
		return 0, 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, 0, nil
}

func DefrayReturn(ctx context.Context, user *db.User, defray *db.Defray, must bool) (int64, errors.WTError) {
	if defray.Status != db.DefrayWaitReturn {
		return 0, errors.Errorf("double return defray")
	}

	if !defray.UserId.Valid {
		return 0, errors.Errorf("not user id")
	}

	walletModel := db.NewWalletModel(mysql.MySQLConn)
	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if wallet.Id != defray.WalletId.Int64 {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletDefray, defray.DefrayId, "用户消费退款")

	defray.Status = db.DefrayReturn
	defray.ReturnAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	websiteFunding := &db.WebsiteFunding{
		WebId:       defray.SupplierId,
		Type:        db.WebsiteFundingDefrayReturn,
		FundingId:   defray.DefrayId,
		Profit:      defray.RealPrice.Int64,
		Expenditure: 0,
		Year:        int64(defray.ReturnAt.Time.Year()),
		Month:       int64(defray.ReturnAt.Time.Month()),
		Day:         int64(defray.ReturnAt.Time.Day()),
		PayAt:       defray.ReturnAt.Time,
	}

	wallet.Balance += defray.RealPrice.Int64
	wallet.NotBilled -= defray.RealPrice.Int64

	walletRecord.Balance = wallet.Balance
	walletRecord.NotBilled = wallet.NotBilled

	if !must && wallet.NotBilled < 0 {
		return 0, InsufficientQuota.New()
	}

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)
		defrayModel := db.NewDefrayModelWithSession(session)
		websiteFundingModel := db.NewWebsiteFundingModelWithSession(session)

		_, err := websiteFundingModel.Insert(ctx, websiteFunding)
		if err != nil {
			return errors.WarpQuick(err)
		}

		_, err = walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = defrayModel.Update(ctx, defray)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, nil
}

func WaitBack(ctx context.Context, mysql sqlx.Session, user *db.User, canWithdraw bool, get int64, reason string) (int64, errors.WTError) {
	walletModel := db.NewWalletModelWithSession(mysql)

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletBack, "", fmt.Sprintf("返现即将入账: %s", reason))

	if canWithdraw {
		wallet.WaitWithdraw += get
	} else {
		wallet.WaitBalance += get
	}

	walletRecord.WaitBalance = wallet.Balance
	walletRecord.WaitWithdraw = wallet.Withdraw

	walletRecordModel := db.NewWalletRecordModelWithSession(mysql)
	_, err = walletRecordModel.Insert(ctx, &walletRecord)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	err = walletModel.Update(ctx, wallet)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, nil
}

func BackWithInsert(ctx context.Context, user *db.User, back *db.Back, reason string, mysql sqlx.Session) (int64, errors.WTError) {
	walletModel := db.NewWalletModelWithSession(mysql)

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if wallet.Id != back.WalletId {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletBack, back.BackId, fmt.Sprintf("返现入账: %s", reason))

	now := time.Now() // 此时还没有create_at

	websiteFunding := &db.WebsiteFunding{
		WebId:       back.SupplierId,
		Type:        db.WebsiteFundingBack,
		FundingId:   back.BackId,
		Profit:      back.Get,
		Expenditure: 0,
		Year:        int64(now.Year()),
		Month:       int64(now.Month()),
		Day:         int64(now.Day()),
		PayAt:       now,
	}

	if back.CanWithdraw {
		wallet.WaitWithdraw -= back.Get
		wallet.Withdraw += back.Get
		wallet.NotWithdraw += back.Get

		if wallet.WaitWithdraw < 0 {
			wallet.WaitWithdraw = 0
		}
	} else {
		wallet.WaitBalance -= back.Get
		wallet.Balance += back.Get

		if wallet.WaitBalance < 0 {
			wallet.WaitBalance = 0
		}
	}

	walletRecord.Balance = wallet.Balance
	walletRecord.WaitBalance = wallet.Balance
	walletRecord.Withdraw = wallet.Withdraw
	walletRecord.WaitWithdraw = wallet.Withdraw
	walletRecord.NotWithdraw = wallet.NotWithdraw

	if wallet.Balance < 0 {
		return 0, Insufficient.New()
	}

	back.Balance = wallet.Balance

	walletRecordModel := db.NewWalletRecordModelWithSession(mysql)
	backModel := db.NewBackModelWithSession(mysql)
	websiteFundingModel := db.NewWebsiteFundingModelWithSession(mysql)

	_, err = websiteFundingModel.Insert(ctx, websiteFunding)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	_, err = walletRecordModel.Insert(ctx, &walletRecord)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	res, err := backModel.Insert(ctx, back)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	back.Id, err = res.LastInsertId()
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	err = walletModel.Update(ctx, wallet)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, nil
}

func WithdrawWithInsert(ctx context.Context, user *db.User, withdraw *db.Withdraw) (int64, errors.WTError) {
	if withdraw.Status != db.WithdrawWait {
		return 0, errors.Errorf("double pay")
	}

	walletModel := db.NewWalletModel(mysql.MySQLConn)

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if withdraw.WalletId != wallet.Id {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletWithdraw, withdraw.WithdrawId, "用户提现")

	websiteFunding := &db.WebsiteFunding{
		WebId:       warp.UserCenterWebsite,
		Type:        db.WebsiteFundingWithdraw,
		FundingId:   withdraw.WithdrawId,
		Profit:      0,
		Expenditure: withdraw.Cny,
		Year:        int64(withdraw.WithdrawAt.Year()),
		Month:       int64(withdraw.WithdrawAt.Month()),
		Day:         int64(withdraw.WithdrawAt.Day()),
		PayAt:       withdraw.WithdrawAt,
	}

	wallet.NotWithdraw -= withdraw.Cny
	wallet.HasWithdraw += withdraw.Cny

	if wallet.NotWithdraw < 0 {
		return 0, Insufficient.New()
	}

	walletRecord.NotWithdraw = wallet.NotWithdraw
	walletRecord.HasWithdraw = wallet.HasWithdraw

	withdraw.Balance = sql.NullInt64{
		Valid: true,
		Int64: wallet.NotWithdraw,
	}

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)
		withdrawModel := db.NewWithdrawModelWithSession(session)
		websiteFundingModel := db.NewWebsiteFundingModelWithSession(session)

		_, err := websiteFundingModel.Insert(ctx, websiteFunding)
		if err != nil {
			return errors.WarpQuick(err)
		}

		_, err = walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		res, err := withdrawModel.Insert(ctx, withdraw)
		if err != nil {
			return errors.WarpQuick(err)
		}

		withdraw.Id, err = res.LastInsertId()
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, nil
}

func WithdrawFail(ctx context.Context, user *db.User, withdraw *db.Withdraw) (int64, errors.WTError) {
	if withdraw.Status != db.WithdrawWait {
		return 0, errors.Errorf("double pay")
	}

	walletModel := db.NewWalletModel(mysql.MySQLConn)

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if withdraw.WalletId != wallet.Id {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletWithdraw, withdraw.WithdrawId, "用户提现失败")

	now := time.Now()
	websiteFunding := &db.WebsiteFunding{
		WebId:       warp.UserCenterWebsite,
		Type:        db.WebsiteFundingWithdrawFail,
		FundingId:   withdraw.WithdrawId,
		Profit:      withdraw.Cny,
		Expenditure: 0,
		Year:        int64(now.Year()),
		Month:       int64(now.Month()),
		Day:         int64(now.Day()),
		PayAt:       now,
	}

	wallet.NotWithdraw += withdraw.Cny
	wallet.HasWithdraw -= withdraw.Cny

	walletRecord.NotWithdraw = wallet.NotWithdraw
	walletRecord.HasWithdraw = wallet.HasWithdraw

	withdraw.Status = db.WithdrawFail

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)
		withdrawModel := db.NewWithdrawModelWithSession(session)
		websiteFundingModel := db.NewWebsiteFundingModelWithSession(session)

		_, err := websiteFundingModel.Insert(ctx, websiteFunding)
		if err != nil {
			return errors.WarpQuick(err)
		}

		_, err = walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = withdrawModel.Update(ctx, withdraw)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.Balance, nil
}

func InvoiceWithInsert(ctx context.Context, user *db.User, invoice *db.Invoice) (int64, errors.WTError) {
	if invoice.Status != db.InvoiceWait {
		return 0, errors.Errorf("double invoice")
	}

	walletModel := db.NewWalletModel(mysql.MySQLConn)

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if wallet.Id != invoice.WalletId {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletInvoice, invoice.InvoiceId, "用户开票")

	wallet.NotBilled -= invoice.Amount
	wallet.HasBilled += invoice.Amount

	walletRecord.NotBilled = wallet.NotBilled
	walletRecord.HasBilled = wallet.HasBilled

	if wallet.NotBilled < 0 {
		return 0, InsufficientQuota.New()
	}

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)
		invoiceModel := db.NewInvoiceModelWithSession(session)

		_, err := walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		res, err := invoiceModel.Insert(ctx, invoice)
		if err != nil {
			return errors.WarpQuick(err)
		}

		invoice.Id, err = res.LastInsertId()
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.NotBilled, nil
}

func InvoiceReturn(ctx context.Context, user *db.User, invoice *db.Invoice, status int64) (int64, errors.WTError) {
	if status != db.InvoiceBad && status != db.InvoiceReturn && status != db.InvoiceRedFlush && status != db.InvoiceWaitReturn {
		return 0, errors.Errorf("bad status")
	}

	if invoice.Status != db.InvoiceOK && invoice.Status != db.InvoiceWait {
		return 0, errors.Errorf("bad status")
	}

	walletModel := db.NewWalletModel(mysql.MySQLConn)
	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if wallet.Id != invoice.WalletId {
		return 0, errors.Errorf("wallet not found")
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletInvoice, invoice.InvoiceId, "用户退票")

	wallet.NotBilled += invoice.Amount
	wallet.HasBilled -= invoice.Amount

	walletRecord.NotBilled = wallet.NotBilled
	walletRecord.HasBilled = wallet.HasBilled

	if wallet.HasBilled < 0 {
		wallet.HasBilled = 0
	}

	invoice.Status = status

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)
		invoiceModel := db.NewInvoiceModelWithSession(session)

		_, err := walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = invoiceModel.Update(ctx, invoice)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.NotBilled, nil
}

func InvoiceAdd(ctx context.Context, user *db.User, amount int64) (int64, errors.WTError) {
	walletModel := db.NewWalletModel(mysql.MySQLConn)

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletAdmin, "", "管理员增加可开票金额")

	wallet.Billed += amount
	wallet.NotBilled += amount

	walletRecord.Billed = wallet.Billed
	walletRecord.NotBilled = wallet.NotBilled

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)

		_, err := walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.NotBilled, nil
}

func InvoiceSub(ctx context.Context, user *db.User, amount int64) (int64, errors.WTError) {
	walletModel := db.NewWalletModel(mysql.MySQLConn)

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("wallet not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	}

	keyWallet := fmt.Sprintf("wallet:%d", wallet.Id)
	if !redis.AcquireLockMore(ctx, keyWallet, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(keyWallet)

	walletRecord := NewWalletRecord(wallet, user, db.WalletAdmin, "", "管理员减少可开票金额")

	wallet.Billed -= amount
	wallet.NotBilled -= amount

	walletRecord.Billed = wallet.Billed
	walletRecord.NotBilled = wallet.NotBilled

	err = mysql.MySQLConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		walletModel := db.NewWalletModelWithSession(session)
		walletRecordModel := db.NewWalletRecordModelWithSession(session)

		_, err := walletRecordModel.Insert(ctx, &walletRecord)
		if err != nil {
			return errors.WarpQuick(err)
		}

		err = walletModel.Update(ctx, wallet)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	go UpdateWallet(wallet, nil)
	return wallet.NotBilled, nil
}
