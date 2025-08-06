package inboundcharge

import (
	"encoding/json"
	"time"

	inboundChargeDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound_charge"
	eventServiceDto "github.com/mohamadrezamomeni/momo/dto/service/event"
	"github.com/mohamadrezamomeni/momo/entity"
	chargeEvent "github.com/mohamadrezamomeni/momo/event/charge"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

type InboundCharge struct {
	inboundChargeRepo InboundChargeRepository
	vpnPackageSvc     VPNPackageService
	inboundRepo       InboundRepository
	chargeRepo        ChargeRepository
	inboundSvc        InboundService
	eventSvc          EventService
	chargeSvc         ChargeService
}

type InboundRepository interface {
	RetriveFinishedInbounds() ([]*entity.Inbound, error)
	FindInboundByID(string) (*entity.Inbound, error)
}

type ChargeRepository interface {
	RetriveAvailbleChargesForInbounds(inboundIDs []string) ([]*entity.Charge, error)
	RetriveChargesApprovedWithoutInbound() ([]*entity.Charge, error)
	FindChargeByID(id string) (*entity.Charge, error)
}

type InboundChargeRepository interface {
	AssignChargeToInbound(
		inbound *entity.Inbound,
		charge *entity.Charge,
		vpnPackage *entity.VPNPackage,
	) error
	CreateInbound(
		chargeID string,
		createInboundByCharge *inboundChargeDto.CreateInboundByCharge,
	) error
}

type InboundService interface {
	GenerateInboundTag(string, string) string
}

type VPNPackageService interface {
	FindVPNPackageByID(string) (*entity.VPNPackage, error)
}

type EventService interface {
	Create(*eventServiceDto.CreateEventDto)
}

type ChargeService interface {
	Approve(*entity.Charge) error
}

func New(
	inboundChargeRepo InboundChargeRepository,
	vpnPackageSvc VPNPackageService,
	inboundRepo InboundRepository,
	chargeRepo ChargeRepository,
	inboundSvc InboundService,
	eventSvc EventService,
	chargeSvc ChargeService,
) *InboundCharge {
	return &InboundCharge{
		chargeSvc:         chargeSvc,
		eventSvc:          eventSvc,
		inboundChargeRepo: inboundChargeRepo,
		vpnPackageSvc:     vpnPackageSvc,
		inboundRepo:       inboundRepo,
		chargeRepo:        chargeRepo,
		inboundSvc:        inboundSvc,
	}
}

func (ic *InboundCharge) ChargeInbound(charge *entity.Charge) error {
	vpnPackage, err := ic.vpnPackageSvc.FindVPNPackageByID(charge.PackageID)
	if err != nil {
		return err
	}

	inbound, err := ic.inboundRepo.FindInboundByID(charge.InboundID)
	if err != nil {
		return err
	}

	err = ic.inboundChargeRepo.AssignChargeToInbound(inbound, charge, vpnPackage)
	if err != nil {
		return err
	}
	return nil
}

func (ic *InboundCharge) ChargeInbounds() {
	inbounds, err := ic.inboundRepo.RetriveFinishedInbounds()
	if err != nil {
		return
	}
	inboundIDs := ic.getInboundIDs(inbounds)

	charges, err := ic.chargeRepo.RetriveAvailbleChargesForInbounds(inboundIDs)

	for _, charge := range charges {
		ic.ChargeInbound(charge)
	}
}

func (ic *InboundCharge) getInboundIDs(inbounds []*entity.Inbound) []string {
	inboundIDs := []string{}
	for _, inbound := range inbounds {
		inboundIDs = append(inboundIDs, inbound.ID)
	}
	return inboundIDs
}

func (ic *InboundCharge) ApproveCharge(id string) error {
	scope := "inboundCharge.approveCharge"

	var err error
	charge, err := ic.chargeRepo.FindChargeByID(id)
	if err != nil {
		return err
	}
	if charge.InboundID != "" {
		err = ic.chargeSvc.Approve(charge)
	} else {
		err = ic.createInboundByCharge(charge)
	}

	if err != nil {
		return err
	}

	chargeApprovingString, err := json.Marshal(chargeEvent.ApproveChargeEvent{
		ID: id,
	})
	if err != nil {
		momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	ic.eventSvc.Create(&eventServiceDto.CreateEventDto{
		Data: string(chargeApprovingString),
		Name: "charge_approve",
	})
	return nil
}

func (ic *InboundCharge) createInboundByCharge(charge *entity.Charge) error {
	pkg, err := ic.vpnPackageSvc.FindVPNPackageByID(charge.PackageID)
	if err != nil {
		return nil
	}
	now := time.Now()
	end := now.AddDate(0, int(pkg.Months), int(pkg.Days))

	tag := ic.inboundSvc.GenerateInboundTag(charge.Country, charge.UserID)

	err = ic.inboundChargeRepo.CreateInbound(charge.ID, &inboundChargeDto.CreateInboundByCharge{
		Tag:          tag,
		TrafficLimit: pkg.TrafficLimit,
		Start:        now,
		End:          end,
		VPNType:      charge.VPNType,
		Country:      charge.Country,
		UserID:       charge.UserID,
		Protocol:     "vmess",
	})
	if err != nil {
		return err
	}

	return nil
}
