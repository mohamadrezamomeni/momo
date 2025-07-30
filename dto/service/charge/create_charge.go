package charge

import "github.com/mohamadrezamomeni/momo/entity"

type CreateChargeDto struct {
	UserID    string
	PackageID string
	InboundID string
	Detail    string
	VPNType   entity.VPNType
	VPNSource *entity.VPNSource
}
