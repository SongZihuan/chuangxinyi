package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"gitee.com/wuntsong-auth/backend/src/rand"
	errors "github.com/wuntsong-org/wterrors"
)

func pkcs7Padding(data []byte, blockSize int) []byte {
	// 判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	// 补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, errors.WTError) {
	length := len(data)
	if length == 0 {
		return nil, errors.Errorf("key error")
	}
	// 获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func AesEncrypt(data []byte, key []byte) (string, errors.WTError) {
	// 创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.WarpQuick(err)
	}
	// 判断加密快的大小
	blockSize := block.BlockSize()
	// 填充
	encryptBytes := pkcs7Padding(data, blockSize)
	// 初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	// 使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return hex.EncodeToString(crypted), nil
}

func AesDecrypt(dataString string, key []byte) ([]byte, errors.WTError) {
	data, err := hex.DecodeString(dataString)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	// 创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	// 获取块的大小
	blockSize := block.BlockSize()
	// 使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// 初始化解密数据接收切片
	crypted := make([]byte, len(data))
	// 执行解密
	blockMode.CryptBlocks(crypted, data)
	// 去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	return crypted, nil
}

func AesEncryptBase62(data []byte, key []byte) (string, errors.WTError) {
	// 创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.WarpQuick(err)
	}
	// 判断加密快的大小
	blockSize := block.BlockSize()
	// 填充
	encryptBytes := pkcs7Padding(data, blockSize)
	// 初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	// 使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return EncodeToBase62(crypted), nil
}

func AesDecryptBase62(dataString string, key []byte) ([]byte, errors.WTError) {
	data := DecodeFromBase62(dataString)

	// 创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	// 获取块的大小
	blockSize := block.BlockSize()
	// 使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// 初始化解密数据接收切片
	crypted := make([]byte, len(data))
	// 执行解密
	blockMode.CryptBlocks(crypted, data)
	// 去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	return crypted, nil
}

const AesKeySize = 32

func GenerateAESKey(keySize int) ([]byte, errors.WTError) {
	key := make([]byte, keySize)
	_, err := rand.GlobalRander.Read(key)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	return key, nil
}
