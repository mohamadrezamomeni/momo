package vpn

import (
	"momo/entity"
	vpnProxy "momo/proxy/vpn"
	"momo/repository/sqllite/vpnmanager/dto"
)

type vpnRepo interface {
	Filter(*dto.FilterVPNs) ([]*entity.VPN, error)
}

type VPNService struct {
	vpnRepo vpnRepo
}

func New(repo vpnRepo) *VPNService {
	return &VPNService{
		vpnRepo: repo,
	}
}

func (v *VPNService) MakeProxy() (*vpnProxy.ProxyVPN, error) {
	vpns, err := v.vpnRepo.Filter(&dto.FilterVPNs{})
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
