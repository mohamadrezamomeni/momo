package tier

import (
	"os"
	"testing"

	tieruser "github.com/mohamadrezamomeni/momo/dto/repository/tier"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

var tierRepo *Tier

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-tier.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	tierRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreateInbound(t *testing.T) {
	defer tierRepo.DeleteAll()
	tier1, err := tierRepo.Create(tieruser.CreateTier{
		Name:    "silver",
		Default: true,
	})
	if err != nil {
		t.Errorf("something wrong has happend the problem was %v", err)
	}
	if tier1.Name != "silver" ||
		tier1.IsDefault != true {
		t.Error("data wasn't saved currectly")
	}
}

func TestFindingTierByID(t *testing.T) {
	defer tierRepo.DeleteAll()
	tierRepo.Create(tieruser.CreateTier{
		Name:    "silver",
		Default: true,
	})
	tier, err := tierRepo.FindByName("silver")
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	if tier.Name != "silver" {
		t.Fatal("error to compare data")
	}
}

func TestFilterTiers(t *testing.T) {
	defer tierRepo.DeleteAll()
	tierRepo.Create(tieruser.CreateTier{
		Name:    "silver",
		Default: true,
	})
	tierRepo.Create(tieruser.CreateTier{
		Name:    "gold",
		Default: true,
	})
	tiers, err := tierRepo.Filter()
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(tiers) != 2 {
		t.Fatal("error to compare data")
	}
}

func TestUpdateTier(t *testing.T) {
	defer tierRepo.DeleteAll()
	tierRepo.Create(tieruser.CreateTier{
		Name:    "silver",
		Default: true,
	})

	isDefault := false
	err := tierRepo.Update("silver", &tieruser.Update{
		IsDefault: &isDefault,
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	tier, err := tierRepo.FindByName("silver")
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	if tier.IsDefault != false {
		t.Fatal("error to compare data")
	}
}
