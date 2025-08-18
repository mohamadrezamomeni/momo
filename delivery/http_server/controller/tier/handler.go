package tier

import (
	tierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/tier"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"

	tierValidation "github.com/mohamadrezamomeni/momo/validator/tier"
)

type Handler struct {
	tierSvc      TierService
	authSvc      *auth.Auth
	tierValiator *tierValidation.Validator
}

type TierService interface {
	Create(*tierServiceDto.CreateTier) (*entity.Tier, error)
	Filter() ([]*entity.Tier, error)
	Update(string, *tierServiceDto.Update) error
}

func New(
	tierSvc TierService,
	authSvc *auth.Auth,
	tierValiator *tierValidation.Validator,
) *Handler {
	return &Handler{
		tierSvc:      tierSvc,
		authSvc:      authSvc,
		tierValiator: tierValiator,
	}
}
