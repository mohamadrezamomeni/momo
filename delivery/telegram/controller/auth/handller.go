package auth

import (
	authServiceDto "github.com/mohamadrezamomeni/momo/dto/service/auth"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Handler struct {
	authSvc AuthService
	userSvc UserService
}

type AuthService interface {
	Register(*authServiceDto.RegisterDto) (string, error)
}

type UserService interface {
	FindByUsername(string) (*entity.User, error)
	FindByTelegramID(string) (*entity.User, error)
}

func New(authSvc AuthService, userSvc UserService) *Handler {
	return &Handler{
		authSvc: authSvc,
		userSvc: userSvc,
	}
}
