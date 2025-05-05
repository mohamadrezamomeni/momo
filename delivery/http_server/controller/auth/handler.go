package auth

import (
	authServiceDto "github.com/mohamadrezamomeni/momo/dto/service/auth"
	authValidation "github.com/mohamadrezamomeni/momo/validator/auth"
)

type Handler struct {
	authSvc    AuthService
	validation *authValidation.Validation
}

type AuthService interface {
	Login(*authServiceDto.LoginDto) (string, string, error)
}

func New(authSvc AuthService, validation *authValidation.Validation) *Handler {
	return &Handler{
		authSvc:    authSvc,
		validation: validation,
	}
}
