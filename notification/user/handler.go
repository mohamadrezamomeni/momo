package user

import "github.com/mohamadrezamomeni/momo/entity"

type Handler struct {
	userSvc UserService
}

type UserService interface {
	FindByID(id string) (*entity.User, error)
}

func New(userSvc UserService) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}
