package charge

import "github.com/mohamadrezamomeni/momo/service/auth"

type Handler struct {
	chargeSvc ChargeService
	authSvc   *auth.Auth
}

type ChargeService interface {
	ApproveCharge(string) error
}

func New(chargeSvc ChargeService, authSvc *auth.Auth) *Handler {
	return &Handler{
		chargeSvc: chargeSvc,
		authSvc:   authSvc,
	}
}
