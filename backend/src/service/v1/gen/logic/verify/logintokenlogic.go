package verify

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type LoginTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginTokenLogic {
	return &LoginTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginTokenLogic) LoginToken(req *types.CheckLoginTokenReq) (resp *types.CheckLoginTokenResp, err error) {
	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	loginData, err := jwt.ParserLoginToken(l.ctx, req.Token)
	if err != nil || loginData.WebID != web.ID {
		return &types.CheckLoginTokenResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析Token失败"),
			Data: types.CheckLoginTokenData{
				IsLogin: false,
			},
		}, nil
	}

	userModel := db.NewUserModel(mysql.MySQLConn)

	user, err := userModel.FindOneByUidWithoutDelete(l.ctx, loginData.UserID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &types.CheckLoginTokenResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
				Data: types.CheckLoginTokenData{
					IsLogin: false,
				},
			}, nil
		}
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if db.IsBanned(user) {
		return &types.CheckLoginTokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			Data: types.CheckLoginTokenData{
				IsLogin: false,
			},
		}, nil
	}

	data, err := utils.GetUserInfoWebsite(l.ctx, user)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.CheckLoginTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.CheckLoginTokenData{
			IsLogin: true,
			User:    data.User,
			Info:    data.InfoEasy,
			Data:    data.Data,
		},
	}, nil
}
