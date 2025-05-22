package charge

import "github.com/mohamadrezamomeni/momo/entity"

type FilterChargesDto struct {
	UserID    string
	InboundID string
	Status    entity.ChargeStatus
}
