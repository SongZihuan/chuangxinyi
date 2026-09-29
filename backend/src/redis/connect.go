package redis

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/signalexit"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

var client *redis.Client
var redcache *redis.Client
var memcache *cache.Cache

const KeepTTL = redis.KeepTTL

func InitRedis() errors.WTError {
	if len(config.BackendConfig.Redis.Addr) == 0 {
		return errors.Errorf("redis addr must be given")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     config.BackendConfig.Redis.Addr,
		Username: config.BackendConfig.Redis.UserName,
		Password: config.BackendConfig.Redis.Password,
		DB:       int(config.BackendConfig.Redis.DB),
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		_ = client.Close()
		client = nil
		return errors.Warp(err, "redis is fail to connect")
	}

	if len(config.BackendConfig.Cache.Addr) == 0 {
		redcache = nil
	} else {
		redcache = redis.NewClient(&redis.Options{
			Addr:     config.BackendConfig.Cache.Addr,
			Username: config.BackendConfig.Cache.UserName,
			Password: config.BackendConfig.Redis.Password,
			DB:       int(config.BackendConfig.Cache.DB),
		})

		err = redcache.Ping(context.Background()).Err()
		if err != nil {
			_ = client.Close()
			client = nil

			_ = redcache.Close()
			redcache = nil

			return errors.Warp(err, "cache is fail to connect")
		}
	}

	memcache = cache.New(time.Minute*30, time.Minute*60)

	signalexit.AddExitByFunc(lockExitFunc)

	return nil
}

func CloseRedis() {
	if client != nil {
		_ = client.Close()
		client = nil
	}

	if redcache != nil {
		_ = redcache.Close()
		redcache = nil
	}
}
