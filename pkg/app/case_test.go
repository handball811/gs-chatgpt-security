package app

import (
	"testing"
)

type TestStruct struct {
	Field1 string
	Field2 int
	Field3 bool
}

func TestEncrypt(t *testing.T) {
	data := TestStruct{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}
	key := []byte("this is a test k")

	encrypted, err := Encrypt(data, key)
	if err != nil {
		t.Error("failed to encrypt data: ", err)
	}
	if len(encrypted) == 0 {
		t.Error("encrypted data is empty")
	}

	var decryptedData TestStruct
	err = Decrypt(encrypted, key, &decryptedData)
	if err != nil {
		t.Error("failed to decrypt data: ", err)
	}

	if data != decryptedData {
		t.Error("decrypted data does not match original data")
	}
}

func TestEncryptAndGenerateKey(t *testing.T) {
	data := TestStruct{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}

	encrypted, key, err := EncryptAndGenerateKey(data)
	if err != nil {
		t.Error("failed to encrypt and generate key: ", err)
	}
	if len(encrypted) == 0 {
		t.Error("encrypted data is empty")
	}
	if len(key) == 0 {
		t.Error("generated key is empty")
	}

	var decryptedData TestStruct
	err = Decrypt(encrypted, key, &decryptedData)
	if err != nil {
		t.Error("failed to decrypt data: ", err)
	}

	if data != decryptedData {
		t.Error("decrypted data does not match original data")
	}
}

func TestDecrypt(t *testing.T) {
	data := TestStruct{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}
	key := []byte("this is a test k")

	encrypted, err := Encrypt(data, key)
	if err != nil {
		t.Error("failed to encrypt data: ", err)
	}

	var decryptedData TestStruct
	err = Decrypt(encrypted, key, &decryptedData)
	if err != nil {
		t.Error("failed to decrypt data: ", err)
	}

	if data != decryptedData {
		t.Error("decrypted data does not match original data")
	}
}
