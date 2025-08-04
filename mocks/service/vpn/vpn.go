package vpn

import (
	"strconv"

	"github.com/mohamadrezamomeni/momo/adapter"
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	mockProxyVPN "github.com/mohamadrezamomeni/momo/mocks/proxy/vpn"
)

type MockVPN struct {
	vpns []*entity.VPN
	idx  int
}

func New() *MockVPN {
	return &MockVPN{
		vpns: make([]*entity.VPN, 0),
		idx:  0,
	}
}

func (mv *MockVPN) Create(createVPNDto *vpnServiceDto.CreateVPN) (*entity.VPN, error) {
	vpn := &entity.VPN{
		ID:        strconv.Itoa(mv.idx),
		Domain:    createVPNDto.Domain,
		VPNType:   createVPNDto.VpnType,
		Country:   createVPNDto.Country,
		UserCount: 20,
		IsActive:  true,
		StartPort: createVPNDto.StartPort,
		EndPort:   createVPNDto.EndPort,
	}
	mv.idx += 1
	mv.vpns = append(mv.vpns, vpn)
	return vpn, nil
}

func (mv *MockVPN) MakeProxy() (adapter.ProxyVPN, error) {
	return &mockProxyVPN.MockProxy{}, nil
}

func (mv *MockVPN) GetAvailableVPNSourceDomains(vpnSources []string, vpnTypes []entity.VPNType) ([]*entity.VPN, error) {
	ret := []*entity.VPN{}

	vpnSourcesMap := make(map[string]struct{})
	for _, vpnSource := range vpnSources {
		vpnSourcesMap[vpnSource] = struct{}{}
	}
	vpnTypeMap := make(map[entity.VPNType]struct{})
	for _, vpnType := range vpnTypes {
		vpnTypeMap[vpnType] = struct{}{}
	}

	for _, vpn := range mv.vpns {
		_, isExistSource := vpnSourcesMap[vpn.Country]
		_, isExistVPNType := vpnTypeMap[vpn.VPNType]
		if isExistSource && isExistVPNType {
			ret = append(ret, vpn)
		}
	}
	return ret, nil
}

func (mv *MockVPN) DeleteAll() {
	mv.vpns = []*entity.VPN{}
}
