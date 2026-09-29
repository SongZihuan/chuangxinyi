package password

import (
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
)

func GetPasswordHash(password string, userID string) string {
	return GetPasswordSecondHash(GetPasswordFirstHash(password), userID)
}

func GetPasswordSecondHash(passwordHash string, userID string) string {
	return utils.HashSHA256(fmt.Sprintf("%s:%s:%s", config.BackendConfig.Password.Salt, passwordHash, userID))
}

func GetPasswordFirstHash(password string) string {
	return utils.HashSHA256(fmt.Sprintf("%s:%s", config.BackendConfig.Password.FrontSalt, password))
}
