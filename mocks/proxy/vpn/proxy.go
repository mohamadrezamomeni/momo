package vpn

import (
	"fmt"

	proxyVpnDto "github.com/mohamadrezamomeni/momo/dto/proxy/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	vpnSerializer "github.com/mohamadrezamomeni/momo/proxy/vpn/serializer"
)

type MockProxy struct {
	addInboundData    *proxyVpnDto.Inbound
	addInboundVPNType entity.VPNType

	disableInboundData    *proxyVpnDto.Inbound
	disableInboundVPNType entity.VPNType

	getTrafficInboundData    *proxyVpnDto.Inbound
	getTrafficInboundVPNType entity.VPNType

	Count int
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

func (mp *MockProxy) Test(inpt *proxyVpnDto.Monitor) error {
	defer func() {
		mp.Count += 1
	}()
	if mp.Count%2 == 0 {
		return fmt.Errorf("")
	}
	return nil
}

func (mp *MockProxy) Close() {
	mp.addInboundData = nil

	mp.disableInboundData = nil

	mp.getTrafficInboundData = nil
}

func (mp *MockProxy) IsInboundActive(inpt *proxyVpnDto.Inbound) (bool, error) {
	return true, nil
}
