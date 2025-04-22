package vpnmanager

import (
	"os"
	"testing"

	"momo/entity"
	"momo/pkg/config"
	"momo/repository/migrate"
	"momo/repository/sqllite"
	"momo/repository/sqllite/vpn_manager/dto"
)

var (
	vpnRepo *VPN

	vpnExample1 = &dto.Add_VPN{
		Domain:         "joi.com",
		ApiPort:        "62733",
		StartRangePort: 1000,
		EndRangePort:   2000,
		VPNType:        entity.XRAY_VPN,
		IsActive:       false,
	}

	vpnExample2 = &dto.Add_VPN{
		Domain:         "joi.com",
		ApiPort:        "62733",
		StartRangePort: 1000,
		EndRangePort:   2500,
		VPNType:        entity.XRAY_VPN,
		IsActive:       true,
	}

	vpnExample3 = &dto.Add_VPN{
		Domain:         "jordan.com",
		ApiPort:        "62733",
		StartRangePort: 3000,
		EndRangePort:   3500,
		VPNType:        entity.XRAY_VPN,
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
	deleteVPNs(v1.ID)
}

func TestFilterVPNs(t *testing.T) {
	v1, _ := vpnRepo.Create(vpnExample1)
	v2, _ := vpnRepo.Create(vpnExample2)
	v3, _ := vpnRepo.Create(vpnExample3)

	isActication := true
	vpns, err := vpnRepo.Filter(&dto.FilterVPNs{
		IsActive: &isActication,
	})
	if err != nil {
		t.Errorf("1. something wrong has happend that was %v", err)
	}

	if len(vpns) != 2 {
		t.Errorf("1. the number of vpns must be 2 but the result was %v", len(vpns))
	}

	vpns, err = vpnRepo.Filter(&dto.FilterVPNs{
		Domain: "joi.com",
	})
	if err != nil {
		t.Errorf("2. something wrong has happend that was %v", err)
	}

	if len(vpns) != 2 {
		t.Errorf("2. the number of vpns must be 2 but the result was %v", len(vpns))
	}

	vpns, err = vpnRepo.Filter(&dto.FilterVPNs{
		VPNType: entity.XRAY_VPN,
	})
	if err != nil {
		t.Errorf("3. something wrong has happend that was %v", err)
	}

	if len(vpns) != 3 {
		t.Errorf("3. the number of vpns must be 3 but the result was %v", len(vpns))
	}

	vpns, err = vpnRepo.Filter(&dto.FilterVPNs{
		VPNType: entity.XRAY_VPN,
		Domain:  "joi.com",
	})
	if err != nil {
		t.Errorf("4. something wrong has happend that was %v", err)
	}

	if len(vpns) != 2 {
		t.Errorf("4. the number of vpns must be 2 but the result was %v", len(vpns))
	}

	deleteVPNs(v1.ID, v2.ID, v3.ID)
}

func deleteVPNs(ids ...int) {
	for _, id := range ids {
		vpnRepo.Delete(id)
	}
}
