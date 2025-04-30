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

func (mp *MockProxy) AddInbound(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) error {
	mp.addInboundData = inpt
	mp.addInboundVPNType = VPNType
	return nil
}

func (mp *MockProxy) DisableInbound(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) error {
	mp.disableInboundData = inpt
	mp.disableInboundVPNType = VPNType
	return nil
}

func (mp *MockProxy) GetTraffic(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) (*vpnSerializer.Traffic, error) {
	mp.getTrafficInboundData = inpt
	mp.getTrafficInboundVPNType = VPNType
	return &vpnSerializer.Traffic{
		Download: 20,
		Upload:   20,
	}, nil
}

func (mp *MockProxy) Close() {
	mp.addInboundData = nil
	mp.addInboundVPNType = 0

	mp.disableInboundData = nil
	mp.addInboundVPNType = 0

	mp.getTrafficInboundData = nil
	mp.getTrafficInboundVPNType = 0
}
