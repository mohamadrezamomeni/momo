package usertier

import (
	userTierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user_tier"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"
)

type Handler struct {
	userTierSvc UserTierService
	authSvc     *auth.Auth
}

type UserTierService interface {
	Create(*userTierServiceDto.Create) error
	FilterTiersByUser(string) ([]*entity.Tier, error)
	Delete(*userTierServiceDto.IdentifyUserTier) error
}

func New(userTierSvc UserTierService, authSvc *auth.Auth) *Handler {
	return &Handler{
		userTierSvc: userTierSvc,
		authSvc:     authSvc,
	}
}
