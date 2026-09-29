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

type GetNephewListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNephewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNephewListLogic {
	return &GetNephewListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNephewListLogic) GetNephewList(req *types.GetUncleList) (resp *types.GetUncleListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	nephewList, err := userModel.GetNephewUserEasyList(l.ctx, user.Id, req.Status, req.Src)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.UncleUserEasy, 0, len(nephewList))
	for _, s := range nephewList {
		res := types.UncleUserEasy{
			UncleStatus: s.UncleStatus,
			UncleTag:    s.UncleTag,
		}

		uncle, err := action.GetUncleUserEasyOther(l.ctx, &s)
		if err != nil {
			continue
		}
		res.UserEasy = uncle

		respList = append(respList, res)
	}

	return &types.GetUncleListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetUncleListData{
			User: respList,
		},
	}, nil
}
