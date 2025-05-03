package adapter

import (
	proxyVpnDto "momo/dto/proxy/vpn"
	"momo/entity"
	vpnProxy "momo/proxy/vpn"
	vpnSerializer "momo/proxy/vpn/serializer"
)

type ProxyVPN interface {
	AddInbound(*proxyVpnDto.Inbound) error
	DisableInbound(*proxyVpnDto.Inbound) error
	GetTraffic(*proxyVpnDto.Inbound) (*vpnSerializer.Traffic, error)
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
