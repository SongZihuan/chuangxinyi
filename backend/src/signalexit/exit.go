package signalexit

import (
	"context"
	errors "github.com/wuntsong-org/wterrors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ExitFunc func(context.Context, os.Signal) context.Context
type DeferFunc func()

var ExitFuncList = make([]*ExitFunc, 0, 10)
var ExitMutex sync.Mutex

func InitSignalExit(exitCode int64) errors.WTError {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	go func() {
		s := <-sig

		ExitMutex.TryLock() // 不用管实际是否lock
		for i := len(ExitFuncList) - 1; i >= 0; i-- {
			f := ExitFuncList[i]
			if f == nil {
				continue
			}

			func() {
				defer func() {
					_ = recover()
				}()
				ctx = (*f)(ctx, s)
			}()
		}

		os.Exit(int(exitCode))
	}()

	return nil
}

func AddExitFunc(f *ExitFunc) {
	ExitMutex.Lock()
	defer ExitMutex.Unlock()

	ExitFuncList = append(ExitFuncList, f)
}

func AddExitByFunc(f ExitFunc) {
	ExitMutex.Lock()
	defer ExitMutex.Unlock()

	ExitFuncList = append(ExitFuncList, &f)
}

func AddExitFuncAsDefer(f DeferFunc) DeferFunc {
	var exitFunc ExitFunc = func(ctx context.Context, _ os.Signal) context.Context {
		f()
		return ctx
	}

	AddExitFunc(&exitFunc)

	return func() {
		DeleteExitFunc(&exitFunc)
		f()
	}
}

func DeleteExitFunc(f *ExitFunc) {
	ExitMutex.Lock()
	defer ExitMutex.Unlock()

	for i, v := range ExitFuncList {
		if v == f {
			ExitFuncList = append(append(make([]*ExitFunc, 0, len(ExitFuncList)-1), ExitFuncList[:i]...), ExitFuncList[i+1:]...)
		}
	}
}
