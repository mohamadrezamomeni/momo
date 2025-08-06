package inboundcharge

import (
	"os"
	"testing"
	"time"

	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"

	chargeRepository "github.com/mohamadrezamomeni/momo/repository/sqllite/charge"
	inboundRepository "github.com/mohamadrezamomeni/momo/repository/sqllite/inbound"
	vpnPackageRepository "github.com/mohamadrezamomeni/momo/repository/sqllite/vpn_package"
)

var (
	inboundChargeRepo *InboundCharge
	inboundRepo       *inboundRepository.Inbound
	chargeRepo        *chargeRepository.Charge
	vpnPackageRepo    *vpnPackageRepository.VPNPackage
)

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-inbound-charge.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	inboundChargeRepo = New(db)
	inboundRepo = inboundRepository.New(db)
	chargeRepo = chargeRepository.New(db)
	vpnPackageRepo = vpnPackageRepository.New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestChargeInbound(t *testing.T) {
	defer inboundRepo.DeleteAll()
	defer vpnPackageRepo.DeleteAll()
	defer chargeRepo.DeleteAll()
	inbound, _ := inboundRepo.Create(inbound1)
	vpnPackage, _ := vpnPackageRepo.Create(vpnPackage1)
	charge1.InboundID = inbound.ID
	charge1.PackageID = vpnPackage.ID
	charge, _ := chargeRepo.Create(charge1)
	err := inboundChargeRepo.AssignChargeToInbound(inbound, charge, vpnPackage)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	chargeFound, _ := chargeRepo.FindChargeByID(charge.ID)
	inboundFound, _ := inboundRepo.FindInboundByID(charge.ID)

	if chargeFound.Status != entity.AssignedCharged {
		t.Fatalf("error to compare data charge status was wrong")
	}

	yn, mn, dn := time.Now().Date()
	yp, mp, dp := time.Now().AddDate(0, int(vpnPackage.Months), int(vpnPackage.Days)).Date()
	ys, ms, ds := inboundFound.Start.Date()
	ye, me, de := inboundFound.End.Date()

	if !((yp == ye && mp == me && de == dp) && (yn == ys && mn == ms && dn == ds)) {
		t.Fatal("error to comapre data")
	}
}

func TestCreateInbound(t *testing.T) {
	vpnPackage, _ := vpnPackageRepo.Create(vpnPackage1)
	charge1.PackageID = vpnPackage.ID
	charge, _ := chargeRepo.Create(charge1)
	err := inboundChargeRepo.CreateInbound(charge.ID, inbound2)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	chargeFound, _ := chargeRepo.FindChargeByID(charge.ID)
	if chargeFound.Status != entity.AssignedCharged {
		t.Fatal("error to compare data")
	}
}
