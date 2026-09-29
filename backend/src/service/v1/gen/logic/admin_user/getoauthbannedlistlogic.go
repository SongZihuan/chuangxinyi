package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetOauthBannedListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOauthBannedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOauthBannedListLogic {
	return &GetOauthBannedListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOauthBannedListLogic) GetOauthBannedList(req *types.AdminGetOauthBannedListReq) (resp *types.AdminGetOauthBannedListResp, err error) {
	bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
	if req.ID == 0 && len(req.UID) == 0 {
		return &types.AdminGetOauthBannedListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	}

	limit := int64(len(model.Websites()))

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetOauthBannedListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	bannedList, err := bannedModel.GetList(l.ctx, srcUser.Id, limit)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	m := make(map[int64]types.AdminOauth2Banned, limit)
	for _, w := range model.WebsiteList() {
		if w.ID == warp.UnknownWebsite {
			continue
		}
		m[w.ID] = types.AdminOauth2Banned{
			UserID:      srcUser.Id,
			WebId:       w.ID,
			WebName:     w.Name,
			AllowLogin:  false,
			AllowDefray: false,
			AllowMsg:    false,
		}
	}

	for _, b := range bannedList {
		w, ok := m[b.WebId]
		if !ok {
			continue
		}

		w.AllowLogin = b.AllowLogin
		w.AllowDefray = b.AllowLogin && b.AllowDefray
		w.AllowMsg = b.AllowLogin && b.AllowMsg
		m[b.WebId] = w
	}

	res := make([]types.AdminOauth2Banned, 0, limit)
	for _, w := range m {
		res = append(res, w)
	}

	return &types.AdminGetOauthBannedListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetOauthBannedListData{
			Record: res,
		},
	}, nil
}
