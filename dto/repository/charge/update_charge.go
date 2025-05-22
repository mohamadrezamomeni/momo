package charge

import "github.com/mohamadrezamomeni/momo/entity"

type UpdateChargeDto struct {
	Status       entity.ChargeStatus
	Detail       string
	AdminComment string
}
