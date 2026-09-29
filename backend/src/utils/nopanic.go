package utils

type BaseFunc func()

func RecoverFunc(f BaseFunc, l Logger, msg string) BaseFunc {
	return func() {
		defer Recover(l, nil, msg)
		f() // 正常调用
	}
}
