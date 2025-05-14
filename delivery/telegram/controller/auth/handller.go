package auth

import authServiceDto "github.com/mohamadrezamomeni/momo/dto/service/auth"

type Handler struct {
	authSvc AuthService
}

type AuthService interface {
	Register(*authServiceDto.RegisterDto) (string, error)
}

func New(userSvc AuthService) *Handler {
	return &Handler{
		authSvc: userSvc,
	}
}
