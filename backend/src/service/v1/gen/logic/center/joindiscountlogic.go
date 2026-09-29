package center

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

func (l *JoinDiscountLogic) JoinDiscount(req *types.JoinDiscountReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	discountModel := db.NewDiscountModel(mysql.MySQLConn)
	d, err := discountModel.FindOneWithoutDelete(l.ctx, req.DiscountID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DiscountNotFound, "优惠包未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if !d.Show {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DiscountNotFound, "优惠包非公开"),
		}, nil
	}

	_, err = discount.Join(l.ctx, user, d)
	switch true {
	case err == nil:
		// pass
	case errors.Is(err, discount.PurchaseLimit):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PurchaseLimit, "超过限额"),
		}, nil
	case errors.Is(err, discount.NeedVerify):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedVerify, "没有使用人实名"),
		}, nil
	case errors.Is(err, discount.NeedCompany):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedCompany, "没有企业实名"),
		}, nil
	case errors.Is(err, discount.NeedUserOrigin):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedUserOrigin, "没有上传使用人信息"),
		}, nil
	case errors.Is(err, discount.NeedCompanyOrigin):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedCompanyOrigin, "没有上传企业信息"),
		}, nil
	case errors.Is(err, discount.NeedUserFace):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedUserFace, "没有使用人人脸验证"),
		}, nil
	case errors.Is(err, discount.NeedCompanyFace):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NeedCompanyFace, "没有企业人脸验证"),
		}, nil
	default:
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.DiscountBuyFail, errors.WarpQuick(err), "加入优惠失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户加入优惠成功：%d", req.DiscountID)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
