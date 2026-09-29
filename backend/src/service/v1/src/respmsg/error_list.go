package respmsg

import (
	errors "github.com/wuntsong-org/wterrors"
)

var MySQLSystemError = errors.NewClass("mysql_error")
var RedisSystemError = errors.NewClass("redis_error")
var SMSError = errors.NewClass("sms_error")
var JWTError = errors.NewClass("jwt_error")
var OSSError = errors.NewClass("oss_error")
var BadContextError = errors.NewClass("bad_context_error")
var AlipayError = errors.NewClass("alipay_error")
