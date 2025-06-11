package vpnsource

import (
	"os"
	"testing"

	vpnsource "github.com/mohamadrezamomeni/momo/dto/repository/vpn_source"
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
	defer VPNSourceRepo.DeleteAll()
	vpnsource, err := VPNSourceRepo.Create(vpnsource1)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if vpnsource.Country != vpnsource.Country ||
		vpnsource.English != vpnsource1.English {
		t.Fatal("error to compare data")
	}
}

func TestFindVPNSource(t *testing.T) {
	defer VPNSourceRepo.DeleteAll()
	vpnsourceCreated, _ := VPNSourceRepo.Create(vpnsource1)

	vpnsource, err := VPNSourceRepo.Find(vpnsourceCreated.Country)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	if vpnsource.Country != vpnsourceCreated.Country {
		t.Fatal("error to compare data")
	}
}

func TestUpdateVPNSource(t *testing.T) {
	defer VPNSourceRepo.DeleteAll()
	vpnsourceCreated, _ := VPNSourceRepo.Create(vpnsource1)

	newEnglishTranslation := "united-state"
	err := VPNSourceRepo.Update(vpnsourceCreated.Country, &vpnsource.UpdateVPNSourceDto{
		English: newEnglishTranslation,
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	vpnSource, _ := VPNSourceRepo.Find(vpnsourceCreated.Country)
	if vpnSource.English != newEnglishTranslation {
		t.Fatalf("error to compare data")
	}
}

func TestFilterVPNSource(t *testing.T) {
	vpnsourceCreated1, _ := VPNSourceRepo.Create(vpnsource1)
	vpnsourceCreated2, _ := VPNSourceRepo.Create(vpnsource2)
	vpnsourceCreated3, _ := VPNSourceRepo.Create(vpnsource3)
	vpnsourcesRefrences := map[string]struct{}{}

	vpnsourcesRefrences[vpnsourceCreated1.Country] = struct{}{}
	vpnsourcesRefrences[vpnsourceCreated2.Country] = struct{}{}
	vpnsourcesRefrences[vpnsourceCreated3.Country] = struct{}{}

	vpnsources, err := VPNSourceRepo.Filter(&vpnsource.FilterVPNSources{})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	if len(vpnsources) != 3 {
		t.Fatalf("we expected the lengh of vpnsource be 3 but we got %d", len(vpnsources))
	}

	for _, vpnsource := range vpnsources {
		if _, isExist := vpnsourcesRefrences[vpnsource.Country]; !isExist {
			t.Fatalf("miss the vpnsource with country %s", vpnsource.Country)
		}
	}

	vpnsources, err = VPNSourceRepo.Filter(&vpnsource.FilterVPNSources{
		Countries: []string{vpnsourceCreated1.Country, vpnsourceCreated2.Country},
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	if len(vpnsources) != 2 {
		t.Fatalf("we expected the lengh of vpnsource be 3 but we got %d", len(vpnsources))
	}

	for _, vpnsource := range vpnsources {
		if !(vpnsource.Country == vpnsourceCreated1.Country || vpnsourceCreated2.Country == vpnsource.Country) {
			t.Fatalf("we get unexpected vpnsource with %s country", vpnsource.Country)
		}
	}
}
