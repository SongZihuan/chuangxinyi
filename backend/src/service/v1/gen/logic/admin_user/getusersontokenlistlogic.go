package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserSonTokenListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserSonTokenListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSonTokenListLogic {
	return &GetUserSonTokenListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSonTokenListLogic) GetUserSonTokenList(req *types.AdminGetUserTokenReq) (resp *types.AdminGetAllSonTokenResp, err error) {
	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetAllSonTokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	}

	tokens, err := jwt.GetAllUserSonTokenGeo(l.ctx, srcUser.Uid)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	sonMap := make(map[string]types.UserEasy, min(req.Limit, int64(len(tokens))))

	tokenResp := make([]types.AdminSonTokenIPGeo, 0, len(tokens))
	for _, t := range tokens {
		if int64(len(tokenResp)) > req.Limit {
			break
		}

		son, ok := sonMap[t.UserID]
		if !ok {
			son, err = action.GetUserEasy(l.ctx, 0, t.UserID)
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			sonMap[t.UserID] = son
		}

		isLogin := jwt.IsLoginToken(l.ctx, t.Token)
		if req.IsLogin && !isLogin {
			continue
		}

		tokenResp = append(tokenResp, types.AdminSonTokenIPGeo{
			UserID:      t.UserID,
			IP:          t.IP,
			Geo:         t.Geo,
			NowIP:       t.NowIP,
			NowGeo:      t.NowGeo,
			User:        son,
			IsLogin:     isLogin,
			Token:       t.Token,
			DeleteToken: "",
			SubType:     t.SubType,
		})
	}

	return &types.AdminGetAllSonTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetAllSonTokenData{
			Token: tokenResp,
		},
	}, nil
}
