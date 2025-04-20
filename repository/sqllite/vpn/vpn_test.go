package vpn

import (
	"os"
	"testing"

	"momo/pkg/config"
	"momo/proxy/vpn"
	"momo/repository/migrate"
	"momo/repository/sqllite"
	"momo/repository/sqllite/vpn/dto"
)

var (
	vpnRepo *VPN

	vpnExample1 = &dto.Add_VPN{
		Domain:         "joi.com",
		ApiPort:        "62733",
		StartRangePort: 1000,
		EndRangePort:   2000,
		VPNType:        vpn.XRAY_VPN,
		IsActive:       true,
	}

	vpnExample2 = &dto.Add_VPN{
		Domain:         "joi.com",
		ApiPort:        "62733",
		StartRangePort: 1000,
		EndRangePort:   2500,
		VPNType:        vpn.XRAY_VPN,
		IsActive:       true,
	}

	vpnExample3 = &dto.Add_VPN{
		Domain:         "jordan.com",
		ApiPort:        "62733",
		StartRangePort: 3000,
		EndRangePort:   3500,
		VPNType:        vpn.XRAY_VPN,
		IsActive:       true,
	}
)

func TestMain(m *testing.M) {
	cfg, err := config.Load("config_test.yaml")
	if err != nil {
		os.Exit(1)
	}
	db := sqllite.New(&cfg.DB)

	migrate := migrate.New(&cfg.DB)

	migrate.UP()

	vpnRepo = New(db)

	code := m.Run()
	os.Exit(code)
}

func TestCreateVPN(t *testing.T) {
	v1, err := vpnRepo.Create(vpnExample1)
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}
	if v1.ApiPort != vpnExample1.ApiPort ||
		v1.Domain != vpnExample1.Domain ||
		v1.StartRangePort != vpnExample1.StartRangePort ||
		v1.EndRangePort != vpnExample1.EndRangePort ||
		v1.VPNType != vpnExample1.VPNType {
		t.Error("the output wasn't matched by original data")
	}

	deleteVPNs(v1.ID)
}

func deleteVPNs(ids ...int) {
	for _, id := range ids {
		vpnRepo.Delete(id)
	}
}

func TestChangeStatus(t *testing.T) {
	v1, _ := vpnRepo.Create(vpnExample1)

	err := vpnRepo.activeVPN(v1.ID)
	if err != nil {
		t.Fatalf("the error has happend that was %v", err)
	}

	err = vpnRepo.deactiveVPN(v1.ID)
	if err != nil {
		t.Fatalf("the error has happend that was %v", err)
	}
}
