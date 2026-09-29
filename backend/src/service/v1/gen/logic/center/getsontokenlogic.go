package center

import (
	"context"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetSonTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSonTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSonTokenLogic {
	return &GetSonTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSonTokenLogic) GetSonToken(req *types.GetSonTokenReq) (resp *types.SuccessResp, err error) {
	father, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	fatherToken, ok := l.ctx.Value("X-Token").(string)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	uncleModel := db.NewUncleModel(mysql.MySQLConn)
	son, err := utils2.FindUser(l.ctx, req.ID, false)
	if errors.Is(err, utils2.UserNotFound) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SonNotFound, "子用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	tmp := son
	for tmp.FatherId.Valid && tmp.FatherId.Int64 != father.Id && tmp.Id > tmp.FatherId.Int64 {
		// 向上寻找
		tmp, err = userModel.FindOneByIDWithoutDelete(l.ctx, tmp.FatherId.Int64)
		if errors.Is(err, db.ErrNotFound) {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SonNotFound, "子用户上级父用户未找到"),
			}, nil
		}
	}

	var token string
	var subType string
	if tmp.FatherId.Valid && tmp.FatherId.Int64 == father.Id {
		if father.FatherId.Valid {
			token, err = jwt.CreateUserToken(l.ctx, son.Uid, true, father.TokenExpiration, jwt2.UserFatherToken, fatherToken, 0)
			subType = jwt2.UserFatherTokenString
		} else {
			token, err = jwt.CreateUserToken(l.ctx, son.Uid, true, father.TokenExpiration, jwt2.UserRootFatherToken, fatherToken, 0)
			subType = jwt2.UserRootFatherTokenString
		}
	} else {
		u, err := uncleModel.FindByUserIDWithoutDelete(l.ctx, son.Id, father.Id)
		if errors.Is(err, db.ErrNotFound) {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SonNotFound, "协作账号未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		} else if u.Status != db.UncleOK {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SonNotFound, "协作人未确认"),
			}, nil
		}

		token, err = jwt.CreateUserToken(l.ctx, son.Uid, true, father.TokenExpiration, jwt2.UserUncleToken, fatherToken, 0)
		subType = jwt2.UserUncleTokenString
	}

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:    UserToken,
			Token:   token,
			SubType: subType,
		},
	}, nil
}
