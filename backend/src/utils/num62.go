package utils

import (
	"math/big"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeToBase62(data []byte) string {
	var encoded strings.Builder
	quotient := new(big.Int)
	remainder := new(big.Int)

	num := new(big.Int).SetBytes(data)
	base := big.NewInt(62)

	for num.Sign() > 0 {
		quotient.DivMod(num, base, remainder)
		num.Set(quotient)

		encoded.WriteByte(base62Chars[remainder.Int64()])
	}

	// Reverse the encoded string
	result := encoded.String()
	runes := []rune(result)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func DecodeFromBase62(encoded string) []byte {
	num := new(big.Int)
	base := big.NewInt(62)

	for _, char := range encoded {
		digit := strings.IndexRune(base62Chars, char)
		num.Mul(num, base)
		num.Add(num, big.NewInt(int64(digit)))
	}

	return num.Bytes()
}
