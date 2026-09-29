package utils

func StringIn(dst string, lst []string) bool {
	for _, l := range lst {
		if l == dst {
			return true
		}
	}

	return false
}
