package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

func ReadBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
func RSAEncrypt(publicKeyFile string, str string) (string, error) {
	data := []byte(str)
	publicKey, err := ReadBytes(publicKeyFile)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	result, err := rsa.EncryptPKCS1v15(rand.Reader, pubInterface.(*rsa.PublicKey), data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}
func RSADecrypt(privateKeyFile string, str string) (string, error) {
	decode, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	data := []byte(decode)
	privateKey, err := ReadBytes(privateKeyFile)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	result, err := rsa.DecryptPKCS1v15(rand.Reader, priv, data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
