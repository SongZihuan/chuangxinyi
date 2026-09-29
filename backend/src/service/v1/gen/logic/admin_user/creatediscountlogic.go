package admin_user

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/discount"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CreateDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDiscountLogic {
	return &CreateDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDiscountLogic) CreateDiscount(req *types.AdminCreateDiscountReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	discountModel := db.NewDiscountModel(mysql.MySQLConn)
	quota, err := utils.JsonMarshal(discount.DiscountQuota{
		Amount:      req.Quota.Amount,
		CanWithdraw: req.Quota.CanWithdraw,

		Type:     req.Quota.Type,
		Send:     req.Quota.Send,
		Discount: req.Quota.Discount,
		Pre:      req.Quota.Pre,
	})
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadQuota, errors.WarpQuick(err), "编码quota错误"),
		}, nil
	}

	_, err = discountModel.Insert(l.ctx, &db.Discount{
		Name:          req.Name,
		Describe:      req.Describe,
		ShortDescribe: req.ShortDescribe,
		Type:          req.Type,
		Quota:         string(quota),
		DayLimit: sql.NullInt64{
			Valid: req.DayLimit != 0,
			Int64: req.DayLimit,
		},
		MonthLimit: sql.NullInt64{
			Valid: req.MonthLimit != 0,
			Int64: req.MonthLimit,
		},
		YearLimit: sql.NullInt64{
			Valid: req.YearLimit != 0,
			Int64: req.DayLimit,
		},
		Limit: sql.NullInt64{
			Valid: req.Limit != 0,
			Int64: req.Limit,
		},
		NeedVerify:        req.NeedVerify,
		NeedCompany:       req.NeedCompany,
		NeedUserOrigin:    req.NeedUserOrigin,
		NeedCompanyOrigin: req.NeedCompanyOrigin,
		NeedUserFace:      req.NeedUserFace,
		NeedCompanyFace:   req.NeedCompanyFace,
		Show:              req.Show,
		Remark:            req.Remark,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员创建优惠（%s）成功", req.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
