package utils

func HasDuplicate[T string | int64](strings []T) bool {
	seen := make(map[T]bool)
	for _, str := range strings {
		if seen[str] {
			return true
		}
		seen[str] = true
	}
	return false
}
