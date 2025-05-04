package crypt

import (
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"golang.org/x/crypto/bcrypt"
)

type Crypt struct{}

func New() *Crypt {
	return &Crypt{}
}

func (h *Crypt) Hash(input string) (string, error) {
	scope := "crypt.hash"

	bytes, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("the inpt is %s", input)
	}
	return string(bytes), nil
}

func (h *Crypt) CheckHash(input string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}
