package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/google/uuid"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

type InvoiceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInvoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InvoiceLogic {
	return &InvoiceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InvoiceLogic) Invoice(req *types.InvoiceReq) (resp *types.InvoiceResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	titleModel := db.NewTitleModel(mysql.MySQLConn)
	addressModel := db.NewAddressModel(mysql.MySQLConn)

	InvoiceIDUUID, success := redis.GenerateUUIDMore(l.ctx, "invoice", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		invoiceModel := db.NewInvoiceModel(mysql.MySQLConn)
		_, err := invoiceModel.FindByInvoiceID(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return &types.InvoiceResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InvoiceFail, "生成发票uuid失败"),
		}, nil
	}

	InvoiceID := InvoiceIDUUID.String()

	title, err := titleModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.InvoiceResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TitleNotFound, "用户没设置抬头"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	address, err := addressModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.InvoiceResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.AddressNotFound, "用户没设置地址"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if req.Type == db.PersonalInvoice {
		if len(title.Name.String) == 0 || !utils.IsValidChineseName(title.Name.String) {
			return &types.InvoiceResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTitle, "抬头名称不正确"),
			}, nil
		}

		if len(title.TaxId.String) == 0 || !utils.IsValidIDCard(title.TaxId.String) {
			return &types.InvoiceResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTitle, "抬头税号（身份证号）不正确"),
			}, nil
		}
	} else if req.Type == db.CompanyInvoice || req.Type == db.CompanySpecializedInvoice {
		if len(title.Name.String) == 0 || !utils.IsValidChineseCompanyName(title.Name.String) {
			return &types.InvoiceResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTitle, "抬头名称不正确"),
			}, nil
		}

		if len(title.TaxId.String) == 0 || !utils.IsValidCreditCode(title.TaxId.String) {
			return &types.InvoiceResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTitle, "抬头税号（身份证号）不正确"),
			}, nil
		}
	} else {
		return &types.InvoiceResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InvoiceFail, "发票类型错误"),
		}, nil
	}

	invoice := &db.Invoice{
		UserId:    user.Id,
		WalletId:  user.WalletId,
		InvoiceId: InvoiceID,
		Amount:    req.Amount,
		Type:      req.Type,

		Name:   title.Name,
		TaxId:  title.TaxId,
		BankId: title.BankId,
		Bank:   title.Bank,

		Recipient: address.Name,
		Phone:     address.Phone,
		Email:     address.Email,
		Province:  address.Province,
		City:      address.City,
		District:  address.District,
		Address:   address.Address,
		Status:    db.InvoiceWait,
	}

	notBilled, err := balance.InvoiceWithInsert(l.ctx, user, invoice) // 自带insert
	if errors.Is(err, balance.InsufficientQuota) {
		return &types.InvoiceResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InsufficientQuota, "开票额度不足"),
		}, nil
	} else if err != nil {
		return &types.InvoiceResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.InvoiceFail, errors.WarpQuick(err), "开票错误"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户开票成功")

	return &types.InvoiceResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.InvoiceData{
			InvoiceID: InvoiceID,
			NotBilled: notBilled,
		},
	}, nil
}
