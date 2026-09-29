package init

import (
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
)

func CloseAll() {
	mysql.CloseMysql()
	redis.CloseRedis()
}
