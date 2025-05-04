package user

import (
	userServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user"
	"github.com/mohamadrezamomeni/momo/entity"
	userValidator "github.com/mohamadrezamomeni/momo/validator/user"
)

type Handler struct {
	userSvc   UserService
	validator userValidator.Validator
}

type UserService interface {
	Create(*userServiceDto.AddUser) (*entity.User, error)
}

func New(userSvc UserService) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}
