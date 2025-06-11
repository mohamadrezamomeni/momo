package vpnmanager

import (
	"github.com/mohamadrezamomeni/momo/adapter"
	vpnProxyDto "github.com/mohamadrezamomeni/momo/dto/proxy/vpn"
	vpnManagerRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_manager"
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
)

type VPNRepo interface {
	Filter(*vpnManagerRepositoryDto.FilterVPNs) ([]*entity.VPN, error)
	ActiveVPN(int) error
	DeactiveVPN(int) error
	Create(*vpnManagerRepositoryDto.AddVPN) (*entity.VPN, error)
}

type AdaptedVPNProxy func(adapterConfigs []*adapter.AdapterVPnProxyigFactoryConfig) adapter.ProxyVPN

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

func (v *VPNService) Create(createVPNDto *vpnServiceDto.CreateVPN) (*entity.VPN, error) {
	return v.vpnRepo.Create(&vpnManagerRepositoryDto.AddVPN{
		ApiPort:   createVPNDto.Port,
		VPNType:   createVPNDto.VpnType,
		Domain:    createVPNDto.Domain,
		IsActive:  false,
		UserCount: createVPNDto.UserCount,
		Country:   createVPNDto.Country,
	})
}

func (v *VPNService) MonitorVPNs() {
	vpns, err := v.load()
	if err != nil {
		return
	}
	cfgs := make([]*adapter.AdapterVPnProxyigFactoryConfig, 0)

	for _, vpn := range vpns {
		cfgs = append(cfgs, &adapter.AdapterVPnProxyigFactoryConfig{
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

func (v *VPNService) Filter(vpnFilterDto *vpnServiceDto.FilterVPNs) ([]*entity.VPN, error) {
	return v.vpnRepo.Filter(&vpnManagerRepositoryDto.FilterVPNs{
		Domain:  vpnFilterDto.Domain,
		VPNType: vpnFilterDto.VPNType,
	})
}

func (v *VPNService) MakeProxy() (adapter.ProxyVPN, error) {
	vpns, err := v.load()
	if err != nil {
		return nil, err
	}

	cfgs := make([]*adapter.AdapterVPnProxyigFactoryConfig, 0)

	for _, vpn := range vpns {
		cfgs = append(cfgs, &adapter.AdapterVPnProxyigFactoryConfig{
			Port:    vpn.ApiPort,
			Domain:  vpn.Domain,
			VPNType: vpn.VPNType,
		})
	}

	return v.adaptedVPNProxy(cfgs), nil
}

func (v *VPNService) load() ([]*entity.VPN, error) {
	vpns, err := v.vpnRepo.Filter(&vpnManagerRepositoryDto.FilterVPNs{})
	if err != nil {
		return nil, err
	}
	return vpns, nil
}
