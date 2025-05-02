package vpnmanager

import (
	vpnProxyDto "momo/dto/proxy/vpn"
	vpnManagerDto "momo/dto/repository/vpn_manager"
	"momo/entity"
	vpnProxy "momo/proxy/vpn"
)

type VPNRepo interface {
	Filter(*vpnManagerDto.FilterVPNs) ([]*entity.VPN, error)
	ActiveVPN(int) error
	DeactiveVPN(int) error
}

type AdaptedVPNProxy func(adapterConfigs []*AdapterVPnProxyigFactoryConfig) ProxyVPN

type VPNService struct {
	vpnRepo         VPNRepo
	adaptedVPNProxy AdaptedVPNProxy
}

func New(repo VPNRepo, adaptedVPNProxy AdaptedVPNProxy) *VPNService {
	return &VPNService{
		vpnRepo:         repo,
		adaptedVPNProxy: adaptedVPNProxy,
	}
}

func (v *VPNService) MonitorVPNs() {
	vpns, err := v.load()
	if err != nil {
		return
	}
	cfgs := make([]*AdapterVPnProxyigFactoryConfig, 0)

	for _, vpn := range vpns {
		cfgs = append(cfgs, &AdapterVPnProxyigFactoryConfig{
			Port:    vpn.ApiPort,
			Domain:  vpn.Domain,
			VPNType: vpn.VPNType,
		})
	}
	vpnProxy := v.adaptedVPNProxy(cfgs)

	for _, vpn := range vpns {
		err := vpnProxy.Test(&vpnProxyDto.Monitor{
			Address: vpn.Domain,
			VPNType: vpn.VPNType,
		})
		if err != nil {
			v.vpnRepo.ActiveVPN(vpn.ID)
		} else {
			v.vpnRepo.DeactiveVPN(vpn.ID)
		}
	}
}

func (v *VPNService) MakeProxy() (vpnProxy.IProxyVPN, error) {
	vpns, err := v.load()
	if err != nil {
		return nil, err
	}

	cfgs := make([]*AdapterVPnProxyigFactoryConfig, 0)

	for _, vpn := range vpns {
		cfgs = append(cfgs, &AdapterVPnProxyigFactoryConfig{
			Port:    vpn.ApiPort,
			Domain:  vpn.Domain,
			VPNType: vpn.VPNType,
		})
	}

	return v.adaptedVPNProxy(cfgs), nil
}

func (v *VPNService) load() ([]*entity.VPN, error) {
	vpns, err := v.vpnRepo.Filter(&vpnManagerDto.FilterVPNs{})
	if err != nil {
		return nil, err
	}
	return vpns, nil
}
