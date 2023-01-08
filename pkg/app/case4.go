//go:build case4
// +build case4

package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
)

func Encrypt(data interface{}, key []byte) ([]byte, error) {
	// JSONエンコード
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 暗号化
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// base64エンコード
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func EncryptAndGenerateKey(data interface{}) ([]byte, []byte, error) {
	// 暗号化キーの生成
	key := make([]byte, 32) // 256 bits
	if _, err := rand.Read(key); err != nil {
		return nil, nil, err
	}

	encrypted, err := Encrypt(data, key)
	if err != nil {
		return nil, nil, err
	}

	return encrypted, key, nil
}

func Decrypt(encrypted []byte, key []byte, data interface{}) error {
	// base64デコード
	ciphertext, err := base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		return err
	}

	// 復号化
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	if len(ciphertext) < aes.BlockSize {
		return err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	// JSONデコード
	return json.Unmarshal(ciphertext, data)
}
