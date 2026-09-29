package utils

import (
	"github.com/xlzd/gotp"
	"time"
)

func CheckTOTP(key string, code string) bool {
	totp := gotp.NewDefaultTOTP(key)

	currentTime := time.Now().Unix()
	beforeTime := time.Now().Add(-30 * time.Second).Unix()
	afterTime := time.Now().Add(30 * time.Second).Unix()

	return totp.Verify(code, currentTime) || totp.Verify(code, beforeTime) || totp.Verify(code, afterTime)
}

func GenerateTotpURL(key string, accountName string, issuerName string) string {
	totp := gotp.NewDefaultTOTP(key)

	if len(issuerName) == 0 {
		issuerName = "unknown"
	}

	if len(accountName) == 0 {
		accountName = "anonymous"
	}

	return totp.ProvisioningUri(accountName, issuerName)
}

func GetRandomSecret() string {
	return gotp.RandomSecret(64)
}
