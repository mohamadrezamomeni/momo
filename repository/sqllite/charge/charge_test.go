package charge

import (
	"os"
	"testing"

	chargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/charge"
	"github.com/mohamadrezamomeni/momo/entity"
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
		chargeCreated.UserID != charge1.UserID ||
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

func TestUpdateCharge(t *testing.T) {
	defer chargeRepo.DeleteAll()
	chargeCreated, _ := chargeRepo.Create(charge1)

	newDetail := "goodBye"
	adminComment := "sorry"
	err := chargeRepo.UpdateCharge(chargeCreated.ID, &chargeRepositoryDto.UpdateChargeDto{
		Detail:       newDetail,
		AdminComment: adminComment,
		Status:       entity.RegejectedStatusCharge,
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	chargeFound, _ := chargeRepo.FindChargeByID(chargeCreated.ID)
	if chargeFound.AdminComment != adminComment ||
		chargeFound.Detail != newDetail ||
		chargeFound.Status != entity.RegejectedStatusCharge {
		t.Fatal("error to comapre data")
	}
}
