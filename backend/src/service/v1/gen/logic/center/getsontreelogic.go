package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetSonTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSonTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSonTreeLogic {
	return &GetSonTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSonTreeLogic) GetSonTree(req *types.GetSonTree) (resp *types.GetSonTreeResp, err error) {
	father, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	son, err := GetSon(l.ctx, father.Id, userModel, config.BackendConfig.MySQL.SonUserLimit, req.Status, req.Src)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.GetSonTreeResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetSonTreeData{
			User: son,
		},
	}, nil
}

func GetSon(ctx context.Context, id int64, userModel db.UserModel, n int64, status []string, searchSrc string) ([]types.FatherUser, error) {
	if n < 0 {
		return []types.FatherUser{}, nil
	}

	sons1, err := userModel.GetSonUserEasyList(ctx, id, status, searchSrc)
	if err != nil {
		return []types.FatherUser{}, err
	}

	sons2, err := userModel.GetNephewUserEasyList(ctx, id, status, searchSrc)
	if err != nil {
		return []types.FatherUser{}, err
	}

	res := make([]types.FatherUser, 0, len(sons1)+len(sons2))
	for _, s := range sons1 {
		sSons, err := GetSon(ctx, s.ID, userModel, n-1, status, searchSrc)
		if err != nil {
			return []types.FatherUser{}, err
		}

		userEasy, err := action.GetUserEasyOther(ctx, &s)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		res = append(res, types.FatherUser{
			UserEasy:     userEasy,
			Lineal:       true,
			Son:          sSons,
			NephewStatus: 0,
		})
	}

	for _, s := range sons2 {
		userEasy, err := action.GetUncleUserEasyOther(ctx, &s)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		res = append(res, types.FatherUser{
			UserEasy:     userEasy,
			Lineal:       false,
			NephewStatus: s.UncleStatus,
		})
	}

	return res, nil
}
