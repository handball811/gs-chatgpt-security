//go:build case1
// +build case1

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

// Encrypt takes in a struct and a key as []byte and returns the encrypted []byte
func Encrypt(data interface{}, key []byte) ([]byte, error) {
	// convert data to []byte
	b, err := encode(data)
	if err != nil {
		return nil, errors.New("failed to encode data: " + err.Error())
	}

	// generate nonce
	nonce, err := generateNonce()
	if err != nil {
		return nil, errors.New("failed to generate nonce: " + err.Error())
	}

	// create block cipher with key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create block cipher: " + err.Error())
	}

	// create cipher stream using block cipher and nonce
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create cipher stream: " + err.Error())
	}

	// encrypt data using cipher stream
	ciphertext := aesgcm.Seal(nil, nonce, b, nil)

	// return encrypted data with nonce appended
	return append(nonce, ciphertext...), nil
}

// EncryptAndGenerateKey takes in a struct and returns the encrypted []byte and the generated key as []byte
func EncryptAndGenerateKey(data interface{}) ([]byte, []byte, error) {
	// generate key
	key, err := generateKey()
	if err != nil {
		return nil, nil, errors.New("failed to generate key: " + err.Error())
	}

	// encrypt data using key
	encrypted, err := Encrypt(data, key)
	if err != nil {
		return nil, nil, errors.New("failed to encrypt data: " + err.Error())
	}

	return encrypted, key, nil
}

// Decrypt takes in an encrypted []byte and a key as []byte and returns the decrypted struct
func Decrypt(encrypted []byte, key []byte, data interface{}) error {
	// extract nonce from encrypted data
	nonce, ciphertext := encrypted[:12], encrypted[12:]

	// create block cipher with key
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("failed to create block cipher: " + err.Error())
	}

	// create cipher stream using block cipher and nonce
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return errors.New("failed to create cipher stream: " + err.Error())
	}

	// decrypt ciphertext using cipher stream
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("failed to decrypt ciphertext: " + err.Error())
	}

	// decode plaintext into data struct
	err = decode(plaintext, data)
	if err != nil {
		return errors.New("failed to decode data: " + err.Error())
	}

	return nil
}

func generateNonce() ([]byte, error) {
	nonce := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, errors.New("failed to generate nonce: " + err.Error())
	}
	return nonce, nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, errors.New("failed to generate key: " + err.Error())
	}
	return key, nil
}

func encode(data interface{}) ([]byte, error) {
	// use gob package to encode data into []byte
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(data)
	if err != nil {
		return nil, errors.New("failed to encode data: " + err.Error())
	}
	return b.Bytes(), nil
}

func decode(b []byte, data interface{}) error {
	// use gob package to decode []byte into data
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(data)
	if err != nil {
		return errors.New("failed to decode data: " + err.Error())
	}
	return nil
}
