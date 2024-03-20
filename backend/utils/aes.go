package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func EncryptAES2Base64(key, text []byte) (string, error) {
	encrypt, err := EncryptAES(key, text)
	if err == nil {
		return base64.StdEncoding.EncodeToString(encrypt), err
	}
	return "", err
}
func DecryptAES2Base64(key []byte, text string) ([]byte, error) {
	encrypt, err := base64.StdEncoding.DecodeString(text)
	if err == nil {
		return DecryptAES(key, encrypt)

	}
	return nil, err
}
func EncryptAES(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	msg := pad(text)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], msg)

	return ciphertext, nil
}

func DecryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return unPad(ciphertext), nil
}

// generateKey 生成一个 AES 密钥。
func generateKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256 需要 32 字节长的密钥
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// pad 对数据进行填充，使其长度为块大小的倍数。
func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// unPad 移除填充数据。
func unPad(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}
