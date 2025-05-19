package vpnpackage

import (
	"os"
	"testing"

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
