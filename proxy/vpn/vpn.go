package vpn

import (
	"sync"

	proxyVpnDto "momo/dto/proxy/vpn"
	"momo/entity"
	momoError "momo/pkg/error"
	"momo/proxy/vpn/internal/xray"
	vpnSerializer "momo/proxy/vpn/serializer"
)

type VPN interface {
	Add(*proxyVpnDto.Inbound) error
	Disable(*proxyVpnDto.Inbound) error
	GetTraffic(*proxyVpnDto.Inbound) (*vpnSerializer.Traffic, error)
	DoesExist(*proxyVpnDto.Inbound) (bool, error)
	GetAddress() string
	Test() error
	Close()
	GetType() entity.VPNType
}

type ProxyVPN struct {
	vpns []VPN
}

type IProxyVPN interface {
	AddInbound(*proxyVpnDto.Inbound) error
	DisableInbound(*proxyVpnDto.Inbound) error
	GetTraffic(*proxyVpnDto.Inbound) (*vpnSerializer.Traffic, error)
	Close()
	Test(*proxyVpnDto.Monitor) error
}

func New(cfgs []*VPNConfig) IProxyVPN {
	vpns := make([]VPN, 0)
	v := make(chan VPN, len(cfgs))

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

func addToVPN(cfg *VPNConfig, v chan<- VPN, wg *sync.WaitGroup) {
	defer wg.Done()
	switch cfg.VPNType {
	case entity.XRAY_VPN:
		x, err := xray.New(&xray.XrayConfig{
			Address: cfg.Domain,
			ApiPort: cfg.Port,
		})
		if err == nil {
			v <- x
		}
	}
}

func (p *ProxyVPN) retriveVPN(address string, VPNType entity.VPNType) VPN {
	for _, v := range p.vpns {
		if v.GetAddress() == address && VPNType == v.GetType() {
			return v
		}
	}
	return nil
}

func (p *ProxyVPN) AddInbound(inpt *proxyVpnDto.Inbound) error {
	scope := "vpnProxy.addInbound"

	v := p.retriveVPN(inpt.Address, inpt.VPNType)
	if v == nil {
		return momoError.Scope(scope).DebuggingErrorf("the result wasn't found the input is %+v", inpt)
	}
	return v.Add(inpt)
}

func (p *ProxyVPN) DisableInbound(inpt *proxyVpnDto.Inbound) error {
	scope := "vpnProxy.DisableInbound"

	v := p.retriveVPN(inpt.Address, inpt.VPNType)
	if v == nil {
		return momoError.Scope(scope).DebuggingErrorf("the result wasn't found the input is %+v", inpt)
	}
	return v.Disable(inpt)
}

func (p *ProxyVPN) GetTraffic(inpt *proxyVpnDto.Inbound) (*vpnSerializer.Traffic, error) {
	scope := "vpnProxy.GetTraffic"

	v := p.retriveVPN(inpt.Address, inpt.VPNType)
	if v == nil {
		return nil, momoError.Scope(scope).DebuggingErrorf("the result wasn't found the input is %+v", inpt)
	}
	return v.GetTraffic(inpt)
}

func (p *ProxyVPN) Test(inpt *proxyVpnDto.Monitor) error {
	scope := "vpnProxy.Test"

	v := p.retriveVPN(inpt.Address, inpt.VPNType)
	if v == nil {
		return momoError.Scope(scope).DebuggingErrorf("the result wasn't found the input is %+v", inpt)
	}
	return v.Test()
}

func (p *ProxyVPN) Close() {
	for _, vpn := range p.vpns {
		vpn.Close()
	}
}
