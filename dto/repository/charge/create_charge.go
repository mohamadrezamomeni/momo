package charge

import "github.com/mohamadrezamomeni/momo/entity"

type CreateDto struct {
	InboundID string
	Detail    string
	Status    entity.ChargeStatus
	UserID    string
	PackageID string
	Country   string
	VPNType   entity.VPNType
}
