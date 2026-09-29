package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/discount"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type JoinDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinDiscountLogic {
	return &JoinDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinDiscountLogic) JoinDiscount(req *types.AdminJoinDiscountReq) (resp *types.RespEmpty, err error) {
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

	discountModel := db.NewDiscountModel(mysql.MySQLConn)
	d, err := discountModel.FindOneWithoutDelete(l.ctx, req.DiscountID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DiscountNotFound, "优惠包未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	_, err = discount.Join(l.ctx, srcUser, d)
	switch true {
	case err == nil:
		// pass
	case errors.Is(err, discount.PurchaseLimit):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PurchaseLimit, "超过使用限额"),
		}, nil
	case errors.Is(err, discount.NeedVerify):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedVerify, "需要使用人实名"),
		}, nil
	case errors.Is(err, discount.NeedCompany):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedCompany, "需要企业实名"),
		}, nil
	case errors.Is(err, discount.NeedUserOrigin):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedUserOrigin, "需要使用人上传信息"),
		}, nil
	case errors.Is(err, discount.NeedCompanyOrigin):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedCompanyOrigin, "需要企业上传信息"),
		}, nil
	case errors.Is(err, discount.NeedUserFace):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedUserFace, "需要使用人人脸遇难者"),
		}, nil
	case errors.Is(err, discount.NeedCompanyFace):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedCompanyFace, "需要企业人脸验证"),
		}, nil
	default:
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.DiscountBuyFail, errors.WarpQuick(err), "将用户加入优惠包失败"),
		}, nil
	}

	audit.NewAdminAudit(user.Id, "管理员将用户（%s）加入优惠（%s）", srcUser.Uid, d.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
