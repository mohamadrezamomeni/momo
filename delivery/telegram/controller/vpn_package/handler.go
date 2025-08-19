package vpnpackage

import (
	vpnPackageServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_package"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Handler struct {
	vpnPackageSvc VPNPackageService
}

type VPNPackageService interface {
	FilterByUserID(string, *vpnPackageServiceDto.FilterVPNPackage) ([]*entity.VPNPackage, error)
	FindVPNPackageByID(id string) (*entity.VPNPackage, error)
}

func New(vpnPackageSvc VPNPackageService) *Handler {
	return &Handler{
		vpnPackageSvc: vpnPackageSvc,
	}
}
