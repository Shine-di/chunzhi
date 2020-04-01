package rsa

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

func MD5Hex(val string) string {
	md5Obj := md5.New()
	md5Obj.Write([]byte(val[:]))
	return hex.EncodeToString(md5Obj.Sum(nil))
}

func Sign(val, privateKey string) (string, error) {
	signHex := md5.MD5Hex(val)
	keyByte, _ := pem.Decode([]byte(privateKey))
	if keyByte == nil {
		return "", errors.New("private key error")
	}

	key, errKey := x509.ParsePKCS8PrivateKey(keyByte.Bytes)
	if errKey != nil {
		return "", errKey
	}
	h := crypto.SHA256.New()
	h.Write([]byte(signHex))
	hash := h.Sum(nil)

	signature, errSign := rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), crypto.SHA256, hash)
	if errSign != nil {
		return "", errSign
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func VerifySign(val, sign, publicKey string) error {
	signHex := md5.MD5Hex(val)

	keyByte, _ := pem.Decode([]byte(publicKey))
	if keyByte == nil {
		return errors.New("public key error")
	}

	key, errKey := x509.ParsePKIXPublicKey(keyByte.Bytes)
	if errKey != nil {
		return errKey
	}

	h := crypto.SHA256.New()
	h.Write([]byte(signHex))
	hash := h.Sum(nil)

	signDecode, errDecode := base64.StdEncoding.DecodeString(sign)
	if errDecode != nil {
		return errDecode
	}

	return rsa.VerifyPKCS1v15(key.(*rsa.PublicKey), crypto.SHA256, hash, signDecode)
}
