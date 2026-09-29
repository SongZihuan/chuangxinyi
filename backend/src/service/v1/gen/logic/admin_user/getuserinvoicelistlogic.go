package admin_user

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

type GetUserInvoiceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInvoiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInvoiceListLogic {
	return &GetUserInvoiceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInvoiceListLogic) GetUserInvoiceList(req *types.AdminGetInvoiceListReq) (resp *types.AdminGetInvoiceListResp, err error) {
	var invoiceList []db.Invoice
	var count int64

	invoiceModel := db.NewInvoiceModel(mysql.MySQLConn)
	if req.ID == 0 && len(req.UID) == 0 {
		invoiceList, err = invoiceModel.GetList(l.ctx, 0, req.Type, req.Status, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = invoiceModel.GetCount(l.ctx, 0, req.Type, req.Status, req.Src, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		user, err := GetUser(l.ctx, req.ID, req.UID, true)
		if errors.Is(err, UserNotFound) {
			return &types.AdminGetInvoiceListResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		invoiceList, err = invoiceModel.GetList(l.ctx, user.WalletId, req.Type, req.Status, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = invoiceModel.GetCount(l.ctx, user.WalletId, req.Type, req.Status, req.Src, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	userMap := make(map[int64]types.UserEasy, len(invoiceList))

	respList := make([]types.AdminInvoice, 0, len(invoiceList))
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

		respList = append(respList, types.AdminInvoice{
			UserID:    i.UserId,
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

			InvoiceNumber:       i.InvoiceNumber.String,
			InvoiceCode:         i.InvoiceCode.String,
			InvoiceCheckCdoe:    i.InvoiceCheckCode.String,
			RedInvoiceNumber:    i.RedInvoiceNumber.String,
			RedInvoiceCode:      i.RedInvoiceCode.String,
			RedInvoiceCheckCode: i.RedInvoiceCheckCode.String,

			Remark:      i.Remark,
			Amount:      i.Amount,
			Status:      i.Status,
			CreateAt:    i.CreateAt.Unix(),
			BillingAt:   billingAt,
			ReturnAt:    returnAt,
			IssuerAt:    issuerAt,
			RedIssuerAt: redIssuerAt,
		})
	}

	return &types.AdminGetInvoiceListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetInvoiceListData{
			Count:   count,
			Invoice: respList,
		},
	}, nil
}
