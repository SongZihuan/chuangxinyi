package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	errors "github.com/wuntsong-org/wterrors"
	"os"
)

func GetCertificateSerialNumber(certificate x509.Certificate) string {
	return fmt.Sprintf("%X", certificate.SerialNumber.Bytes())
}

func LoadCertificateWithPath(path string) (*x509.Certificate, errors.WTError) {
	certificateBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf("read certificate pem file err:%s", err.Error())
	}
	return LoadCertificate(string(certificateBytes))
}

func LoadCertificate(certificateStr string) (*x509.Certificate, errors.WTError) {
	block, _ := pem.Decode([]byte(certificateStr))
	if block == nil {
		return nil, errors.Errorf("decode certificate err")
	}
	if block.Type != "CERTIFICATE" {
		return nil, errors.Errorf("the kind of PEM should be CERTIFICATE")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Errorf("parse certificate err:%s", err.Error())
	}
	return certificate, nil
}
