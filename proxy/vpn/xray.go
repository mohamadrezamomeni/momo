package vpn

import (
	momoError "momo/pkg/error"
	"momo/proxy/vpn/dto"
	vpnSerializer "momo/proxy/vpn/serializer"
)

func (p *ProxyVPN) getXray(address string) V {
	for _, xray := range p.xrays {
		if xray.GetAddress() == address {
			return xray
		}
	}
	return nil
}

func (p *ProxyVPN) addXray(address string, inpt *dto.Inbound) error {
	v := p.getXray(address)
	if v != nil {
		return momoError.DebuggingErrorf("%s address has'nt been introuduced", address)
	}
	return v.Add(inpt)
}

func (p *ProxyVPN) disableXray(address string, inpt *dto.Inbound) error {
	v := p.getXray(address)
	if v != nil {
		return momoError.DebuggingErrorf("%s address has'nt been introuduced", address)
	}
	return v.Disable(inpt)
}

func (p *ProxyVPN) getTrafficXray(address string, inpt *dto.Inbound) (*vpnSerializer.Traffic, error) {
	v := p.getXray(address)
	if v != nil {
		return &vpnSerializer.Traffic{}, momoError.DebuggingErrorf("%s address has'nt been introuduced", address)
	}
	return v.GetTraffic(inpt)
}
