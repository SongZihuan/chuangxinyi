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

type UpdateDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDiscountLogic {
	return &UpdateDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDiscountLogic) UpdateDiscount(req *types.AdminUpdateDiscountReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	discountModel := db.NewDiscountModel(mysql.MySQLConn)
	d, err := discountModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DiscountNotFound, "优惠包未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

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

	d.Name = req.Name
	d.Describe = req.Describe
	d.ShortDescribe = req.ShortDescribe
	d.Type = req.Type
	d.Quota = string(quota)
	d.DayLimit = sql.NullInt64{
		Valid: req.DayLimit != 0,
		Int64: req.DayLimit,
	}
	d.MonthLimit = sql.NullInt64{
		Valid: req.MonthLimit != 0,
		Int64: req.MonthLimit,
	}
	d.YearLimit = sql.NullInt64{
		Valid: req.YearLimit != 0,
		Int64: req.YearLimit,
	}
	d.Limit = sql.NullInt64{
		Valid: req.Limit != 0,
		Int64: req.Limit,
	}
	d.NeedVerify = req.NeedVerify
	d.NeedCompany = req.NeedCompany
	d.NeedUserOrigin = req.NeedUserOrigin
	d.NeedCompanyOrigin = req.NeedCompanyOrigin
	d.NeedUserFace = req.NeedUserFace
	d.NeedCompanyFace = req.NeedCompanyFace
	d.Show = req.Show
	d.Remark = req.Remark

	err = discountModel.Update(l.ctx, d)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新优惠（%s）成功", d.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
