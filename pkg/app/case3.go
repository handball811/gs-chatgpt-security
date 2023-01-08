//go:build case3
// +build case3

package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"io"
)

func Encrypt(data interface{}, key []byte) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// Gobエンコーダーを使用して構造体をバイト列にエンコードする
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	// 共通鍵を使用して暗号化する
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+buf.Len())
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], buf.Bytes())

	return ciphertext, nil
}

func EncryptAndGenerateKey(data interface{}) ([]byte, []byte, error) {
	// 鍵を生成する
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, nil, err
	}

	ciphertext, err := Encrypt(data, key)
	if err != nil {
		return nil, nil, err
	}

	return ciphertext, key, nil
}

func Decrypt(ciphertext []byte, key []byte, data interface{}) error {
	// 共通鍵を使用して復号する
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return errors.New("ciphertext is too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	// Gobデコーダーを使用してバイト列を構造体にデコードする
	decoder := gob.NewDecoder(bytes.NewReader(ciphertext))
	return decoder.Decode(data)
}
