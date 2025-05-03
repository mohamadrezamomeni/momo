package vpnmanager

import (
	"testing"

	"github.com/mohamadrezamomeni/momo/adapter"

	vpnProxyMock "github.com/mohamadrezamomeni/momo/mocks/proxy/vpn"
	vpnRepositoryMock "github.com/mohamadrezamomeni/momo/mocks/repository/vpn"
)

func register() (*VPNService, *vpnRepositoryMock.MockVPN) {
	vpnRepo := vpnRepositoryMock.New()

	vpnSvc := New(vpnRepo, func(adapterConfigs []*adapter.AdapterVPnProxyigFactoryConfig) adapter.ProxyVPN {
		return &vpnProxyMock.MockProxy{
			Count: 0,
		}
	})
	return vpnSvc, vpnRepo
}

func TestMonitorVPNs(t *testing.T) {
	vpnSvc, vpnRepo := register()

	vpnCreated1, _ := vpnRepo.Create(vpn1)
	vpnCreated2, _ := vpnRepo.Create(vpn2)
	vpnCreated3, _ := vpnRepo.Create(vpn3)

	vpnSvc.MonitorVPNs()

	vpnFound1 := vpnRepo.FindByID(vpnCreated1.ID)
	vpnFound2 := vpnRepo.FindByID(vpnCreated2.ID)
	vpnFound3 := vpnRepo.FindByID(vpnCreated3.ID)

	if vpnFound3.IsActive != true {
		t.Error("vpn3 must be active")
	}

	if vpnFound2.IsActive != false {
		t.Error("vpn2 must be deactive")
	}

	if vpnFound1.IsActive != true {
		t.Error("vpn1 must be active")
	}
}
