package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserTokenListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserTokenListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTokenListLogic {
	return &GetUserTokenListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserTokenListLogic) GetUserTokenList(req *types.AdminGetUserTokenReq) (resp *types.AdminGetAllTokenResp, err error) {
	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetAllTokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	}

	tokens, err := jwt.GetAllUserTokenGeo(l.ctx, srcUser.Uid)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	fatherTokens, err := jwt.GetAllUserFatherTokenGeo(l.ctx, srcUser.Uid)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	websiteTokens, err := jwt.GetAllWebsiteTokenGeo(l.ctx, srcUser.Uid)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	tokenResp := make([]types.AdminTokenIPGeo, 0, min(req.Limit, int64(len(tokens)+len(fatherTokens)+len(websiteTokens))))

	for _, t := range tokens {
		if int64(len(tokenResp)) > req.Limit {
			break
		}

		dt, err := jwt.CreateDeleteToken(t.Token, jwt.TypeDeleteUserToken)
		if err != nil {
			continue
		}

		isLogin := jwt.IsLoginToken(l.ctx, t.Token)
		if req.IsLogin && !isLogin {
			continue
		}

		tokenResp = append(tokenResp, types.AdminTokenIPGeo{
			TokenType:   jwt.TokenTypeUser,
			IP:          t.IP,
			Geo:         t.Geo,
			NowIP:       t.NowIP,
			NowGeo:      t.NowGeo,
			IsLogin:     isLogin,
			Token:       t.Token,
			DeleteToken: dt,
			SubType:     t.SubType,
		})
	}

	fatherMap := make(map[string]types.UserEasy, len(fatherTokens))

	for _, t := range fatherTokens {
		if int64(len(tokenResp)) > req.Limit {
			break
		}

		father, ok := fatherMap[t.Father]
		if !ok {
			father, err = action.GetUserEasy(l.ctx, 0, t.Father)
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			fatherMap[t.Father] = father
		}

		dt, err := jwt.CreateDeleteToken(t.Token, jwt.TypeDeleteUserToken)
		if err != nil {
			continue
		}

		isLogin := jwt.IsLoginToken(l.ctx, t.Token)
		if req.IsLogin && !isLogin {
			continue
		}

		tokenResp = append(tokenResp, types.AdminTokenIPGeo{
			TokenType:   jwt.TokenTypeFather,
			IP:          t.IP,
			Geo:         t.Geo,
			NowIP:       t.NowIP,
			NowGeo:      t.NowGeo,
			Father:      father,
			FatherToken: t.FatherToken,
			IsLogin:     isLogin,
			Token:       t.Token,
			DeleteToken: dt,
			SubType:     t.SubType,
		})
	}

	for _, t := range websiteTokens {
		if int64(len(tokenResp)) > req.Limit {
			break
		}

		web := action.GetWebsite(t.WebID)
		if web.Status == db.WebsiteStatusBanned {
			continue
		}

		dt, err := jwt.CreateDeleteToken(t.Token, jwt.TypeDeleteUserToken)
		if err != nil {
			continue
		}

		isLogin := jwt.IsLoginToken(l.ctx, t.Token)
		if req.IsLogin && !isLogin {
			continue
		}

		tokenResp = append(tokenResp, types.AdminTokenIPGeo{
			TokenType:   jwt.TokenTypeWebsite,
			IP:          t.IP,
			Geo:         t.Geo,
			NowIP:       t.NowIP,
			NowGeo:      t.NowGeo,
			WebID:       web.ID,
			WebName:     web.Name,
			IsLogin:     isLogin,
			Token:       t.Token,
			DeleteToken: dt,
			SubType:     t.SubType,
		})
	}

	return &types.AdminGetAllTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetAllTokenData{
			Token: tokenResp,
		},
	}, nil
}
