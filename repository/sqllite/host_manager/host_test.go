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

var (
	hostExample1 = &hostmanagerDto.AddHost{
		Domain: "google.com",
		Port:   "62789",
		Status: entity.Deactive,
	}
	hostExample2 = &hostmanagerDto.AddHost{
		Domain: "yahoo.com",
		Port:   "62780",
		Status: entity.High,
	}
	hostExample3 = &hostmanagerDto.AddHost{
		Domain: "facebook.com",
		Port:   "62780",
		Status: entity.Deactive,
	}
	hostExample4 = &hostmanagerDto.AddHost{
		Domain: "twitter.com",
		Port:   "62780",
		Status: entity.Medium,
	}
	hostExample5 = &hostmanagerDto.AddHost{
		Domain: "github.com",
		Port:   "62780",
		Status: entity.Low,
	}
)

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

func TestFilterHosts(t *testing.T) {
	h1, _ := hostRepo.Create(hostExample2)
	h2, _ := hostRepo.Create(hostExample3)
	h3, _ := hostRepo.Create(hostExample4)
	h4, _ := hostRepo.Create(hostExample5)

	hosts, err := hostRepo.Filter(&hostmanagerDto.FilterHosts{})
	if err != nil {
		t.Errorf("1.something wrong has happend that was %v", err)
	}

	if len(hosts) != 4 {
		t.Errorf("1. we expected 4 items but we got %v", len(hosts))
	}

	hosts, err = hostRepo.Filter(&hostmanagerDto.FilterHosts{
		Statuses: []entity.HostStatus{entity.High},
	})
	if err != nil {
		t.Errorf("2. something wrong has happend that was %v", err)
	}

	if len(hosts) != 1 {
		t.Errorf("2. we expected 1 items but we got %v", len(hosts))
	}

	if hosts[0].Status != entity.High {
		t.Errorf("2. we expected status be high")
	}

	hosts, err = hostRepo.Filter(&hostmanagerDto.FilterHosts{
		Statuses: []entity.HostStatus{entity.Medium},
	})
	if err != nil {
		t.Errorf("3. something wrong has happend that was %v", err)
	}

	if len(hosts) != 1 {
		t.Errorf("3. we expected 1 items but we got %v", len(hosts))
	}
	if hosts[0].Status != entity.Medium {
		t.Errorf("3. we expected status be high")
	}

	hosts, err = hostRepo.Filter(&hostmanagerDto.FilterHosts{
		Statuses: []entity.HostStatus{entity.Low},
	})
	if err != nil {
		t.Errorf("4. something wrong has happend that was %v", err)
	}

	if len(hosts) != 1 {
		t.Errorf("4. we expected 1 items but we got %v", len(hosts))
	}
	if hosts[0].Status != entity.Low {
		t.Errorf("3. we expected status be low")
	}

	hosts, err = hostRepo.Filter(&hostmanagerDto.FilterHosts{
		Statuses: []entity.HostStatus{entity.High, entity.Medium, entity.Low},
	})
	if err != nil {
		t.Errorf("something wrong has happend that was %v", err)
	}

	if len(hosts) != 3 {
		t.Errorf("we expected 3 items but we got %v", len(hosts))
	}

	deleteHosts(h1.ID, h2.ID, h3.ID, h4.ID)
}

func deleteHosts(ids ...int) {
	for _, id := range ids {
		hostRepo.Delete(id)
	}
}
