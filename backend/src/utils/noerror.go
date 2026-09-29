package utils

import errors "github.com/wuntsong-org/wterrors"

type Logger interface {
	Error(string, ...any)
}

func MustNotError(err error) {
	if err != nil {
		panic(err)
	}
}

func Recover(logger Logger, err *errors.WTError, msg string) {
	e := recover()
	if e != nil {
		logger.Error("Error (%s): %v", msg, e)

		if err != nil {
			*err = errors.Errorf("error (%s): %v", msg, e)
		}
	}
}
