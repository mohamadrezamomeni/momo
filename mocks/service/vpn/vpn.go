package vpn

import (
	"fmt"

	"github.com/mohamadrezamomeni/momo/adapter"
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	mockProxyVPN "github.com/mohamadrezamomeni/momo/mocks/proxy/vpn"
)

type MockVPN struct {
	vpns []*entity.VPN
}

func New() *MockVPN {
	return &MockVPN{}
}

func (mv *MockVPN) Create(createVPNDto *vpnServiceDto.CreateVPN) (*entity.VPN, error) {
	id := 0
	if len(mv.vpns) != 0 {
		id = mv.vpns[len(mv.vpns)-1].ID
	}
	vpn := &entity.VPN{
		ID:        id,
		Domain:    createVPNDto.Domain,
		VPNType:   createVPNDto.VpnType,
		Country:   createVPNDto.Country,
		UserCount: 20,
		IsActive:  true,
	}
	mv.vpns = append(mv.vpns, vpn)
	return vpn, nil
}

func (mv *MockVPN) MakeProxy() (adapter.ProxyVPN, error) {
	return &mockProxyVPN.MockProxy{}, nil
}

func (mv *MockVPN) GetAvailableVPNSourceDomains(vpnSources []string) (map[string][]string, error) {
	ret := map[string][]string{}
	for _, vpnSource := range vpnSources {
		for _, v := range mv.vpns {
			fmt.Println(vpnSource, v.Country)
			if v.Country == vpnSource {
				ret[vpnSource] = append(ret[vpnSource], v.Domain)
			}
		}
	}
	return ret, nil
}

func (mv *MockVPN) DeleteAll() {
	mv.vpns = []*entity.VPN{}
}
