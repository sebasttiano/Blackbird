package common

import (
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestUnmarshalRSAPrivate(t *testing.T) {
	privateKey, err := os.ReadFile("../../files/rsa_private")
	if err != nil {
		t.Errorf("failed to read crypto key, %v", err)
	}
	RSAPrivateType := UnmarshalRSAPrivate(privateKey)
	t.Run("test unmarshall private", func(t *testing.T) {
		assert.Equal(t, RSAPrivateType.Size(), 512)
		assert.Equal(t, RSAPrivateType.E, 65537)
	})
	nillKey := UnmarshalRSAPrivate([]byte("1, 2, 3"))
	t.Run("test unmarshall nil private", func(t *testing.T) {
		var k *rsa.PrivateKey = nil
		assert.Equal(t, nillKey, k)
	})
}

func TestUnmarshalRSAPublic(t *testing.T) {
	publicKey, err := os.ReadFile("../../files/rsa_public")
	if err != nil {
		t.Errorf("failed to read crypto key, %v", err)
	}
	RSAPublicType := UnmarshalRSAPublic(publicKey)
	t.Run("test unmarshall public", func(t *testing.T) {
		assert.Equal(t, RSAPublicType.Size(), 512)
		assert.Equal(t, RSAPublicType.E, 65537)
	})
	nillKey := UnmarshalRSAPublic([]byte("9, 13, 33"))
	t.Run("test unmarshall nil public", func(t *testing.T) {
		var k *rsa.PublicKey = nil
		assert.Equal(t, nillKey, k)
	})
}

func TestEncryptRSA(t *testing.T) {
	testTable := []struct {
		name string
		msg  string
		//expected string
	}{
		{name: "OK", msg: "hello world"},
		{name: "OK2", msg: "FOO BAR"},
		{name: "OK3", msg: "фыыфввфы"},
	}

	for _, tt := range testTable {
		publicKey, err := os.ReadFile("../../files/rsa_public")
		if err != nil {
			t.Errorf("failed to read crypto key, %v", err)
		}
		RSAPublicType := UnmarshalRSAPublic(publicKey)
		privateKey, err := os.ReadFile("../../files/rsa_private")
		if err != nil {
			t.Errorf("failed to read crypto key, %v", err)
		}
		RSAPrivateType := UnmarshalRSAPrivate(privateKey)
		t.Run(tt.name, func(t *testing.T) {
			encMsg, _ := EncryptRSA(tt.msg, RSAPublicType)
			origMsg, _ := DecryptRSA(encMsg, RSAPrivateType)
			assert.Equal(t, origMsg, tt.msg)
		})
	}
}
