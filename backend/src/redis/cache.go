package redis

import (
	"context"
	"fmt"
	"time"
)

func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	memcache.Set(key, fmt.Sprintf("%v", value), expiration)

	if redcache != nil {
		_ = redcache.Set(ctx, key, fmt.Sprintf("%v", value), expiration)
	}
}

func GetCache(ctx context.Context, key string) (string, bool) {
	res1, ok := memcache.Get(key)
	if ok {
		return fmt.Sprintf("%v", res1), true
	}

	if redcache != nil {
		res2, err := redcache.Get(ctx, key).Result()
		if err != nil || len(res2) == 0 {
			return "", false
		}

		return res2, true
	}

	return "", false
}

func DelCache(ctx context.Context, key string) {
	memcache.Delete(key)

	if redcache != nil {
		_ = redcache.Del(ctx, key)
	}
}
