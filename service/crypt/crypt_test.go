package crypt

import (
	"testing"
)

func TestHashing(t *testing.T) {
	crypt := New(&CryptConfig{
		EncryptKey: "qFQGz0yE0nYBxqi9",
	})
	password := "zoyxFGPRgL"

	passwordHashed, err := crypt.Hash(password)
	if err != nil || password == passwordHashed {
		t.Error("error to hash password")
	}

	isSame := crypt.CheckHash(password, passwordHashed)

	if !isSame {
		t.Error("error to compare password")
	}

	wrongPassword := "qeZgIRR8jV"

	isSame = crypt.CheckHash(wrongPassword, passwordHashed)
	if isSame {
		t.Error("the compare must be wrong")
	}
}

func TestEncryptin(t *testing.T) {
	crypt := New(&CryptConfig{
		EncryptKey: "qFQGz0yE0nYBxqi9",
	})
	plainText := "hello this is a test for encrypting data :)"

	encryptedData, err := crypt.Encrypt(plainText)
	if err != nil || encryptedData == plainText {
		t.Error("error to encrypt data")
	}

	originalText, err := crypt.Decrypt(encryptedData)

	if originalText != plainText {
		t.Error("error to encryption processing")
	}
}
