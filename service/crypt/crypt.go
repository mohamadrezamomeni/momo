package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"golang.org/x/crypto/bcrypt"
)

type Crypt struct {
	config *CryptConfig
}

func New(cfg *CryptConfig) *Crypt {
	return &Crypt{
		config: cfg,
	}
}

func (c *Crypt) Hash(input string) (string, error) {
	scope := "crypt.hash"

	bytes, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("the inpt is %s", input)
	}
	return string(bytes), nil
}

func (c *Crypt) CheckHash(input string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}

func (c *Crypt) Encrypt(plaintext string) (string, error) {
	scope := "crypt.encrypt"

	key := []byte(c.config.EncryptKey)

	plaintextBytes := []byte(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("err to newCypher the input is %s", plaintext)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("err to newGCM the input is %s", plaintext)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("err to readFull the input is %s", plaintext)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintextBytes, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Crypt) Decrypt(ciphertextBase64 string) (string, error) {
	scope := "crypt.decrypt"

	key := []byte(c.config.EncryptKey)
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("error to decodeString the input is %s", ciphertextBase64)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("error to newCipher the input is %s", ciphertextBase64)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("error to newGCM the input is %s", ciphertextBase64)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", momoError.Scope(scope).Errorf("length of ciphertext is bigger than nonceSize the input is %s", ciphertextBase64)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	bytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("error to aesGCM open the input is %s", ciphertextBase64)
	}
	return string(bytes), nil
}
