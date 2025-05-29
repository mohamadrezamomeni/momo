package inboundcharge

import (
	"strconv"

	"github.com/mohamadrezamomeni/momo/entity"
)

type InboundCharge struct {
	inboundChargeRepo InboundChargeRepository
	vpnPackageSvc     VPNPackageService
	inboundRepo       InboundRepository
	chargeRepo        ChargeRepository
}

type InboundRepository interface {
	RetriveFinishedInbounds() ([]*entity.Inbound, error)
	FindInboundByID(string) (*entity.Inbound, error)
}

type ChargeRepository interface {
	RetriveAvailbleChargesForInbounds(inboundIDs []string) ([]*entity.Charge, error)
}

type InboundChargeRepository interface {
	AssignChargeToInbound(
		inbound *entity.Inbound,
		charge *entity.Charge,
		vpnPackage *entity.VPNPackage,
	) error
}

type VPNPackageService interface {
	FindVPNPackageByID(string) (*entity.VPNPackage, error)
}

func New(
	inboundChargeRepo InboundChargeRepository,
	vpnPackageSvc VPNPackageService,
	inboundRepo InboundRepository,
	chargeRepo ChargeRepository,
) *InboundCharge {
	return &InboundCharge{
		inboundChargeRepo: inboundChargeRepo,
		vpnPackageSvc:     vpnPackageSvc,
		inboundRepo:       inboundRepo,
		chargeRepo:        chargeRepo,
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
		inboundIDs = append(inboundIDs, strconv.Itoa(inbound.ID))
	}
	return inboundIDs
}
