package hostmanager

import (
	"os"
	"testing"

	hostmanagerDto "momo/dto/repository/host_manager"
	"momo/entity"
	"momo/pkg/config"
	"momo/repository/migrate"
	"momo/repository/sqllite"
)

var hostRepo *Host

func TestMain(m *testing.M) {
	cfg, err := config.Load("config_test.yaml")
	if err != nil {
		os.Exit(1)
	}
	db := sqllite.New(&cfg.DB)

	migrate := migrate.New(&cfg.DB)

	migrate.UP()

	hostRepo = New(db)

	code := m.Run()
	os.Exit(code)
}

var hostExample1 = &hostmanagerDto.AddHost{
	Domain: "google.com",
	Port:   "62789",
}

func TestCreateHost(t *testing.T) {
	host, err := hostRepo.Create(hostExample1)
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}

	if host.Domain != hostExample1.Domain ||
		host.Port != hostExample1.Port ||
		host.Status != entity.Deactive {
		t.Error("the out put of creating was wrong")
	}

	hostRepo.Delete(host.ID)
}
