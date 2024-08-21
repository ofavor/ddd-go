package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/ofavor/ddd-go/pkg/log"
)

// AES encrypt a string
func AesEncrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("[util-aes] encrypt error: ", err)
		return
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Error("[util-aes] encrypt error: ", err)
		return
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

// AES decrypt a string
func AesDecrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		log.Error("[util-aes] decrypt error: ", err)
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("[util-aes] decrypt error: ", err)
		return
	}
	if len(cipherText) < aes.BlockSize {
		err = errors.New("cliphertext block size is too short")
		log.Error("[util-aes] decrypt error: ", err)
		return
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	decodedmess = string(cipherText)
	return
}
