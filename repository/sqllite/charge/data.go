package charge

import (
	chargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/charge"
	"github.com/mohamadrezamomeni/momo/entity"
)

var charge1 = &chargeRepositoryDto.CreateDto{
	Status:    entity.PendingStatusCharge,
	Detail:    "hello",
	InboundID: "12",
}
