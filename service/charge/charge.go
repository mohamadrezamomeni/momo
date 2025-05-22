package charge

import (
	"encoding/json"

	chargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/charge"
	chargeServiceDto "github.com/mohamadrezamomeni/momo/dto/service/charge"
	eventServiceDto "github.com/mohamadrezamomeni/momo/dto/service/event"
	entity "github.com/mohamadrezamomeni/momo/entity"
	chargeEvent "github.com/mohamadrezamomeni/momo/event/charge"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

type Charge struct {
	chargeRepo ChargeRepository
	event      EventService
}

type ChargeRepository interface {
	Create(*chargeRepositoryDto.CreateDto) (*entity.Charge, error)
	FindChargeByID(string) (*entity.Charge, error)
	UpdateCharge(string, *chargeRepositoryDto.UpdateChargeDto) error
	FilterCharges(*chargeRepositoryDto.FilterChargesDto) ([]*entity.Charge, error)
}

type EventService interface {
	Create(*eventServiceDto.CreateEventDto)
}

func New(event EventService, chargeRepo ChargeRepository) *Charge {
	return &Charge{
		chargeRepo: chargeRepo,
		event:      event,
	}
}

func (c *Charge) Create(createChargeDto *chargeServiceDto.CreateChargeDto) (*entity.Charge, error) {
	scope := "charge.service.create"
	chargeCreated, err := c.chargeRepo.Create(&chargeRepositoryDto.CreateDto{
		Detail:    createChargeDto.Detail,
		Status:    entity.PendingStatusCharge,
		InboundID: createChargeDto.InboundID,
		PackageID: createChargeDto.PackageID,
		UserID:    createChargeDto.UserID,
	})
	if err != nil {
		return nil, err
	}
	chargeCreatedEvent := chargeEvent.CreateChargeEvent{
		ID: chargeCreated.ID,
	}

	chargeCreatedEventString, err := json.Marshal(chargeCreatedEvent)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).ErrorWrite()
	}
	c.event.Create(&eventServiceDto.CreateEventDto{
		Name: "chargeCreating",
		Data: string(chargeCreatedEventString),
	})

	return chargeCreated, nil
}
