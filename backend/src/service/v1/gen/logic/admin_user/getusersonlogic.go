package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserSonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserSonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSonLogic {
	return &GetUserSonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSonLogic) GetUserSon(req *types.AdminGetUserSonReq) (resp *types.AdminGetSonTreeResp, err error) {
	user, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetSonTreeResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sons, err := GetSon(l.ctx, user.Id, 100, req.Status, req.Src)

	return &types.AdminGetSonTreeResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetSonTreeData{
			User: sons,
		},
	}, nil
}

func GetSon(ctx context.Context, id int64, n int64, status []string, searchSrc string) ([]types.AdminFatherUser, error) {
	if n == 0 {
		return []types.AdminFatherUser{}, nil
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	sons1, err := userModel.GetSonUserEasyList(ctx, id, status, searchSrc)
	if err != nil {
		return []types.AdminFatherUser{}, err
	}

	sons2, err := userModel.GetNephewUserEasyList(ctx, id, status, searchSrc)
	if err != nil {
		return []types.AdminFatherUser{}, err
	}

	res := make([]types.AdminFatherUser, 0, len(sons1)+len(sons2))

	for _, s := range sons1 {
		sSons, err := GetSon(ctx, s.ID, n-1, status, searchSrc)
		if err != nil {
			return []types.AdminFatherUser{}, err
		}

		r := action.GetRole(s.ID, s.IsAdmin)

		userEasy, err := action.GetUserEasyOther(ctx, &s)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		res = append(res, types.AdminFatherUser{
			UserID:       s.ID,
			UserEasy:     userEasy,
			RoleID:       r.ID,
			RoleName:     r.Name,
			RoleSign:     r.Sign,
			Lineal:       true,
			NephewStatus: 0,
			Son:          sSons,
		})
	}

	for _, s := range sons2 {
		r := action.GetRole(s.ID, s.IsAdmin)

		userEasy, err := action.GetUncleUserEasyOther(ctx, &s)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		res = append(res, types.AdminFatherUser{
			UserID:       s.ID,
			UserEasy:     userEasy,
			RoleID:       r.ID,
			RoleName:     r.Name,
			RoleSign:     r.Sign,
			Lineal:       false,
			NephewStatus: s.UncleStatus,
		})
	}

	return res, nil
}
