package vpnsource

import (
	"os"
	"testing"

	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

var VPNSourceRepo *VPNSource

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-vpn-source.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	VPNSourceRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreateingVPNSource(t *testing.T) {
	vpnsource, err := VPNSourceRepo.Create(vpnsource1)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if vpnsource.Title != vpnsource.Title ||
		vpnsource.English != vpnsource1.English ||
		vpnsource.ID == "" {
		t.Fatal("error to compare data")
	}
}

func TestFindVPNSource(t *testing.T) {
	vpnsourceCreated, _ := VPNSourceRepo.Create(vpnsource1)

	vpnsource, err := VPNSourceRepo.Find(vpnsourceCreated.ID)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	if vpnsource.ID != vpnsourceCreated.ID {
		t.Fatal("error to compare data")
	}
}
