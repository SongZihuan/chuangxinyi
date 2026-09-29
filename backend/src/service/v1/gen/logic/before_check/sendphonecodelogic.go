package before_check

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/checkcode"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/sms"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type SendPhoneCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPhoneCodeLogic {
	return &SendPhoneCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPhoneCodeLogic) SendPhoneCode(req *types.SendPhoneCodeReq) (resp *types.RespEmpty, err error) {
	if !utils.IsPhoneNumber(req.Phone) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPhone, "错误的手机号"),
		}, nil
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	limitKey := fmt.Sprintf("sms:sendlimit:%s", req.Phone)
	res, err := redis.Get(l.ctx, limitKey).Result()
	if err == nil && res == "1" {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.SendLimit, "sms频繁发送限制", "太频繁"),
		}, nil
	}

	key := fmt.Sprintf("code:phone:%s", req.Phone)
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

	keyBusy := fmt.Sprintf("code:phone:busy:%s", req.Phone)
	resBusy := redis.Exists(l.ctx, keyBusy)
	busy, err := resBusy.Result()
	if err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	} else if busy == 1 { // 不使用runMode跳过检查，因为会导致阿里云流限
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

		go func(code int64, phone string) {
			defer utils.Recover(logger.Logger, nil, "")
			err = sms.SendImportCode(code, phone)
			if !errors.Is(err, sms.SMSSendLimit) && err != nil {
				logger.Logger.Error("Send sms error: %s, %s", err.Error(), phone)
			}
		}(code, req.Phone)
	} else {
		res2 := redis.Set(l.ctx, key, fmt.Sprintf("%06d;%s;%d", code, checkcode.NormalCode, web.ID), time.Minute*5)
		err = res2.Err()
		if err != nil {
			return nil, respmsg.RedisSystemError.WarpQuick(err)
		}

		go func(code int64, phone string) {
			defer utils.Recover(logger.Logger, nil, "")
			err = sms.SendCode(code, phone)
			if !errors.Is(err, sms.SMSSendLimit) && err != nil {
				logger.Logger.Error("Send sms error: %s, %s", err.Error(), phone)
			}
		}(code, req.Phone)
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
