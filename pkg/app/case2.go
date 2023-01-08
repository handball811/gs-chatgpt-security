//go:build case2
// +build case2

package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"io"
)

// Encrypt takes in a struct and a byte slice as input and returns the encrypted version of the struct as a byte slice.
func Encrypt(data interface{}, key []byte) ([]byte, error) {
	// convert the struct to a byte slice using gob encoding
	dataBytes := new(bytes.Buffer)
	enc := gob.NewEncoder(dataBytes)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}

	// create a new cipher with the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// generate a random nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// create a new gcm for encryption
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// encrypt the data
	ciphertext := aesgcm.Seal(nil, nonce, dataBytes.Bytes(), nil)

	// append the nonce to the beginning of the ciphertext
	encrypted := append(nonce, ciphertext...)

	return encrypted, nil
}

// EncryptAndGenerateKey takes in a struct as input and returns the encrypted version of the struct as a byte slice.
// It also generates a random key internally and returns it as part of the output.
func EncryptAndGenerateKey(data interface{}) ([]byte, []byte, error) {
	// generate a random key
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, nil, err
	}

	// encrypt the data using the generated key
	encrypted, err := Encrypt(data, key)
	if err != nil {
		return nil, nil, err
	}

	return encrypted, key, nil
}

// Decrypt takes in an encrypted byte slice and the key used to encrypt it, and returns the original struct.
func Decrypt(encrypted []byte, key []byte, data interface{}) error {
	// get the nonce from the beginning of the encrypted slice
	nonce := encrypted[:12]
	ciphertext := encrypted[12:]

	// create a new cipher with the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// create a new gcm for decryption
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// decrypt the data
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	// convert the decrypted byte slice back into the original struct using gob decoding
	dataBytes := bytes.NewBuffer(plaintext)
	dec := gob.NewDecoder(dataBytes)
	err = dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}
