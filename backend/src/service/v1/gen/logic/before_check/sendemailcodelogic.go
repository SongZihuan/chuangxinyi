package before_check

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/email"
	"gitee.com/wuntsong-auth/backend/src/global/checkcode"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type SendEmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailCodeLogic {
	return &SendEmailCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendEmailCodeLogic) SendEmailCode(req *types.SendEmailCodeReq) (resp *types.RespEmpty, err error) {
	if !utils.IsEmailAddress(req.Email) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadEmail, "错误的邮箱地址"),
		}, nil
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	key := fmt.Sprintf("code:email:%s", req.Email)
	res1 := redis.TTL(l.ctx, key)
	t, err := res1.Result()
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	if config.BackendConfig.GetMode() != config.RunModeDevelop && t >= time.Minute*4 {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.TooBusy, "上次申请的验证码还未过期", "太频繁"),
		}, nil
	}

	keyBusy := fmt.Sprintf("code:email:busy:%s", req.Email)
	resBusy := redis.Exists(l.ctx, keyBusy)
	busy, err := resBusy.Result()
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	} else if config.BackendConfig.GetMode() != config.RunModeDevelop && busy == 1 {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.TooBusy, "请求过于频繁", "太频繁"),
		}, nil
	}

	code := utils.GenerateRandomInt(0, 999999)

	if web.ID == warp.UserCenterWebsite {
		res2 := redis.Set(l.ctx, key, fmt.Sprintf("%06d;%s;%d", code, checkcode.ImportCode, web.ID), time.Minute*5)
		err = res2.Err()
		if err != nil {
			return nil, respmsg.RedisSystemError.WarpQuick(err)
		}

		go func(code int64, emailAddress string) {
			defer utils.Recover(logger.Logger, nil, "")
			err := email.SendImportCode(code, emailAddress)
			if err != nil {
				logger.Logger.Error("Send email error: %s, %s", err.Error(), emailAddress)
			}
		}(code, req.Email)
	} else {
		res2 := redis.Set(l.ctx, key, fmt.Sprintf("%06d;%s;%d", code, checkcode.NormalCode, web.ID), time.Minute*5)
		err = res2.Err()
		if err != nil {
			return nil, respmsg.RedisSystemError.WarpQuick(err)
		}

		go func(code int64, emailAddress string) {
			defer utils.Recover(logger.Logger, nil, "")
			err := email.SendCode(code, emailAddress)
			if err != nil {
				logger.Logger.Error("Send email error: %s, %s", err.Error(), emailAddress)
			}
		}(code, req.Email)
	}

	res3 := redis.Set(l.ctx, keyBusy, "1", time.Minute*1)
	err = res3.Err()
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccessWithDebug(l.ctx, "获取的验证码是：%d", code),
	}, nil
}
