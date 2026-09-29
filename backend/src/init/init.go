package init

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/afs"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/auth"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/dbinit"
	"gitee.com/wuntsong-auth/backend/src/defray"
	"gitee.com/wuntsong-auth/backend/src/email"
	"gitee.com/wuntsong-auth/backend/src/fuwuhao"
	"gitee.com/wuntsong-auth/backend/src/global/peername"
	"gitee.com/wuntsong-auth/backend/src/ip"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/ocr"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/rand"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/signalexit"
	"gitee.com/wuntsong-auth/backend/src/sms"
	"gitee.com/wuntsong-auth/backend/src/wechat"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	errors "github.com/wuntsong-org/wterrors"
	"os"
)

func InitUserCenter() errors.WTError {
	var err errors.WTError

	err = signalexit.InitSignalExit(0)
	if err != nil {
		return errors.Warp(err, "signal error")
	}

	signalexit.AddExitByFunc(func(ctx context.Context, _ os.Signal) context.Context {
		CloseAll()
		return context.WithValue(ctx, "InitClose", true)
	})

	err = rand.InitRander()
	if err != nil {
		return errors.Warp(err, "rander error")
	}

	err = peername.InitPeerName(config.EnvPrefix)
	if err != nil {
		return errors.Warp(err, "redis error")
	}

	err = redis.InitRedis()
	if err != nil {
		return errors.Warp(err, "redis error")
	}

	err = mysql.InitMysql()
	if err != nil {
		return errors.Warp(err, "mysql error")
	}

	err = sms.InitSMS()
	if err != nil {
		return errors.Warp(err, "sms error")
	}

	err = email.InitSmtp()
	if err != nil {
		return errors.Warp(err, "email error")
	}

	err = jwt.InitJWT()
	if err != nil {
		return errors.Warp(err, "jwt error")
	}

	err = logger.InitLogger(config.BackendConfig.User.LogServiceName)
	if err != nil {
		return errors.Warp(err, "logger error")
	}

	err = ocr.InitOcr()
	if err != nil {
		return errors.Warp(err, "ocr error")
	}

	err = yundun.InitYunDun()
	if err != nil {
		return errors.Warp(err, "yundun error")
	}

	err = oss.InitOss()
	if err != nil {
		return errors.Warp(err, "oss error")
	}

	err = cron.InitCron()
	if err != nil {
		return errors.Warp(err, "cron error")
	}

	err = afs.InitAFS()
	if err != nil {
		return errors.Warp(err, "afs error")
	}

	err = wechat.InitWeChat()
	if err != nil {
		return errors.Warp(err, "wechat error")
	}

	err = alipay.InitAlipay()
	if err != nil {
		return errors.Warp(err, "alipay error")
	}

	err = wechatpay.InitWeChatPay()
	if err != nil {
		return errors.Warp(err, "wechatpay error")
	}

	err = defray.InitDefray()
	if err != nil {
		return errors.Warp(err, "defray error")
	}

	err = fuwuhao.InitFuWuHao()
	if err != nil {
		return errors.Warp(err, "fuwuhao error")
	}

	err = auth.InitAuth()
	if err != nil {
		return errors.Warp(err, "auth error")
	}

	err = dbinit.CreateFooter()
	if err != nil {
		return errors.Warp(err, "rander error")
	}

	err = ip.InitYunIP()
	if err != nil {
		return errors.Warp(err, "rander error")
	}

	return nil
}

func InitSqlClear() errors.WTError {
	var err errors.WTError

	err = signalexit.InitSignalExit(1)
	if err != nil {
		return errors.Warp(err, "signal error")
	}

	signalexit.AddExitByFunc(func(ctx context.Context, _ os.Signal) context.Context {
		CloseAll()
		return context.WithValue(ctx, "InitClose", true)
	})

	err = peername.InitPeerName(config.EnvPrefix)
	if err != nil {
		return errors.Warp(err, "redis error")
	}

	err = rand.InitRander()
	if err != nil {
		return errors.Warp(err, "rander error")
	}

	err = redis.InitRedis()
	if err != nil {
		return errors.Warp(err, "redis error")
	}

	err = mysql.InitMysql()
	if err != nil {
		return errors.Warp(err, "mysql error")
	}

	err = logger.InitLogger(config.BackendConfig.SqlClear.LogServiceName)
	if err != nil {
		return errors.Warp(err, "logger error")
	}

	return nil
}
