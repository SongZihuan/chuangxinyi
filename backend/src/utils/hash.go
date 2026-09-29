package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256WithBase62(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hash := hasher.Sum(nil)
	return EncodeToBase62(hash)
}

func HashSHA256(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func HashSHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
