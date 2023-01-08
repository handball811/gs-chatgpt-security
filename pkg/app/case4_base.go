//go:build case4_base
// +build case4_base

package app

import (
	"encoding/json"
	"math/rand"
)

func Encrypt(data interface{}, key []byte) ([]byte, error) {
	encrypted, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

func EncryptAndGenerateKey(data interface{}) ([]byte, []byte, error) {
	key := make([]byte, 8)
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
	if err := json.Unmarshal(encrypted, data); err != nil {
		return err
	}
	return nil
}
