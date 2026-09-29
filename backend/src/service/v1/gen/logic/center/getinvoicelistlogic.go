package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetInvoiceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInvoiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInvoiceListLogic {
	return &GetInvoiceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInvoiceListLogic) GetInvoiceList(req *types.GetInvoiceListReq) (resp *types.GetInvoiceListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	invoiceModel := db.NewInvoiceModel(mysql.MySQLConn)
	invoiceList, err := invoiceModel.GetList(l.ctx, user.WalletId, req.Type, req.Status, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := invoiceModel.GetCount(l.ctx, user.WalletId, req.Type, req.Status, req.Src, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userMap := make(map[int64]types.UserEasy, len(invoiceList))

	respList := make([]types.Invoice, 0, len(invoiceList))
	for _, i := range invoiceList {
		optUser, ok := userMap[i.UserId]
		if !ok {
			optUser, err = action.GetUserEasy(l.ctx, i.UserId, "")
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			userMap[i.UserId] = optUser
		}

		billingAt := int64(0)
		if i.BillingAt.Valid {
			billingAt = i.BillingAt.Time.Unix()
		}

		returnAt := int64(0)
		if i.ReturnAt.Valid {
			returnAt = i.ReturnAt.Time.Unix()
		}

		issuerAt := int64(0)
		if i.IssuerAt.Valid {
			issuerAt = i.IssuerAt.Time.Unix()
		}

		redIssuerAt := int64(0)
		if i.RedIssuerAt.Valid {
			redIssuerAt = i.RedIssuerAt.Time.Unix()
		}

		respList = append(respList, types.Invoice{
			InvoiceID: i.InvoiceId,
			Type:      i.Type,
			Title: types.Title{
				Name:   i.Name.String,
				TaxID:  i.TaxId.String,
				BankID: i.BankId.String,
				Bank:   i.Bank.String,
			},
			Address: types.InvoiceAddress{
				Name:     i.Recipient.String,
				Phone:    i.Phone.String,
				Email:    i.Email.String,
				Province: i.Province.String,
				City:     i.City.String,
				District: i.District.String,
				Address:  i.Address.String,
			},
			Amount:      i.Amount,
			Status:      i.Status,
			CreateAt:    i.CreateAt.Unix(),
			BillingAt:   billingAt,
			ReturnAt:    returnAt,
			IssuerAt:    issuerAt,
			RedIssuerAt: redIssuerAt,
		})
	}

	return &types.GetInvoiceListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetInvoiceListData{
			Count:   count,
			Invoice: respList,
		},
	}, nil
}
