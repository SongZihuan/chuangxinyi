package admin_user

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type ProcessInvoiceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessInvoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessInvoiceLogic {
	return &ProcessInvoiceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessInvoiceLogic) ProcessInvoice(req *types.AdminProcessInvoiceReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	invoiceModel := db.NewInvoiceModel(mysql.MySQLConn)
	invoice, err := invoiceModel.FindByInvoiceID(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InvoiceNotFound, "发票未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	keyInvoice := fmt.Sprintf("invoice:%s", invoice.InvoiceId)
	if !redis.AcquireLockMore(l.ctx, keyInvoice, time.Minute*2) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InvoiceDoubleOperation, "发票重复操作"),
		}, nil
	}
	defer redis.ReleaseLock(keyInvoice)

	if invoice.Status == db.InvoiceWait && req.Status == db.InvoiceOK {
		invoice.Status = req.Status
		invoice.BillingAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
	} else if (invoice.Status == db.InvoiceWait || invoice.Status == db.InvoiceOK) && req.Status == db.InvoiceWaitReturn {
		invoice.Status = req.Status
	} else if invoice.Status == db.InvoiceWaitReturn && (req.Status == db.InvoiceBad || req.Status == db.InvoiceReturn || req.Status == db.InvoiceRedFlush) {
		if req.Status == db.InvoiceReturn || req.Status == db.InvoiceRedFlush {
			invoice.ReturnAt = sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			}
		} else {
			invoice.ReturnAt = sql.NullTime{
				Valid: false,
			}
		}
	} else if (invoice.Status == db.InvoiceOK || invoice.Status == db.InvoiceWait) && (req.Status == db.InvoiceBad || req.Status == db.InvoiceReturn || req.Status == db.InvoiceRedFlush) {
		if req.Status == db.InvoiceReturn || req.Status == db.InvoiceRedFlush {
			invoice.ReturnAt = sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			}
		} else {
			invoice.ReturnAt = sql.NullTime{
				Valid: false,
			}
		}

		_, err = balance.InvoiceReturn(l.ctx, user, invoice, req.Status)
		if err != nil {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.InvoiceReturnFail, errors.WarpQuick(err), "退票失败"),
			}, nil
		}
	} else if (invoice.Status == db.InvoiceReturn || invoice.Status == db.InvoiceBad || invoice.Status == db.InvoiceRedFlush) && (req.Status == db.InvoiceBad || req.Status == db.InvoiceReturn || req.Status == db.InvoiceRedFlush) {
		invoice.Status = req.Status
		if (req.Status == db.InvoiceReturn || req.Status == db.InvoiceRedFlush) && !invoice.ReturnAt.Valid {
			invoice.ReturnAt = sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			}
		}
	} else {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.InvoiceOperationFail, "发票操作错误"),
		}, nil
	}

	err = invoiceModel.Update(l.ctx, invoice)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员处理开票（%s）", invoice.InvoiceId)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
