package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListLogic {
	return &GetListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListLogic) GetList(req *types.GetOauth2ListReq) (resp *types.GetOauth2ListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	ipgeo, err := jwt.GetAllLoginTokenGeo(l.ctx, user.Uid)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	res := make([]types.Oauth2Record, 0, min(int64(len(ipgeo)), req.Limit))
	for _, i := range ipgeo {
		if int64(len(res)) >= req.Limit {
			break
		}

		web := action.GetWebsite(i.WebID)
		if web.ID == warp.UnknownWebsite {
			continue
		}

		dt, err := jwt.CreateDeleteToken(i.Token, jwt.TypeDeleteLoginToken)
		if err != nil {
			continue
		}

		isLogin := jwt.IsLoginToken(l.ctx, i.UserToken)

		if req.IsLogin && !isLogin {
			continue
		}

		res = append(res, types.Oauth2Record{
			WebID:       web.ID,
			WebName:     web.Name,
			IP:          i.IP,
			Geo:         i.Geo,
			DeleteToken: dt,
			IsLogin:     isLogin,
		})
	}

	return &types.GetOauth2ListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetOauth2ListData{
			Record: res,
		},
	}, nil
}
