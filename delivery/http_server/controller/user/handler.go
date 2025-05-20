package user

import (
	userServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"
	userValidator "github.com/mohamadrezamomeni/momo/validator/user"
)

type Handler struct {
	userSvc   UserService
	validator *userValidator.Validator
	authSvc   *auth.Auth
}

type UserService interface {
	Create(*userServiceDto.AddUser) (*entity.User, error)
	Filter() ([]*entity.User, error)
	ApproveUser(id string) error
}

func New(userSvc UserService, validator *userValidator.Validator, authSvc *auth.Auth) *Handler {
	return &Handler{
		userSvc:   userSvc,
		validator: validator,
		authSvc:   authSvc,
	}
}
