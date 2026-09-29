package admin_user

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/google/uuid"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddInvoiceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddInvoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddInvoiceLogic {
	return &AddInvoiceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddInvoiceLogic) AddInvoice(req *types.AdminAddInvoice) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	InvoiceIDUUID, success := redis.GenerateUUIDMore(l.ctx, "invoice", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		invoiceModel := db.NewInvoiceModel(mysql.MySQLConn)
		_, err := invoiceModel.FindByInvoiceID(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InvoiceFail, "生成发票id失败"),
		}, nil
	}

	InvoiceID := InvoiceIDUUID.String()

	if req.Type == db.PersonalInvoice {
		if len(req.Name) == 0 || !utils.IsValidChineseName(req.Name) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTitle, "错误的发票抬头名字"),
			}, nil
		}

		if len(req.TaxId) == 0 || !utils.IsValidIDCard(req.TaxId) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTitle, "错误的发票税号（身份证号）"),
			}, nil
		}
	} else if req.Type == db.CompanyInvoice || req.Type == db.CompanySpecializedInvoice {
		if len(req.Name) == 0 || !utils.IsValidChineseCompanyName(req.Name) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTitle, "错误的发票抬头名字"),
			}, nil
		}

		if len(req.TaxId) == 0 || !utils.IsValidCreditCode(req.TaxId) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTitle, "错误的发票税号（企业税号）"),
			}, nil
		}
	} else {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InvoiceFail, "错误的发票类型"),
		}, nil
	}

	// 不检查邮箱和手机

	invoice := &db.Invoice{
		UserId:    srcUser.Id,
		WalletId:  srcUser.WalletId,
		InvoiceId: InvoiceID,
		Amount:    req.Amount,
		Type:      req.Type,

		Name: sql.NullString{
			Valid:  len(req.Name) != 0,
			String: req.Name,
		},
		TaxId: sql.NullString{
			Valid:  len(req.TaxId) != 0,
			String: req.TaxId,
		},
		BankId: sql.NullString{
			Valid:  len(req.BankId) != 0,
			String: req.BankId,
		},
		Bank: sql.NullString{
			Valid:  len(req.Bank) != 0,
			String: req.Bank,
		},
		Recipient: sql.NullString{
			Valid:  len(req.Recipient) != 0,
			String: req.Recipient,
		},
		Phone: sql.NullString{
			Valid:  len(req.Phone) != 0,
			String: req.Phone,
		},
		Email: sql.NullString{
			Valid:  len(req.Email) != 0,
			String: req.Email,
		},
		Province: sql.NullString{
			Valid:  len(req.Province) != 0,
			String: req.Province,
		},
		City: sql.NullString{
			Valid:  len(req.City) != 0,
			String: req.City,
		},
		District: sql.NullString{
			Valid:  len(req.District) != 0,
			String: req.District,
		},
		Address: sql.NullString{
			Valid:  len(req.Address) != 0,
			String: req.Address,
		},
		Status: db.InvoiceWait,
	}

	_, err = balance.InvoiceWithInsert(l.ctx, srcUser, invoice) // 自带insert
	if errors.Is(err, balance.InsufficientQuota) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InsufficientQuota, "开票失败，额度不足"),
		}, nil
	} else if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.InvoiceFail, errors.WarpQuick(err), "开票失败"),
		}, nil
	}

	audit.NewAdminAudit(user.Id, "管理员新增用户发票：%s", srcUser.Uid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
