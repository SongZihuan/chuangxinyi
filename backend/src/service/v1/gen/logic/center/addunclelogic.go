package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddUncleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUncleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUncleLogic {
	return &AddUncleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUncleLogic) AddUncle(req *types.AddUncleReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	uncleModel := db.NewUncleModel(mysql.MySQLConn)

	uncle, err := utils2.FindUser(l.ctx, req.UncleID, false)
	if errors.Is(err, utils2.UserNotFound) || uncle.Id == user.Id {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "协作人未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	countUncle, err := uncleModel.GetUncleCount(l.ctx, user.Id)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if countUncle > config.BackendConfig.MySQL.UncleUserLimit {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "协作人过多"),
		}, nil
	}

	countNephew, err := uncleModel.GetNephewCount(l.ctx, uncle.Id)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if countNephew > config.BackendConfig.MySQL.NephewLimit {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "侄账号过多"),
		}, nil
	}

	_, err = uncleModel.FindByUserIDWithoutDelete(l.ctx, user.Id, uncle.Id)
	if !errors.Is(err, db.ErrNotFound) {
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.RepeatAddUncle, "用户已经添加对方为协作人"),
		}, nil
	}

	_, err = uncleModel.Insert(l.ctx, &db.Uncle{
		UserId:   user.Id,
		UncleId:  uncle.Id,
		UncleTag: req.UncleID,
		Status:   db.UncleWaitOk,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.MessageSendUncle(user.Uid, uncle.Id)
	audit.NewUserAudit(user.Id, "用户添加协作人：%s", uncle.Uid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
