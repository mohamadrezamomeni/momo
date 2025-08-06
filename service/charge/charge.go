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
	GetFirstApprovedInboundCharge(string) (*entity.Charge, error)
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
		VPNType:   createChargeDto.VPNType,
		Country:   createChargeDto.VPNSource,
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

func (c *Charge) FindByID(id string) (*entity.Charge, error) {
	return c.chargeRepo.FindChargeByID(id)
}

func (c *Charge) Approve(charge *entity.Charge) error {
	err := c.chargeRepo.UpdateCharge(charge.ID, &chargeRepositoryDto.UpdateChargeDto{
		Status: entity.ApprovedStatusCharge,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Charge) FindAvailbleCharge(inboundID string) (*entity.Charge, error) {
	return c.chargeRepo.GetFirstApprovedInboundCharge(inboundID)
}

func (c *Charge) FilterCharges(filterCharges *chargeServiceDto.FilterCharges) ([]*entity.Charge, error) {
	return c.chargeRepo.FilterCharges(&chargeRepositoryDto.FilterChargesDto{
		InboundID: filterCharges.InboundID,
		UserID:    filterCharges.UserID,
		Statuses:  filterCharges.Statuses,
	})
}
