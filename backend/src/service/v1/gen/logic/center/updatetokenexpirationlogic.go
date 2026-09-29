package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateTokenExpirationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTokenExpirationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTokenExpirationLogic {
	return &UpdateTokenExpirationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTokenExpirationLogic) UpdateTokenExpiration(req *types.UserUpdateTokenExpirationReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.TokenExpiration == 0 || req.TokenExpiration > config.BackendConfig.Jwt.User.ExpiresSecond {
		req.TokenExpiration = config.BackendConfig.Jwt.User.ExpiresSecond
	}

	if time.Second*time.Duration(req.TokenExpiration) < time.Minute*15 {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TimeTooShort, "可登录时间太短"),
		}, nil
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	user.TokenExpiration = req.TokenExpiration
	err = userModel.UpdateWithoutStatus(l.ctx, user) // 不需要通知
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户更新登录在线最长时长成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
