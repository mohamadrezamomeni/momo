package vpnpackage

import (
	authSvc "github.com/mohamadrezamomeni/momo/service/auth"
	vpnPakcageSvc "github.com/mohamadrezamomeni/momo/service/vpn_package"
	vpnPackageValidation "github.com/mohamadrezamomeni/momo/validator/vpn_package"
)

type Handler struct {
	vpnPackageSvc       *vpnPakcageSvc.VPNPackage
	vpnPackageValidator *vpnPackageValidation.Validator
	authSvc             *authSvc.Auth
}

func New(
	vpnPackageSvc *vpnPakcageSvc.VPNPackage,
	validator *vpnPackageValidation.Validator,
	authSvc *authSvc.Auth,
) *Handler {
	return &Handler{
		vpnPackageSvc:       vpnPackageSvc,
		vpnPackageValidator: validator,
		authSvc:             authSvc,
	}
}
