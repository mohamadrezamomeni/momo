package auth

import (
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/crypt"
)

type Auth struct {
	userSvc UserService
	crypt   *crypt.Crypt
	config  *AuthConfig
}

type UserService interface {
	FindByID(string) (*entity.User, error)
	FindByUsername(string) (*entity.User, error)
}

func New(userSvc UserService, crypt *crypt.Crypt, config *AuthConfig) *Auth {
	return &Auth{
		userSvc: userSvc,
		crypt:   crypt,
		config:  config,
	}
}
