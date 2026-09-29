package admin_user

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

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.AdminGetUserListReq) (resp *types.AdminGetUserListResp, err error) {
	userModel := db.NewUserModel(mysql.MySQLConn)
	userList, err := userModel.GetUserEasyList(l.ctx, req.Status, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := userModel.GetUserEasyCount(l.ctx, req.Status, req.Src, req.StartTime, req.EndTime)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userListResp := make([]types.UserEasyWithID, 0, len(userList))
	for _, u := range userList {
		userEasy, err := action.GetUserEasyOther(l.ctx, &u)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		userListResp = append(userListResp, types.UserEasyWithID{
			NumberID: u.ID,
			UserEasy: userEasy,
		})
	}

	return &types.AdminGetUserListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetUserListData{
			User:  userListResp,
			Count: count,
		},
	}, nil
}
