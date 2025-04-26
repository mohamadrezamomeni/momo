package vpnmanager

import (
	"os"
	"testing"

	vpnManagerDto "momo/dto/repository/vpn_manager"
	"momo/entity"
	"momo/repository/migrate"
	"momo/repository/sqllite"
)

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	vpnRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreateVPN(t *testing.T) {
	v1, err := vpnRepo.Create(vpn1)
	defer vpnRepo.DeleteAll()
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}
	if v1.ApiPort != vpn1.ApiPort ||
		v1.Domain != vpn1.Domain ||
		v1.VPNType != vpn1.VPNType {
		t.Error("the output wasn't matched by original data")
	}
}

func TestChangeStatus(t *testing.T) {
	v1, _ := vpnRepo.Create(vpn1)
	defer vpnRepo.DeleteAll()

	err := vpnRepo.activeVPN(v1.ID)
	if err != nil {
		t.Fatalf("the error has happend that was %v", err)
	}

	err = vpnRepo.deactiveVPN(v1.ID)
	if err != nil {
		t.Fatalf("the error has happend that was %v", err)
	}
}

func TestFilterVPNs(t *testing.T) {
	defer vpnRepo.DeleteAll()
	vpnRepo.Create(vpn1)
	vpnRepo.Create(vpn2)
	vpnRepo.Create(vpn3)
	isActication := true
	vpns, err := vpnRepo.Filter(&vpnManagerDto.FilterVPNs{
		IsActive: &isActication,
	})
	if err != nil {
		t.Errorf("1. something wrong has happend that was %v", err)
	}

	if len(vpns) != 2 {
		t.Errorf("1. the number of vpns must be 2 but the result was %v", len(vpns))
	}

	vpns, err = vpnRepo.Filter(&vpnManagerDto.FilterVPNs{
		Domain: "joi.com",
	})
	if err != nil {
		t.Errorf("2. something wrong has happend that was %v", err)
	}

	if len(vpns) != 2 {
		t.Errorf("2. the number of vpns must be 2 but the result was %v", len(vpns))
	}

	vpns, err = vpnRepo.Filter(&vpnManagerDto.FilterVPNs{
		VPNType: entity.XRAY_VPN,
	})
	if err != nil {
		t.Errorf("3. something wrong has happend that was %v", err)
	}

	if len(vpns) != 3 {
		t.Errorf("3. the number of vpns must be 3 but the result was %v", len(vpns))
	}

	vpns, err = vpnRepo.Filter(&vpnManagerDto.FilterVPNs{
		VPNType: entity.XRAY_VPN,
		Domain:  "joi.com",
	})
	if err != nil {
		t.Errorf("4. something wrong has happend that was %v", err)
	}

	if len(vpns) != 2 {
		t.Errorf("4. the number of vpns must be 2 but the result was %v", len(vpns))
	}
}
