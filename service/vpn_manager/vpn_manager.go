package vpnmanager

import (
	vpnManagerDto "momo/dto/repository/vpn_manager"
	"momo/entity"
	vpnProxy "momo/proxy/vpn"
)

type VPNRepo interface {
	Filter(*vpnManagerDto.FilterVPNs) ([]*entity.VPN, error)
}

type VPNService struct {
	vpnRepo VPNRepo
}

func New(repo VPNRepo) *VPNService {
	return &VPNService{
		vpnRepo: repo,
	}
}

func (v *VPNService) MakeProxy() (vpnProxy.IProxyVPN, error) {
	vpns, err := v.vpnRepo.Filter(&vpnManagerDto.FilterVPNs{})
	if err != nil {
		return nil, err
	}

	cfgs := make([]*vpnProxy.VPNConfig, 0)

	for _, vpn := range vpns {
		cfgs = append(cfgs, &vpnProxy.VPNConfig{
			Port:    vpn.ApiPort,
			Domain:  vpn.Domain,
			VPNType: vpn.VPNType,
		})
	}

	return vpnProxy.New(cfgs), nil
}
