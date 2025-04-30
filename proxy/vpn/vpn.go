package vpn

import (
	"sync"

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

type IProxyVPN interface {
	AddInbound(*proxyVpnDto.Inbound, entity.VPNType) error
	DisableInbound(*proxyVpnDto.Inbound, entity.VPNType) error
	GetTraffic(*proxyVpnDto.Inbound, entity.VPNType) (*vpnSerializer.Traffic, error)
	Close()
}

func New(cfgs []*VPNConfig) IProxyVPN {
	vpns := make([]*VPN, 0)
	v := make(chan *VPN, len(cfgs))

	var wg sync.WaitGroup

	for i := 0; i < len(cfgs); i++ {
		cfg := cfgs[i]
		wg.Add(1)
		go addToVPN(cfg, v, &wg)
	}

	go func() {
		wg.Wait()
		close(v)
	}()

	for vpn := range v {
		vpns = append(vpns, vpn)
	}

	return &ProxyVPN{
		vpns: vpns,
	}
}

func addToVPN(cfg *VPNConfig, v chan<- *VPN, wg *sync.WaitGroup) {
	defer wg.Done()
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
