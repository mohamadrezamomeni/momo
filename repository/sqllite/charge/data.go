package charge

import (
	chargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/charge"
	"github.com/mohamadrezamomeni/momo/entity"
)

var charge1 = &chargeRepositoryDto.CreateDto{
	Status:    entity.PendingStatusCharge,
	UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d479",
	Detail:    "hello",
	InboundID: "12",
}
