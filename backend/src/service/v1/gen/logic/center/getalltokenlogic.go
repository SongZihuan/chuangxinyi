package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAllTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllTokenLogic {
	return &GetAllTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllTokenLogic) GetAllToken() (resp *types.GetAllTokenResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	token, ok := l.ctx.Value("X-Token").(string)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token")
	}

	tokens, err := jwt.GetAllUserTokenGeo(l.ctx, user.Uid)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	fatherTokens, err := jwt.GetAllUserFatherTokenGeo(l.ctx, user.Uid)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	websiteTokens, err := jwt.GetAllWebsiteTokenGeo(l.ctx, user.Uid)
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	tokenResp := make([]types.TokenIPGeo, 0, len(tokens)+len(fatherTokens)+len(websiteTokens))

	for _, t := range tokens {
		dt, err := jwt.CreateDeleteToken(t.Token, jwt.TypeDeleteUserToken)
		if err != nil {
			continue
		}

		tokenResp = append(tokenResp, types.TokenIPGeo{
			TokenType:   jwt.TokenTypeUser,
			IP:          t.IP,
			Geo:         t.Geo,
			NowIP:       t.NowIP,
			NowGeo:      t.NowGeo,
			IsLogin:     jwt.IsLoginToken(l.ctx, t.Token),
			DeleteToken: dt,
			IsSelf:      t.Token == token,
		})
	}

	fatherMap := make(map[string]types.UserEasy, len(fatherTokens))

	for _, t := range fatherTokens {
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

		tokenResp = append(tokenResp, types.TokenIPGeo{
			TokenType:   jwt.TokenTypeFather,
			IP:          t.IP,
			Geo:         t.Geo,
			NowIP:       t.NowIP,
			NowGeo:      t.NowGeo,
			Father:      father,
			IsLogin:     jwt.IsLoginToken(l.ctx, t.Token),
			DeleteToken: "",
			IsSelf:      t.Token == token,
		})
	}

	for _, t := range websiteTokens {
		web := action.GetWebsite(t.WebID)
		if web.Status == db.WebsiteStatusBanned {
			continue
		}

		dt, err := jwt.CreateDeleteToken(t.Token, jwt.TypeDeleteUserToken)
		if err != nil {
			continue
		}

		tokenResp = append(tokenResp, types.TokenIPGeo{
			TokenType:   jwt.TokenTypeWebsite,
			IP:          t.IP,
			Geo:         t.Geo,
			NowIP:       t.NowIP,
			NowGeo:      t.NowGeo,
			WebID:       web.ID,
			WebName:     web.Name,
			IsLogin:     jwt.IsLoginToken(l.ctx, t.Token),
			DeleteToken: dt,
			IsSelf:      t.Token == token,
		})
	}

	return &types.GetAllTokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetAllTokenData{
			Token: tokenResp,
		},
	}, nil
}
