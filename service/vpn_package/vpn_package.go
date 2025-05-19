package vpnpackage

import (
	vpnPackageRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_package"
	vpnPackageServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_package"
	"github.com/mohamadrezamomeni/momo/entity"
)

type VPNPackage struct {
	vpnPackageRepo VPNPackageRepository
}
type VPNPackageRepository interface {
	Update(string, *vpnPackageRepositoryDto.UpdateVPNPackage) error
	Create(*vpnPackageRepositoryDto.CreateVPNPackage) (*entity.VPNPackage, error)
	Filter(*vpnPackageRepositoryDto.FilterVPNPackage) ([]*entity.VPNPackage, error)
}

func New(vpnPackageRepo VPNPackageRepository) *VPNPackage {
	return &VPNPackage{
		vpnPackageRepo: vpnPackageRepo,
	}
}

func (vp *VPNPackage) Create(inpt *vpnPackageServiceDto.CreateVPNPackage) (*entity.VPNPackage, error) {
	return vp.vpnPackageRepo.Create(&vpnPackageRepositoryDto.CreateVPNPackage{
		Price:             inpt.Price,
		PriceTitle:        inpt.PriceTitle,
		TrafficLimitTitle: inpt.TrafficLimitTitle,
		TrafficLimit:      inpt.TrafficLimit,
		Days:              inpt.Days,
		Months:            inpt.Months,
		IsActive:          inpt.IsActive,
	})
}

func (vp *VPNPackage) Active(id string) error {
	active := true
	err := vp.vpnPackageRepo.Update(id, &vpnPackageRepositoryDto.UpdateVPNPackage{
		IsActive: &active,
	})

	return err
}

func (vp *VPNPackage) Deactive(id string) error {
	active := true
	err := vp.vpnPackageRepo.Update(id, &vpnPackageRepositoryDto.UpdateVPNPackage{
		IsActive: &active,
	})

	return err
}

func (vp *VPNPackage) Filter(inpt *vpnPackageServiceDto.FilterVPNPackage) ([]*entity.VPNPackage, error) {
	return vp.vpnPackageRepo.Filter(&vpnPackageRepositoryDto.FilterVPNPackage{
		IsActive: inpt.IsActive,
	})
}
