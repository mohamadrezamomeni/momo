package inboundcharge

import (
	"github.com/mohamadrezamomeni/momo/entity"
)

type InboundCharge struct {
	inboundChargeRepo InboundChargeRepository
	vpnPackageSvc     VPNPackageService
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
) *InboundCharge {
	return &InboundCharge{
		inboundChargeRepo: inboundChargeRepo,
		vpnPackageSvc:     vpnPackageSvc,
	}
}

func (ic *InboundCharge) ChargeInbound(inbound *entity.Inbound, charge *entity.Charge) error {
	vpnPackage, err := ic.vpnPackageSvc.FindVPNPackageByID(charge.PackageID)
	if err != nil {
		return err
	}

	err = ic.inboundChargeRepo.AssignChargeToInbound(inbound, charge, vpnPackage)
	if err != nil {
		return err
	}
	return nil
}
