package respmsg

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
)

func InitErrorHandler() {
	httpx.SetErrorHandlerCtx(ErrorHandler)
}

func ErrorHandler(ctx context.Context, e error) (int, any) {
	if e == nil {
		return ReturnUnknown(ctx, errors.Errorf("empty error"))
	}

	err := errors.WarpQuick(e)

	switch true {
	case errors.Is(err, RedisSystemError), errors.Is(err, MySQLSystemError), errors.Is(err, SMSError), errors.Is(err, OSSError), errors.Is(err, BadContextError):
		return Return500(ctx, err)
	case errors.Is(err, JWTError), errors.Is(err, AlipayError):
		return Return(ctx, err)
	}

	return ReturnUnknown(ctx, err)
}

func Return(ctx context.Context, err errors.WTError) (int, any) {
	return http.StatusOK, types.RespEmpty{
		Resp: GetRespByError(ctx, SystemError, err, "请求错误"),
	}
}

func Return500(ctx context.Context, err errors.WTError) (int, any) {
	return http.StatusInternalServerError, types.RespEmpty{
		Resp: GetRespByError(ctx, SystemError, err, "系统错误"),
	}
}

func ReturnUnknown(ctx context.Context, err errors.WTError) (int, any) {
	return http.StatusOK, types.RespEmpty{
		Resp: GetRespByError(ctx, UnknownError, err, "请求错误"),
	}
}
