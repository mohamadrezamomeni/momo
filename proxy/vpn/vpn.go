package vpn

import (
	momoError "momo/pkg/error"
	"momo/proxy/vpn/dto"
	vpnDto "momo/proxy/vpn/dto"
	"momo/proxy/vpn/internal/xray"
	vpnSerializer "momo/proxy/vpn/serializer"
)

type V interface {
	Add(*vpnDto.Inbound) error
	Disable(*vpnDto.Inbound) error
	GetTraffic(*vpnDto.Inbound) (*vpnSerializer.Traffic, error)
	DoesExist(*vpnDto.Inbound) (bool, error)
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
		err = p.addXray(inpt.Address, inpt)
	default:
		err = momoError.DebuggingErrorf("the vpnType that was given was wrong %s", vpnType)
	}
	return
}

func (p *ProxyVPN) DisableInbound(inpt *dto.Inbound, vpnType string) (err error) {
	switch vpnType {
	case xray_vpn:
		err = p.disableXray(inpt.Address, inpt)
	default:
		err = momoError.DebuggingErrorf("the vpnType that was given was wrong %s", vpnType)
	}
	return
}

func (p *ProxyVPN) GetTraffix(inpt *dto.Inbound, vpnType string) (*vpnSerializer.Traffic, error) {
	switch vpnType {
	case xray_vpn:
		ret, err := p.getTrafficXray(inpt.Address, inpt)
		return ret, err
	default:
		return &vpnSerializer.Traffic{}, momoError.DebuggingErrorf("the vpnType that was given was wrong %s", vpnType)
	}
}
