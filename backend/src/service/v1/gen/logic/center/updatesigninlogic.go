package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateSigninLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSigninLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSigninLogic {
	return &UpdateSigninLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSigninLogic) UpdateSignin(req *types.UserUpdateSigninReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	selfToken, ok := l.ctx.Value("X-Token").(string)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token")
	}

	if !user.Signin && req.Signin { // 从非单点登录切换到单点登录
		err = jwt.DeleteAllUserToken(l.ctx, user.Uid, selfToken) // 删除全部token 除了selfToken
		if err != nil {
			return nil, respmsg.RedisSystemError.WarpQuick(err)
		}
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	user.Signin = req.Signin
	err = userModel.UpdateWithoutStatus(l.ctx, user) // 不需要通知
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户更新单点登录成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
