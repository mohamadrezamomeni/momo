package vpn

import (
	momoError "momo/pkg/error"
)

func (p *ProxyVPN) getXray(address string) V {
	for _, xray := range p.xrays {
		if xray.GetAddress() == address {
			return xray
		}
	}
	return nil
}

func (p *ProxyVPN) addXray(address string) error {
	v := p.getXray(address)
	if v != nil {
		return momoError.DebuggingErrorf("%s address has'nt been introuduced", address)
	}
	return v.Add()
}

func (p *ProxyVPN) disableXray(address string) error {
	v := p.getXray(address)
	if v != nil {
		return momoError.DebuggingErrorf("%s address has'nt been introuduced", address)
	}
	return v.Add()
}
