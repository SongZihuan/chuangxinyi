package utils

import (
	"crypto/aes"
	"crypto/cipher"
	errors "github.com/wuntsong-org/wterrors"
)

func DecryptAEAD(data []byte, key, nonce, additionalData string) ([]byte, errors.WTError) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	plaintext, err := aead.Open(nil, []byte(nonce), data, []byte(additionalData))
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return plaintext, nil
}
