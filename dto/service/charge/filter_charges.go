package charge

import "github.com/mohamadrezamomeni/momo/entity"

type FilterCharges struct {
	Statuses  []entity.ChargeStatus
	InboundID string
	UserID    string
}
