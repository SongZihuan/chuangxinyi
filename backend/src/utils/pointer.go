package utils

func GetPointer[T any](data T) *T {
	return &data
}
