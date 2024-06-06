package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// Encrypt encrypts a string using a key.
func Encrypt(toEncrypt, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	plain := []byte(toEncrypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plain))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plain)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a string using a key.
func Decrypt(toDecrypt, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	cipherText, err := base64.URLEncoding.DecodeString(toDecrypt)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", shared.ErrCipherTooShort
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}
