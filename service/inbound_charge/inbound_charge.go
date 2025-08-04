package inboundcharge

import (
	"time"

	inboundChargeDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound_charge"
	"github.com/mohamadrezamomeni/momo/entity"
)

type InboundCharge struct {
	inboundChargeRepo InboundChargeRepository
	vpnPackageSvc     VPNPackageService
	inboundRepo       InboundRepository
	chargeRepo        ChargeRepository
	inboundSvc        InboundService
}

type InboundRepository interface {
	RetriveFinishedInbounds() ([]*entity.Inbound, error)
	FindInboundByID(string) (*entity.Inbound, error)
}

type ChargeRepository interface {
	RetriveAvailbleChargesForInbounds(inboundIDs []string) ([]*entity.Charge, error)
	RetriveChargesApprovedWithoutInbound() ([]*entity.Charge, error)
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

func New(
	inboundChargeRepo InboundChargeRepository,
	vpnPackageSvc VPNPackageService,
	inboundRepo InboundRepository,
	chargeRepo ChargeRepository,
	inboundSvc InboundService,
) *InboundCharge {
	return &InboundCharge{
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

func (ic *InboundCharge) CreateInbounds() {
	charges, err := ic.chargeRepo.RetriveChargesApprovedWithoutInbound()
	if err != nil {
		return
	}
	ic.createInboundsByCharges(charges)
}

func (ic *InboundCharge) createInboundsByCharges(charges []*entity.Charge) {
	for _, charge := range charges {
		ic.createInboundByCharge(charge)
	}
}

func (ic *InboundCharge) createInboundByCharge(charge *entity.Charge) {
	pkg, err := ic.vpnPackageSvc.FindVPNPackageByID(charge.PackageID)
	if err != nil {
		return
	}
	now := time.Now()
	end := now.AddDate(0, int(pkg.Months), int(pkg.Days))

	tag := ic.inboundSvc.GenerateInboundTag(charge.Country, charge.UserID)

	ic.inboundChargeRepo.CreateInbound(charge.ID, &inboundChargeDto.CreateInboundByCharge{
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
		return
	}
}
