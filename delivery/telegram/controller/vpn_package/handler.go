package vpnpackage

import (
	vpnPackageServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_package"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Handler struct {
	vpnPackageSvc VPNPackageService
	userSvc       UserService
}

type VPNPackageService interface {
	Filter(*vpnPackageServiceDto.FilterVPNPackage) ([]*entity.VPNPackage, error)
	FindVPNPackageByID(id string) (*entity.VPNPackage, error)
}

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
}

func New(vpnPackageSvc VPNPackageService, userSvc UserService) *Handler {
	return &Handler{
		vpnPackageSvc: vpnPackageSvc,
		userSvc:       userSvc,
	}
}
