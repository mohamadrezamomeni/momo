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

func TestFilterCharges(t *testing.T) {
	defer chargeRepo.DeleteAll()
	chargeRepo.Create(charge1)
	chargeRepo.Create(charge2)
	chargeRepo.Create(charge3)

	charges, err := chargeRepo.FilterCharges(&chargeRepositoryDto.FilterChargesDto{})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(charges) != 3 {
		t.Fatalf("we expected the lengh of charges be 3 but we got %v", len(charges))
	}

	charges, err = chargeRepo.FilterCharges(&chargeRepositoryDto.FilterChargesDto{
		Status: entity.RegejectedStatusCharge,
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(charges) != 0 {
		t.Fatalf("we expected the lengh of charges be 0 but we got %v", len(charges))
	}

	charges, err = chargeRepo.FilterCharges(&chargeRepositoryDto.FilterChargesDto{
		UserID: "f47ac10b-58cc-4372-a567-0e02b2c3d477",
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(charges) != 2 {
		t.Fatalf("we expected the lengh of charges be 2 but we got %v", len(charges))
	}

	charges, err = chargeRepo.FilterCharges(&chargeRepositoryDto.FilterChargesDto{
		InboundID: "12",
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(charges) != 2 {
		t.Fatalf("we expected the lengh of charges be 2 but we got %v", len(charges))
	}

	charges, err = chargeRepo.FilterCharges(&chargeRepositoryDto.FilterChargesDto{
		InboundID: "12",
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d477",
		Status:    entity.PendingStatusCharge,
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(charges) != 1 {
		t.Fatalf("we expected the lengh of charges be 1 but we got %v", len(charges))
	}
}

func TestGetFirstAvailbleInboundCharge(t *testing.T) {
	defer chargeRepo.DeleteAll()

	chargeCreated, _ := chargeRepo.Create(charge4)
	chargeRepo.Create(charge5)

	chargeFound, err := chargeRepo.GetFirstAvailbleInboundCharge(chargeCreated.InboundID)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if chargeFound.ID != chargeCreated.ID {
		t.Fatal("error to compare data")
	}
}

func TestRetriveAvailbleChargesForInbounds(t *testing.T) {
	defer chargeRepo.DeleteAll()

	chargeCreated1, _ := chargeRepo.Create(charge6)
	chargeRepo.Create(charge7)
	chargeRepo.Create(charge8)
	chargeCreated4, _ := chargeRepo.Create(charge9)

	charges, err := chargeRepo.RetriveAvailbleChargesForInbounds([]string{"15", "16", "17"})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(charges) != 2 {
		t.Fatalf("we expected the lengh of charges be 2 but we got %d", len(charges))
	}

	if !(chargeCreated1.ID == charges[0].ID || chargeCreated1.ID == charges[1].ID) {
		t.Fatal("error to compare data")
	}

	if !(chargeCreated4.ID == charges[0].ID || chargeCreated4.ID == charges[1].ID) {
		t.Fatal("error to compare data")
	}
}
