package inbound

import "github.com/mohamadrezamomeni/momo/entity"

type Validator struct {
	userSvc UserService
}

type UserService interface {
	FindByID(string) (*entity.User, error)
}

func New(userSvc UserService) *Validator {
	return &Validator{
		userSvc: userSvc,
	}
}
