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

type GetInvoiceInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInvoiceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInvoiceInfoLogic {
	return &GetInvoiceInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInvoiceInfoLogic) GetInvoiceInfo(req *types.AdminGetInvoiceInfoReq) (resp *types.AdminGetInvoiceInfoResp, err error) {
	invoiceModel := db.NewInvoiceModel(mysql.MySQLConn)
	invoice, err := invoiceModel.FindByInvoiceID(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.AdminGetInvoiceInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InvoiceNotFound, "发票未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	billingAt := int64(0)
	if invoice.BillingAt.Valid {
		billingAt = invoice.BillingAt.Time.Unix()
	}

	returnAt := int64(0)
	if invoice.ReturnAt.Valid {
		returnAt = invoice.ReturnAt.Time.Unix()
	}

	issuerAt := int64(0)
	if invoice.IssuerAt.Valid {
		issuerAt = invoice.IssuerAt.Time.Unix()
	}

	redIssuerAt := int64(0)
	if invoice.RedIssuerAt.Valid {
		redIssuerAt = invoice.RedIssuerAt.Time.Unix()
	}

	optUser, err := action.GetUserEasy(l.ctx, invoice.UserId, "")
	if errors.Is(err, action.UserEasyNotFound) {
		return &types.AdminGetInvoiceInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.AdminGetInvoiceInfoResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetInvoiceInfoData{
			Invoice: types.AdminInvoice{
				UserID:    invoice.UserId,
				User:      optUser,
				WalletID:  invoice.WalletId,
				InvoiceID: invoice.InvoiceId,
				Type:      invoice.Type,
				Title: types.Title{
					Name:   invoice.Name.String,
					TaxID:  invoice.TaxId.String,
					BankID: invoice.BankId.String,
					Bank:   invoice.Bank.String,
				},
				Address: types.InvoiceAddress{
					Name:     invoice.Recipient.String,
					Phone:    invoice.Phone.String,
					Email:    invoice.Email.String,
					Province: invoice.Province.String,
					City:     invoice.City.String,
					District: invoice.District.String,
					Address:  invoice.Address.String,
				},
				InvoiceNumber:       invoice.InvoiceNumber.String,
				InvoiceCode:         invoice.InvoiceCode.String,
				InvoiceCheckCdoe:    invoice.InvoiceCheckCode.String,
				RedInvoiceNumber:    invoice.RedInvoiceNumber.String,
				RedInvoiceCode:      invoice.RedInvoiceCode.String,
				RedInvoiceCheckCode: invoice.RedInvoiceCheckCode.String,

				Remark:      invoice.Remark,
				Amount:      invoice.Amount,
				Status:      invoice.Status,
				CreateAt:    invoice.CreateAt.Unix(),
				BillingAt:   billingAt,
				ReturnAt:    returnAt,
				IssuerAt:    issuerAt,
				RedIssuerAt: redIssuerAt,
			},
		},
	}, nil
}
