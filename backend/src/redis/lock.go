package redis

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
)

var LockMap = make(map[string]bool, 100)
var LockMapMutex sync.Mutex

func lockExitFunc(ctx context.Context, _ os.Signal) context.Context {
	LockMapMutex.TryLock()
	// 不需要释放锁

	for key, exists := range LockMap {
		if !exists {
			continue
		}
		releaseLock(key)
	}

	return context.WithValue(ctx, "Not-Redis-Lock", true)
}

func AcquireLockMore(ctx context.Context, lockKey string, lockTTL time.Duration) bool {
	notRedisLock, ok := ctx.Value("Not-Redis-Lock").(bool)
	if ok && !notRedisLock {
		return false
	}

	count := 10
	for count >= 0 {
		count -= 1
		res := AcquireLock(ctx, lockKey, lockTTL)
		if res {
			return true
		}
		time.Sleep(time.Millisecond * 300)
	}

	return false
}

func AcquireLock(ctx context.Context, lockKey string, lockTTL time.Duration) bool {
	notRedisLock, ok := ctx.Value("Not-Redis-Lock").(bool)
	if ok && !notRedisLock {
		return false
	}

	LockMapMutex.Lock()
	defer LockMapMutex.Unlock()

	return acquireLock(ctx, lockKey, lockTTL)
}

func ExtendLock(ctx context.Context, lockKey string, lockTTL time.Duration) bool {
	notRedisLock, ok := ctx.Value("Not-Redis-Lock").(bool)
	if ok && !notRedisLock {
		return false
	}

	LockMapMutex.Lock()
	defer LockMapMutex.Unlock()

	return extendLock(ctx, lockKey, lockTTL)
}

func ReleaseLock(lockKey string) {
	go func() {
		LockMapMutex.Lock()
		defer LockMapMutex.Unlock()

		releaseLock(lockKey)
	}()
}

func acquireLock(ctx context.Context, lockKey string, lockTTL time.Duration) bool {
	exists, ok := LockMap[lockKey]
	if ok && exists {
		return false
	}

	lockValue := generateLockValue()

	// 使用SET命令尝试获取锁，并设置过期时间，使用NX选项确保只在键不存在时设置
	result, err := client.SetNX(ctx, fmt.Sprintf("lock:%s", lockKey), lockValue, lockTTL).Result()
	if err != nil {
		logger.Logger.Error("redis error: %s", err.Error())
		return false
	}

	LockMap[lockKey] = true

	return result
}

func extendLock(ctx context.Context, lockKey string, lockTTL time.Duration) bool {
	exists, ok := LockMap[lockKey]
	if !ok || !exists {
		return false
	}

	result, err := client.Expire(ctx, fmt.Sprintf("lock:%s", lockKey), lockTTL).Result()
	if err != nil {
		logger.Logger.Error("redis error: %s", err.Error())
		return false
	}

	if !result {
		return acquireLock(ctx, lockKey, lockTTL)
	}

	return true
}

func releaseLock(lockKey string) {
	exists, ok := LockMap[lockKey]
	if !ok || !exists {
		return
	}

	delete(LockMap, lockKey)

	err := client.Del(context.Background(), fmt.Sprintf("lock:%s", lockKey)).Err() // 不用ctx
	if err != nil {
		logger.Logger.Error("redis error: %s", err.Error())
	}
}

func GenerateUUIDMore(ctx context.Context, lockPrefix string, lockTTL time.Duration, checker func(context.Context, uuid.UUID) bool) (uuid.UUID, bool) {
	count := 10
	for count >= 0 {
		count -= 1
		res, success := GenerateUUID(ctx, lockPrefix, lockTTL, checker)
		if success {
			return res, true
		}
	}

	return uuid.UUID{}, false
}

func GenerateUUID(ctx context.Context, lockPrefix string, lockTTL time.Duration, checker func(context.Context, uuid.UUID) bool) (uuid.UUID, bool) {
	uuidByte, err := uuid.NewRandom()
	if err != nil {
		return uuid.UUID{}, false
	}

	uuidString := uuidByte.String()

	if !checker(ctx, uuidByte) {
		return uuid.UUID{}, false
	}

	lockRes := AcquireLock(ctx, fmt.Sprintf("%s:%s", lockPrefix, uuidString), lockTTL)
	if !lockRes {
		return uuid.UUID{}, false
	}

	return uuidByte, true
}

func GenerateStringMore(ctx context.Context, lockPrefix string, lockTTL time.Duration, checker func(context.Context, string) bool, gen func(context.Context, uuid.UUID) string) (string, bool) {
	count := 10
	for count >= 0 {
		count -= 1
		res, success := GenerateString(ctx, lockPrefix, lockTTL, checker, gen)
		if success {
			return res, true
		}
	}

	return "", false
}

func GenerateString(ctx context.Context, lockPrefix string, lockTTL time.Duration, checker func(context.Context, string) bool, gen func(context.Context, uuid.UUID) string) (string, bool) {
	uuidByte, err := uuid.NewRandom()
	if err != nil {
		return "", false
	}

	uuidString := gen(ctx, uuidByte)
	if len(uuidString) == 0 {
		return "", false
	}

	if !checker(ctx, uuidString) {
		return "", false
	}

	lockRes := AcquireLock(ctx, fmt.Sprintf("%s:%s", lockPrefix, uuidString), lockTTL)
	if !lockRes {
		return "", false
	}

	return uuidString, true
}

func generateLockValue() string {
	return "lock_value_" + time.Now().Format(time.RFC3339Nano)
}
