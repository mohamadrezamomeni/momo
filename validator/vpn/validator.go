package vpn

import "github.com/mohamadrezamomeni/momo/entity"

type Validator struct {
	vpnSourceSvc VPNSourceService
}

type VPNSourceService interface {
	Find(string) (*entity.VPNSource, error)
}

func New(vpnSourceService VPNSourceService) *Validator {
	return &Validator{
		vpnSourceSvc: vpnSourceService,
	}
}
