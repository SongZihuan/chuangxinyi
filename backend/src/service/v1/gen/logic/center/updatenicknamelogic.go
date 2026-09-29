package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
)

type UpdateNicknameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNicknameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNicknameLogic {
	return &UpdateNicknameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNicknameLogic) UpdateNickname(req *types.UserUpdateNicknameReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.IsDelete || len(req.NickName) == 0 {
		nickNameModel := db.NewNicknameModel(mysql.MySQLConn)
		_, err = nickNameModel.InsertWithDelete(context.Background(), &db.Nickname{
			UserId: user.Id,
			Nickname: sql.NullString{
				Valid: false,
			},
		})

		audit.NewUserAudit(user.Id, "用户删除昵称成功")
		return &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		}, nil
	}

	res, err := yundun.CheckName(req.NickName)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadNickName, errors.WarpQuick(err), "昵称审核失败"),
		}, nil
	}

	if !res {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadNickName, "昵称审核不通过"),
		}, nil
	}

	nickNameModel := db.NewNicknameModel(mysql.MySQLConn)

	_, err = nickNameModel.InsertWithDelete(context.Background(), &db.Nickname{
		UserId: user.Id,
		Nickname: sql.NullString{
			Valid:  true,
			String: req.NickName,
		},
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户更新昵称成功")
	sender.MessageSend(user.Id, "昵称更新", "昵称已经更新成功！")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
