package hostmanager

import (
	"os"
	"testing"

	hostmanagerDto "momo/dto/repository/host_manager"
	"momo/entity"
	"momo/repository/migrate"
	"momo/repository/sqllite"
)

var hostRepo *Host

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	hostRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreateHost(t *testing.T) {
	host, err := hostRepo.Create(hostExample1)
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}
	if host.Domain != hostExample1.Domain ||
		host.Port != hostExample1.Port ||
		host.Status != entity.Deactive ||
		host.StartRangePort != hostExample1.StartRangePort ||
		host.EndRangePort != hostExample1.EndRangePort {
		t.Error("the out put of creating was wrong")
	}

	hostRepo.DeleteAll()
}

func TestFilterHosts(t *testing.T) {
	hostRepo.Create(hostExample2)
	hostRepo.Create(hostExample3)
	hostRepo.Create(hostExample4)
	hostRepo.Create(hostExample5)

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

	hostRepo.DeleteAll()
}

func TestUpdateHost(t *testing.T) {
	h1, _ := hostRepo.Create(hostExample6)
	err := hostRepo.Update(
		h1.ID,
		&hostmanagerDto.UpdateHost{Rank: 3, Status: entity.Low},
	)
	if err != nil {
		t.Errorf("error has happend that was %e", err)
	}
	hostRepo.DeleteAll()
}
