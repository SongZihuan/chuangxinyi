package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserFinanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserFinanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFinanceLogic {
	return &GetUserFinanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserFinanceLogic) GetUserFinance(req *types.AdminGetUserReq) (resp *types.AdminGetUserFinanceResp, err error) {
	user, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetUserFinanceResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	titleModel := db.NewTitleModel(mysql.MySQLConn)

	wallet, err := balance.QueryBalance(l.ctx, user.Id)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	title, err := titleModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		title = &db.Title{}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.AdminGetUserFinanceResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetUserFinanceData{
			WalletID:     wallet.Id,
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

			TitleName:   title.Name.String,
			TitleTaxID:  title.TaxId.String,
			TitleBankID: title.BankId.String,
			TitleBank:   title.Bank.String,
		},
	}, nil
}
