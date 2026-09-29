package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/defray"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddDefrayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddDefrayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDefrayLogic {
	return &AddDefrayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddDefrayLogic) AddDefray(req *types.AdminAddDefrayReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	var ownerID = int64(0)
	if req.Owner.ID != 0 || len(req.Owner.UID) != 0 {
		owner, err := GetUser(l.ctx, req.Owner.ID, req.Owner.UID, true)
		if errors.Is(err, UserNotFound) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
		ownerID = owner.Id
	}

	err = defray.NewAdminDefray(l.ctx, defray.AdminDefrayData{
		OwnerID:            ownerID,
		Subject:            req.Subject,
		Price:              req.Price,
		Describe:           req.Describe,
		InvitePre:          req.InvitePre,
		DistributionLevel1: req.DistributionLevel1,
		DistributionLevel2: req.DistributionLevel2,
		DistributionLevel3: req.DistributionLevel3,
		CanWithdraw:        req.CanWithdraw,
		SupplierID:         web.ID,
	}, srcUser)
	if errors.Is(err, defray.Insufficient) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Insufficient, "添加消费失败，额度不足"),
		}, nil
	} else if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.ChangeBalanceFail, errors.WarpQuick(err), "添加消费失败"),
		}, nil
	}

	audit.NewAdminAudit(user.Id, "管理员新增用户消费：%s", srcUser.Uid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
