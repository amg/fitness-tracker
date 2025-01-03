package env

// Copied over from:
//  https://github.com/golang-jwt/jwt/blob/main/test/helpers.go
//

import (
	"crypto"
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func LoadRSAPrivateKeyFromDisk(location string) *rsa.PrivateKey {
	keyData, e := os.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}

	return ParseRSAPrivateKeyFromPEMString(keyData)
}

func ParseRSAPrivateKeyFromPEMString(pem []byte) *rsa.PrivateKey {
	key, e := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func LoadRSAPublicKeyFromDisk(location string) *rsa.PublicKey {
	keyData, e := os.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	return ParseRSAPublicKeyFromPEMString(keyData)
}

func ParseRSAPublicKeyFromPEMString(pem []byte) *rsa.PublicKey {
	key, e := jwt.ParseRSAPublicKeyFromPEM(pem)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func LoadECPrivateKeyFromDisk(location string) crypto.PrivateKey {
	keyData, e := os.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseECPrivateKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func LoadECPublicKeyFromDisk(location string) crypto.PublicKey {
	keyData, e := os.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseECPublicKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}
