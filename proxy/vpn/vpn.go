package vpn

import (
	momoError "momo/pkg/error"
	"momo/proxy/vpn/dto"
	vpnDto "momo/proxy/vpn/dto"
	"momo/proxy/vpn/internal/xray"
	vpnSerializer "momo/proxy/vpn/serializer"
)

type IVPN interface {
	Add(*vpnDto.Inbound) error
	Disable(*vpnDto.Inbound) error
	GetTraffic(*vpnDto.Inbound) (*vpnSerializer.Traffic, error)
	DoesExist(*vpnDto.Inbound) (bool, error)
	GetAddress() string
	Test() error
}

type VPN struct {
	VPNType VPNType
	V       IVPN
}

type ProxyVPN struct {
	vpns []*VPN
}

func New(cfgs []*VPNConfig) *ProxyVPN {
	vpns := make([]*VPN, 0)
	v := make(chan *VPN, len(cfgs))
	errs := make(chan error, len(cfgs))

	for i := 0; i < len(cfgs); i++ {
		cfg := cfgs[i]
		go addToVPN(cfg, v, errs)
	}

	for i := 0; i < len(cfgs); i++ {
		select {
		case vpn := <-v:
			vpns = append(vpns, vpn)
		case <-errs:
		}
	}
	return &ProxyVPN{
		vpns: vpns,
	}
}

func addToVPN(cfg *VPNConfig, v chan<- *VPN, errs chan<- error) {
	switch cfg.VPNType {
	case XRAY_VPN:
		x, err := xray.New(&xray.XrayConfig{
			Address: cfg.Domain,
			ApiPort: cfg.Port,
		})
		if err != nil {
			v <- &VPN{
				VPNType: XRAY_VPN,
				V:       x,
			}
		} else {
			errs <- err
		}
	}
}

func (p *ProxyVPN) retriveVPN(address string, VPNType VPNType) IVPN {
	for _, v := range p.vpns {
		if v.V.GetAddress() == address && VPNType == v.VPNType {
			return v.V
		}
	}
	return nil
}

func (p *ProxyVPN) AddInbound(inpt *dto.Inbound, VPNType string) (err error) {
	v := p.retriveVPN(inpt.Address, VPNType)
	if v == nil {
		return momoError.DebuggingErrorf("%s vpn has'nt been introuduced with address %s", VPNType, inpt.Address)
	}
	return v.Add(inpt)
}

func (p *ProxyVPN) DisableInbound(inpt *dto.Inbound, VPNType string) (err error) {
	v := p.retriveVPN(inpt.Address, VPNType)
	if v == nil {
		return momoError.DebuggingErrorf("%s vpn has'nt been introuduced with address %s", VPNType, inpt.Address)
	}
	return v.Disable(inpt)
}

func (p *ProxyVPN) GetTraffix(inpt *dto.Inbound, VPNType string) (*vpnSerializer.Traffic, error) {
	v := p.retriveVPN(inpt.Address, VPNType)
	if v == nil {
		return &vpnSerializer.Traffic{}, momoError.DebuggingErrorf("%s vpn has'nt been introuduced with address %s", VPNType, inpt.Address)
	}
	return v.GetTraffic(inpt)
}
