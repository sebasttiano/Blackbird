package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// MarshalRSAPrivate
func MarshalRSAPrivate(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
}

// UnmarshalRSAPrivate bytes to private key
func UnmarshalRSAPrivate(data []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}
	return priv
}

// MarshalRSAPublic public key to bytes
func MarshalRSAPublic(pub *rsa.PublicKey) []byte {

	pubPKCS := x509.MarshalPKCS1PublicKey(pub)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubPKCS,
	})

	return pubBytes
}

// UnmarshalRSAPublic bytes to public key type
func UnmarshalRSAPublic(data []byte) *rsa.PublicKey {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil
	}
	return pub
}

// GenerateKey makes private and public RSA keys
func GenerateKey(size int) ([]byte, []byte, error) {
	reader := rand.Reader
	bitSize := size

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

	pubKeyStr := MarshalRSAPublic(&key.PublicKey)
	privKeyStr := MarshalRSAPrivate(key)

	return pubKeyStr, privKeyStr, nil
}

// EncryptRSA encrypts rsa message
func EncryptRSA(msg string, pub *rsa.PublicKey) (string, error) {

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub,
		[]byte(msg),
		nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// DecryptRSA decrypts rsa message
func DecryptRSA(data string, priv *rsa.PrivateKey) (string, error) {

	data2, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	decrypted, err := rsa.DecryptOAEP(sha256.New(),
		rand.Reader, priv, data2, nil)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
