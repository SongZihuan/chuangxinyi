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

type GetLoginTokenListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoginTokenListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginTokenListLogic {
	return &GetLoginTokenListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginTokenListLogic) GetLoginTokenList(req *types.AdminGetUserTokenReq) (resp *types.AdminGetAllOauth2TokenResp, err error) {
	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetAllOauth2TokenResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	}

	ipgeo, err := jwt.GetAllLoginTokenGeo(l.ctx, srcUser.Uid)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	res := make([]types.AdminOauth2Record, 0, min(req.Limit, int64(len(ipgeo))))
	for _, i := range ipgeo {
		if int64(len(res)) > req.Limit {
			break
		}

		web := action.GetWebsite(i.WebID)
		if web.Status == db.WebsiteStatusBanned {
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

		res = append(res, types.AdminOauth2Record{
			UserID:           i.UserID,
			WebID:            web.ID,
			WebName:          web.Name,
			IP:               i.IP,
			Geo:              i.Geo,
			Token:            i.Token,
			DeleteToken:      dt,
			WebsiteUserToken: i.UserToken,
			IsLogin:          isLogin,
		})
	}

	return &types.AdminGetAllOauth2TokenResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetAllOauth2TokenData{
			Token: res,
		},
	}, nil
}
