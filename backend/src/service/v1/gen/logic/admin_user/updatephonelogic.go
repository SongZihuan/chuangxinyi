package admin_user

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdatePhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhoneLogic {
	return &UpdatePhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePhoneLogic) UpdatePhone(req *types.AdminUpdatePhoneReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	phoneModel := db.NewPhoneModel(mysql.MySQLConn)

	oldPhone, err := phoneModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		oldPhone = &db.Phone{
			Phone: "",
		}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if !db.IsBanned(srcUser) {
		lock := fmt.Sprintf("phone:%s", req.Phone)
		if !redis.AcquireLock(l.ctx, lock, time.Minute*2) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneHasBeenRegister, "无法上锁，手机号已被注册"),
			}, nil
		}
		defer redis.ReleaseLock(lock)

		_, err = phoneModel.FindByPhone(l.ctx, req.Phone)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		} else if err == nil {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneHasBeenRegister, "手机号已经存在数据库中"),
			}, nil
		}
	}

	_, err = phoneModel.InsertSafe(l.ctx, &db.Phone{
		UserId:   srcUser.Id,
		Phone:    req.Phone,
		IsDelete: db.IsBanned(srcUser),
	})
	if errors.Is(err, db.PhoneRepeat) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneRepeat, "手机号重复"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新用户（%s）手机成功", srcUser.Uid)

	if !db.IsBanned(srcUser) {
		sender.PhoneSendPhoneChange(oldPhone.Phone, req.Phone)
		sender.EmailSendChange(user.Id, "手机号")
		sender.MessageSendChange(user.Id, "手机号")
		sender.WxrobotSendChange(user.Id, "手机号")
		sender.FuwuhaoSendChange(user.Id, "手机号")
		sender.PhoneSendBind(req.Phone)
		audit.NewUserAudit(user.Id, "用户更新手机号成功")
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
