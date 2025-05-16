package root

import "github.com/mohamadrezamomeni/momo/entity"

type Handler struct {
	userSvc UserService
}

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
}

func New(userSvc UserService) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}
