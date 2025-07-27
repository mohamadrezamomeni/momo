package vpnmanager

import (
	"os"
	"testing"

	vpnManagerDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_manager"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-vpn.db",
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

	err := vpnRepo.ActiveVPN(v1.ID)
	if err != nil {
		t.Fatalf("the error has happend that was %v", err)
	}

	err = vpnRepo.DeactiveVPN(v1.ID)
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
		VPNTypes: []entity.VPNType{entity.XRAY_VPN},
	})
	if err != nil {
		t.Errorf("3. something wrong has happend that was %v", err)
	}

	if len(vpns) != 3 {
		t.Errorf("3. the number of vpns must be 3 but the result was %v", len(vpns))
	}

	vpns, err = vpnRepo.Filter(&vpnManagerDto.FilterVPNs{
		VPNTypes: []entity.VPNType{entity.XRAY_VPN},
		Domain:   "joi.com",
	})
	if err != nil {
		t.Errorf("4. something wrong has happend that was %v", err)
	}

	if len(vpns) != 2 {
		t.Errorf("4. the number of vpns must be 2 but the result was %v", len(vpns))
	}

	vpns, err = vpnRepo.Filter(&vpnManagerDto.FilterVPNs{
		Coountries: []string{"uk"},
	})
	if err != nil {
		t.Errorf("5. something wrong has happend that was %v", err)
	}

	if len(vpns) != 2 {
		t.Errorf("5. the number of vpns must be 2 but the result was %v", len(vpns))
	}
}

func TestGroupingByCountry(t *testing.T) {
	defer vpnRepo.DeleteAll()
	vpnRepo.Create(vpn1)
	vpnRepo.Create(vpn2)
	vpnRepo.Create(vpn3)
	vpnRepo.Create(vpn4)

	countries, err := vpnRepo.GroupAvailbleVPNsByCountry()
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(countries) != 2 {
		t.Fatalf("we expected the lengh of coutnreis be 2 but we got %d", len(countries))
	}
	countriesRefrence := map[string]struct{}{}

	for _, country := range countries {
		countriesRefrence[country] = struct{}{}
	}
	for _, country := range []string{"china", "uk"} {
		if _, isExist := countriesRefrence[country]; !isExist {
			t.Fatalf("we expected china be existed but the response miss that`")
		}
	}
}

func ValidateResponseGroupDomainByVPNSource(
	t *testing.T,
	validResponse map[string][]string,
	VPNSourceDomains map[string][]string,
) {
	for country, domains := range validResponse {
		if _, isExist := VPNSourceDomains[country]; !isExist {
			t.Fatalf("the country %s was missed", country)
		}

		domainsRefrence := map[string]struct{}{}
		for _, domain := range VPNSourceDomains[country] {
			domainsRefrence[domain] = struct{}{}
		}
		for _, domain := range domains {
			if _, isExist := domainsRefrence[domain]; !isExist {
				t.Fatalf("the domain %s was missed", domain)
			}
		}
	}
}
