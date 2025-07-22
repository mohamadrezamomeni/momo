package adapter

import (
	proxyVpnDto "github.com/mohamadrezamomeni/momo/dto/proxy/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	vpnProxy "github.com/mohamadrezamomeni/momo/proxy/vpn"
	vpnSerializer "github.com/mohamadrezamomeni/momo/proxy/vpn/serializer"
)

type ProxyVPN interface {
	AddInbound(*proxyVpnDto.Inbound) error
	DisableInbound(*proxyVpnDto.Inbound) error
	GetTraffic(*proxyVpnDto.Inbound) (*vpnSerializer.Traffic, error)
	IsInboundActive(*proxyVpnDto.Inbound) (bool, error)
	Close()
	Test(*proxyVpnDto.Monitor) error
}

type AdapterVPnProxyigFactoryConfig struct {
	Domain  string
	Port    string
	VPNType entity.VPNType
}

func AdaptedVPNProxyFactory(adapterConfigs []*AdapterVPnProxyigFactoryConfig) ProxyVPN {
	cfgs := make([]*vpnProxy.VPNConfig, 0)

	for _, adapterConfig := range adapterConfigs {
		cfgs = append(cfgs, &vpnProxy.VPNConfig{
			Domain:  adapterConfig.Domain,
			Port:    adapterConfig.Port,
			VPNType: adapterConfig.VPNType,
		})
	}
	return vpnProxy.New(cfgs)
}
