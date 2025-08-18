package usertier

import (
	userTierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user_tier"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"
	userTierValidation "github.com/mohamadrezamomeni/momo/validator/user_tier"
)

type Handler struct {
	userTierSvc       UserTierService
	authSvc           *auth.Auth
	userTierValidator *userTierValidation.Validator
}

type UserTierService interface {
	Create(*userTierServiceDto.Create) error
	FilterTiersByUser(string) ([]*entity.Tier, error)
	Delete(*userTierServiceDto.IdentifyUserTier) error
}

func New(userTierSvc UserTierService, authSvc *auth.Auth, userTierValidator *userTierValidation.Validator) *Handler {
	return &Handler{
		userTierSvc:       userTierSvc,
		authSvc:           authSvc,
		userTierValidator: userTierValidator,
	}
}
