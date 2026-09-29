package balance

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"github.com/wuntsong-org/wterrors"
)

func QueryBalance(ctx context.Context, userID int64) (*db.Wallet, errors.WTError) {
	walletModel := db.NewWalletModel(mysql.MySQLConn)
	userModel := db.NewUserModel(mysql.MySQLConn)

	user, err := userModel.FindOneByIDWithoutDelete(ctx, userID)
	if errors.Is(err, db.ErrNotFound) {
		return nil, errors.Errorf("user not found")
	} else if err != nil {
		return nil, errors.WarpQuick(err)
	}

	wallet, err := walletModel.FindByWalletID(ctx, user.WalletId)
	if errors.Is(err, db.ErrNotFound) {
		return nil, errors.Errorf("wallet not found")
	} else if err != nil {
		return nil, errors.WarpQuick(err)
	}

	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return wallet, nil
}
