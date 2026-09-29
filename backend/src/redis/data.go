package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return client.Set(ctx, key, value, expiration)
}

func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return client.SetNX(ctx, key, value, expiration)
}

func Get(ctx context.Context, key string) *redis.StringCmd {
	return client.Get(ctx, key)
}

func Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return client.Keys(ctx, pattern)
}

func Exists(ctx context.Context, key string) *redis.IntCmd {
	return client.Exists(ctx, key)
}

func Del(ctx context.Context, key ...string) *redis.IntCmd {
	return client.Del(ctx, key...)
}

func TTL(ctx context.Context, key string) *redis.DurationCmd {
	return client.TTL(ctx, key)
}

func Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return client.Expire(ctx, key, expiration)
}

func ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return client.ExpireAt(ctx, key, tm)
}

func Incr(ctx context.Context, key string) *redis.IntCmd {
	return client.Incr(ctx, key)
}
