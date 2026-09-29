package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetInviteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInviteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInviteListLogic {
	return &GetInviteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInviteListLogic) GetInviteList(req *types.GetInviteList) (resp *types.GetInviteListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	inviteList, err := userModel.GetInviteUserEasyList(l.ctx, user.Id, req.Status, req.Src, req.Page, req.PageSize)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	inviteCount, err := userModel.CountInviteUserEasyList(l.ctx, user.Id, req.Status, req.Src)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.UserLessEasy, 0, len(inviteList))
	for _, i := range inviteList {
		easy, err := action.GetUserLessEasyOther(l.ctx, &i)
		if err != nil {
			continue
		}

		respList = append(respList, easy)
	}

	return &types.GetInviteListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetInviteListData{
			User:  respList,
			Count: inviteCount,
		},
	}, nil
}
