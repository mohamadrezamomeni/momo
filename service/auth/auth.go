package auth

import (
	authServiceDto "github.com/mohamadrezamomeni/momo/dto/service/auth"
	userServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user"
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
	Create(*userServiceDto.AddUser) (*entity.User, error)
}

func New(userSvc UserService, crypt *crypt.Crypt, config *AuthConfig) *Auth {
	return &Auth{
		userSvc: userSvc,
		crypt:   crypt,
		config:  config,
	}
}

func (a *Auth) Login(inpt *authServiceDto.LoginDto) (string, string, error) {
	user, err := a.userSvc.FindByUsername(inpt.Username)
	if err != nil {
		return "", "", err
	}

	token, err := a.createToken(user)
	if err != nil {
		return "", "", err
	}

	return token, "", nil
}
