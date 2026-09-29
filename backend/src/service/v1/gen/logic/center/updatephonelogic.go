package center

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
	"time"

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

func (l *UpdatePhoneLogic) UpdatePhone(req *types.UserUpdatePhoneReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	phone, err := jwt.ParserPhoneToken(req.PhoneToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	} else if phone.WebID != warp.UserCenterWebsite {
		return nil, respmsg.JWTError.New("bad website")
	}

	lock := fmt.Sprintf("phone:%s", phone)
	if !redis.AcquireLock(l.ctx, lock, time.Minute*2) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneHasBeenRegister, "上锁失败，手机号重复"),
		}, nil
	}
	defer redis.ReleaseLock(lock)

	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	_, err = phoneModel.FindByPhone(l.ctx, phone.Phone)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if err == nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneHasBeenRegister, "数据库存在手机号，手机号重复"),
		}, nil
	}

	oldPhone, err := phoneModel.FindByUserID(context.Background(), user.Id)
	if errors.Is(err, db.ErrNotFound) {
		logger.Logger.Error("user not phone: %d, %s", user.Id, user.Uid)
		oldPhone = &db.Phone{
			Phone: "",
		}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		phoneModel := db.NewPhoneModelWithSession(session)
		_, err := phoneModel.InsertSafe(l.ctx, &db.Phone{
			UserId: user.Id,
			Phone:  phone.Phone,
		})
		return err
	})
	if errors.Is(err, db.PhoneRepeat) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PhoneRepeat, "手机号重复"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.PhoneSendPhoneChange(oldPhone.Phone, phone.Phone)
	sender.EmailSendChange(user.Id, "手机号")
	sender.MessageSendChange(user.Id, "手机号")
	sender.WxrobotSendChange(user.Id, "手机号")
	sender.FuwuhaoSendChange(user.Id, "手机号")
	sender.PhoneSendBind(phone.Phone)
	audit.NewUserAudit(user.Id, "用户更新手机号成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
