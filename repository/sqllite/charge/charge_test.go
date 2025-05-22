package charge

import (
	"os"
	"testing"

	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

var chargeRepo *Charge

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-charge.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	chargeRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCharge(t *testing.T) {
	defer chargeRepo.DeleteAll()
	chargeCreated, err := chargeRepo.Create(charge1)
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	if chargeCreated.AdminComment != "" ||
		chargeCreated.Detail != charge1.Detail ||
		chargeCreated.InboundID != charge1.InboundID ||
		chargeCreated.Status != charge1.Status ||
		chargeCreated.ID == "" {
		t.Fatal("error to comapre data")
	}
}

func TestFindChargeByID(t *testing.T) {
	defer chargeRepo.DeleteAll()
	chargeCreated, _ := chargeRepo.Create(charge1)

	chargeFound, err := chargeRepo.FindChargeByID(chargeCreated.ID)
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	if chargeFound.ID != chargeCreated.ID {
		t.Fatalf("error to compare data")
	}
}
