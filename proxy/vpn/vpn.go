package vpn

import (
	momoError "momo/pkg/error"
	"momo/proxy/vpn/dto"
	"momo/proxy/vpn/internal/xray"
)

type V interface {
	Add() error
	Disable() error
	GetTraffic()
	GetAddress() string
}

type ProxyVPN struct {
	xrays []V
}

func New() *ProxyVPN {
	return &ProxyVPN{
		xrays: make([]V, 0),
	}
}

func (p *ProxyVPN) AddVpn(vpnType string, domain string, port string) {
	switch vpnType {
	case xray_vpn:
		p.xrays = append(p.xrays, xray.New(&xray.XrayConfig{
			Address: domain,
			ApiPort: port,
		}))
	}
}

func (p *ProxyVPN) AddInbound(inpt *dto.Inbound, vpnType string) (err error) {
	switch vpnType {
	case xray_vpn:
		p.addXray(inpt.Address)
	default:
		err = momoError.DebuggingErrorf("the vpnType that was given was wrong %s", vpnType)
	}
	return
}

func (p *ProxyVPN) DisableInbound(inpt *dto.Inbound, vpnType string) (err error) {
	switch vpnType {
	case xray_vpn:
		p.disableXray(inpt.Address)
	default:
		err = momoError.DebuggingErrorf("the vpnType that was given was wrong %s", vpnType)
	}
	return
}

func (p *ProxyVPN) GetTraffix(inpt *dto.Inbound, vpnType string) error {
	return nil
}
