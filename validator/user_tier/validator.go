package user

import (
	tierService "github.com/mohamadrezamomeni/momo/service/tier"
	userService "github.com/mohamadrezamomeni/momo/service/user"
)

type Validator struct {
	userSvc *userService.User
	tierSvc *tierService.Tier
}

func New(
	userSvc *userService.User,
	tierSvc *tierService.Tier,
) *Validator {
	return &Validator{
		userSvc: userSvc,
		tierSvc: tierSvc,
	}
}
