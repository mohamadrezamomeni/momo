package tier

import (
	tierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/tier"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"
)

type Handler struct {
	tierSvc TierService
	authSvc *auth.Auth
}

type TierService interface {
	Create(*tierServiceDto.CreateTier) (*entity.Tier, error)
	Filter() ([]*entity.Tier, error)
	Update(string, *tierServiceDto.Update) error
}

func New(tierSvc TierService, authSvc *auth.Auth) *Handler {
	return &Handler{
		tierSvc: tierSvc,
		authSvc: authSvc,
	}
}
