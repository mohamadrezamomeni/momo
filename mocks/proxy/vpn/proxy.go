package vpn

import (
	proxyVpnDto "momo/dto/proxy/vpn"
	"momo/entity"
	vpnSerializer "momo/proxy/vpn/serializer"
)

type MockProxy struct{}

func (mp *MockProxy) AddInbound(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) error {
	return nil
}

func (mp *MockProxy) DisableInbound(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) error {
	return nil
}

func (mp *MockProxy) GetTraffic(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) (*vpnSerializer.Traffic, error) {
	return &vpnSerializer.Traffic{
		Download: 20,
		Upload:   20,
	}, nil
}

func (mp *MockProxy) Close() {
}
