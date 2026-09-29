package utils

type Closer interface {
	Close() error
}

func Close(f Closer) {
	_ = f.Close()
}
