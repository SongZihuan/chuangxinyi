package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"sort"

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

func (l *GetOauthBannedListLogic) GetOauthBannedList(req *types.PageReq) (resp *types.GetOauthBannedListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	limit := int64(len(model.Websites()))

	bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
	bannedList, err := bannedModel.GetList(l.ctx, user.Id, limit)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	m := make(map[int64]types.Oauth2Banned, limit)
	for _, w := range model.WebsiteList() {
		if w.ID == warp.UnknownWebsite {
			continue
		}
		m[w.ID] = types.Oauth2Banned{
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

	targetRes := make([]types.Oauth2Banned, 0, limit)
	for _, w := range m {
		targetRes = append(targetRes, w)
	}

	sort.Slice(targetRes, func(i, j int) bool {
		return targetRes[i].WebId < targetRes[j].WebId
	})

	count := int64(len(targetRes))
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= count {
		return &types.GetOauthBannedListResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetOauthBannedListData{
				Record: []types.Oauth2Banned{},
				Count:  count,
			},
		}, nil
	}

	if end > count {
		end = count
	}

	return &types.GetOauthBannedListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetOauthBannedListData{
			Record: targetRes[start:end],
		},
	}, nil
}
