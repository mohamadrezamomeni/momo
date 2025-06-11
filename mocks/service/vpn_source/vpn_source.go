package vpnsource

import (
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

type MockVPNSource struct {
	vpnSourcesRefrence map[string]*entity.VPNSource
}

func New(vpnSources []*entity.VPNSource) *MockVPNSource {
	vpnSourcesRefrence := map[string]*entity.VPNSource{}
	for _, vpnSource := range vpnSources {
		vpnSourcesRefrence[vpnSource.Country] = vpnSource
	}
	return &MockVPNSource{
		vpnSourcesRefrence: vpnSourcesRefrence,
	}
}

func (mv *MockVPNSource) Find(country string) (*entity.VPNSource, error) {
	if vpnSource, isExist := mv.vpnSourcesRefrence[country]; isExist {
		return vpnSource, nil
	}
	return nil, momoError.Scope("").NotFound().ErrorWrite()
}
