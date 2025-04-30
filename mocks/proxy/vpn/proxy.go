package vpn

import (
	proxyVpnDto "momo/dto/proxy/vpn"
	"momo/entity"
	vpnSerializer "momo/proxy/vpn/serializer"
)

type MockProxy struct {
	addInboundData    *proxyVpnDto.Inbound
	addInboundVPNType entity.VPNType

	disableInboundData    *proxyVpnDto.Inbound
	disableInboundVPNType entity.VPNType

	getTrafficInboundData    *proxyVpnDto.Inbound
	getTrafficInboundVPNType entity.VPNType
}

func (mp *MockProxy) AddInbound(inpt *proxyVpnDto.Inbound) error {
	mp.addInboundData = inpt
	return nil
}

func (mp *MockProxy) DisableInbound(inpt *proxyVpnDto.Inbound) error {
	mp.disableInboundData = inpt
	return nil
}

func (mp *MockProxy) GetTraffic(inpt *proxyVpnDto.Inbound) (*vpnSerializer.Traffic, error) {
	mp.getTrafficInboundData = inpt
	return &vpnSerializer.Traffic{
		Download: 20,
		Upload:   20,
	}, nil
}

func (mp *MockProxy) Close() {
	mp.addInboundData = nil

	mp.disableInboundData = nil

	mp.getTrafficInboundData = nil
}
