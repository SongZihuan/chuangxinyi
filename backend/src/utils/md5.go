package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
