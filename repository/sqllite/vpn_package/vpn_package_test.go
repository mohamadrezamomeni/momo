package vpnpackage

import (
	"os"
	"testing"

	vpnpackage "github.com/mohamadrezamomeni/momo/dto/repository/vpn_package"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

var vpnPackageRepo *VPNPackage

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-vpn-package.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	vpnPackageRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreateVPNPackage(t *testing.T) {
	vpnPackageCreated, err := vpnPackageRepo.Create(vpnPackage1)
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}
	if vpnPackageCreated.Days != vpnPackage1.Days ||
		vpnPackageCreated.Months != vpnPackage1.Months ||
		vpnPackageCreated.IsActive != vpnPackage1.IsActive ||
		vpnPackageCreated.Price != vpnPackage1.Price ||
		vpnPackageCreated.PriceTitle != vpnPackage1.PriceTitle ||
		vpnPackageCreated.TrafficLimit != vpnPackage1.TrafficLimit ||
		vpnPackageCreated.TrafficLimitTitle != vpnPackage1.TrafficLimitTitle {
		t.Error("we got unexpected result")
	}
}

func TestFindVPNPackage(t *testing.T) {
	vpnCreated, _ := vpnPackageRepo.Create(vpnPackage1)

	vpnPackageFound, err := vpnPackageRepo.FindVPNPackageByID(vpnCreated.ID)
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	if vpnPackageFound.Days != vpnPackage1.Days ||
		vpnPackageFound.Months != vpnPackage1.Months ||
		vpnPackageFound.IsActive != vpnPackage1.IsActive ||
		vpnPackageFound.Price != vpnPackage1.Price ||
		vpnPackageFound.PriceTitle != vpnPackage1.PriceTitle ||
		vpnPackageFound.TrafficLimit != vpnPackage1.TrafficLimit ||
		vpnPackageFound.TrafficLimitTitle != vpnPackage1.TrafficLimitTitle {
		t.Error("we got unexpected result")
	}
}

func TestUpdateVPNPackage(t *testing.T) {
	vpnCreated, _ := vpnPackageRepo.Create(vpnPackage1)
	deactive := false
	err := vpnPackageRepo.Update(vpnCreated.ID, &vpnpackage.UpdateVPNPackage{
		IsActive: &deactive,
	})
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	vpnPackageFound, err := vpnPackageRepo.FindVPNPackageByID(vpnCreated.ID)

	if vpnPackageFound.IsActive != false {
		t.Fatalf("we expected the vpn was deactive")
	}
}
