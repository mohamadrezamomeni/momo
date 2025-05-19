package inbound

import (
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	vpnPackageServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_package"
	"github.com/mohamadrezamomeni/momo/entity"
	inboundValidator "github.com/mohamadrezamomeni/momo/validator/inbound"
)

type Handler struct {
	userSvc          UserService
	inboundSvc       InboundService
	inboundValidator *inboundValidator.Validator
	vpnPackageSvc    VPNPackageService
}

type InboundService interface {
	Filter(inpt *inboundServiceDto.FilterInbounds) ([]*entity.Inbound, error)
	Create(inpt *inboundServiceDto.CreateInbound) (*entity.Inbound, error)
	FindInboundByID(string) (*entity.Inbound, error)
}

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
}

type VPNPackageService interface {
	Filter(*vpnPackageServiceDto.FilterVPNPackage) ([]*entity.VPNPackage, error)
	FindVPNPackageByID(id string) (*entity.VPNPackage, error)
}

func New(
	userSvc UserService,
	inboundSvc InboundService,
	vpnPackageSvc VPNPackageService,
	inboundValidator *inboundValidator.Validator,
) *Handler {
	return &Handler{
		vpnPackageSvc:    vpnPackageSvc,
		userSvc:          userSvc,
		inboundSvc:       inboundSvc,
		inboundValidator: inboundValidator,
	}
}
