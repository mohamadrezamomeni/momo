package vpn

import (
	proxyVpnDto "momo/dto/proxy/vpn"
	"momo/entity"
	momoError "momo/pkg/error"
	"momo/proxy/vpn/internal/xray"
	vpnSerializer "momo/proxy/vpn/serializer"
)

type IVPN interface {
	Add(*proxyVpnDto.Inbound) error
	Disable(*proxyVpnDto.Inbound) error
	GetTraffic(*proxyVpnDto.Inbound) (*vpnSerializer.Traffic, error)
	DoesExist(*proxyVpnDto.Inbound) (bool, error)
	GetAddress() string
	Test() error
	Close()
}

type VPN struct {
	VPNType entity.VPNType
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
	case entity.XRAY_VPN:
		x, err := xray.New(&xray.XrayConfig{
			Address: cfg.Domain,
			ApiPort: cfg.Port,
		})
		if err == nil {
			v <- &VPN{
				VPNType: entity.XRAY_VPN,
				V:       x,
			}
		} else {
			errs <- err
		}
	}
}

func (p *ProxyVPN) retriveVPN(address string, VPNType entity.VPNType) IVPN {
	for _, v := range p.vpns {
		if v.V.GetAddress() == address && VPNType == v.VPNType {
			return v.V
		}
	}
	return nil
}

func (p *ProxyVPN) AddInbound(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) error {
	v := p.retriveVPN(inpt.Address, VPNType)
	if v == nil {
		return momoError.DebuggingErrorf("%v vpn has'nt been introuduced with address %s", VPNType, inpt.Address)
	}
	return v.Add(inpt)
}

func (p *ProxyVPN) DisableInbound(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) error {
	v := p.retriveVPN(inpt.Address, VPNType)
	if v == nil {
		return momoError.DebuggingErrorf("%v vpn has'nt been introuduced with address %s", VPNType, inpt.Address)
	}
	return v.Disable(inpt)
}

func (p *ProxyVPN) GetTraffic(inpt *proxyVpnDto.Inbound, VPNType entity.VPNType) (*vpnSerializer.Traffic, error) {
	v := p.retriveVPN(inpt.Address, VPNType)
	if v == nil {
		return &vpnSerializer.Traffic{}, momoError.DebuggingErrorf("%v vpn has'nt been introuduced with address %s", VPNType, inpt.Address)
	}
	return v.GetTraffic(inpt)
}

func (p *ProxyVPN) Close() {
	for _, vpn := range p.vpns {
		vpn.V.Close()
	}
}
