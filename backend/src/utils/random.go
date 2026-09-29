package utils

import (
	"gitee.com/wuntsong-auth/backend/src/rand"
	errors "github.com/wuntsong-org/wterrors"
	"strconv"
)

func GenerateRandomInt(min, max int) int64 {
	return int64(rand.GlobalRander.Intn(max-min+1) + min)
}

func GenerateUniqueNumber(length int) (string, errors.WTError) {
	// 生成指定长度的随机字节序列
	randomBytes := make([]byte, length)
	_, err := rand.GlobalRander.Read(randomBytes)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	// 将随机字节序列转换为纯数字的字符串
	uniqueNumber := ""
	for _, b := range randomBytes {
		uniqueNumber += strconv.Itoa(int(b))
	}

	// 截取长度为length的纯数字字符串
	if len(uniqueNumber) > length {
		uniqueNumber = uniqueNumber[:length]
	}

	return uniqueNumber, nil
}
