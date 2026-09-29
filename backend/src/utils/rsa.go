package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	errors "github.com/wuntsong-org/wterrors"
)

func SignRsaHash256Sign(data string, privateKey *rsa.PrivateKey) ([]byte, errors.WTError) {
	hashed := sha256.Sum256([]byte(data))
	res, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return res, nil
}

func VerifyRsaHash256Sign(data string, signature []byte, publicKey *rsa.PublicKey) errors.WTError {
	hashed := sha256.Sum256([]byte(data))
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func ReadRsaPubKeyFromCert(c []byte) (*rsa.PublicKey, errors.WTError) {
	block, _ := pem.Decode(c)
	if block == nil {
		return nil, errors.Errorf("bad cert")
	}

	if block.Type != "CERTIFICATE" {
		return nil, errors.Errorf("bad cert")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Errorf("bad cert")
	}

	pubkey, ok := cert.PublicKey.(rsa.PublicKey)
	if !ok {
		return nil, errors.Errorf("bad cert")
	}

	return &pubkey, nil
}

func ReadRsaPublicKey(c []byte) (*rsa.PublicKey, errors.WTError) {
	block, _ := pem.Decode(c)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.Errorf("bad public key type")
	}

	// 解析RSA公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Errorf("bad public key: " + err.Error())
	}

	rsaPubKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.Errorf("bad public key")
	}

	return rsaPubKey, nil
}

func ReadRsaPrivateKey(c []byte) (*rsa.PrivateKey, errors.WTError) {
	block, _ := pem.Decode(c)
	if block == nil || (block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY") {
		return nil, errors.Errorf("bad private key")
	}

	privateKeyRsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privateKeyRsa, nil
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	privateKeyRsa, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.Errorf("not a rsa private key")
	}

	return privateKeyRsa, nil
}
