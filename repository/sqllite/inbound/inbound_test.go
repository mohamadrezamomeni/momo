package inbound

import (
	"os"
	"testing"

	inboundDto "momo/dto/repository/inbound"
	"momo/repository/migrate"
	"momo/repository/sqllite"
)

var inboundRepo *Inbound

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	inboundRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreateInbound(t *testing.T) {
	ret, err := inboundRepo.Create(inbound1)
	if err != nil {
		t.Errorf("something wrong has happend the problem was %v", err)
	}
	if ret.Domain != inbound1.Domain ||
		ret.Port != inbound1.Port ||
		ret.IsActive != false ||
		ret.IsBlock != inbound1.IsBlock ||
		ret.UserID != inbound1.UserID ||
		ret.VPNType != inbound1.VPNType {
		t.Error("data wasn't saved currectly")
	}

	inboundRepo.DeleteAll()
}

func TestChangeStatus(t *testing.T) {
	ret, err := inboundRepo.Create(inbound1)
	if err != nil {
		t.Errorf("something wrong has happend the problem was %v", err)
	}

	err = inboundRepo.changeStatus(ret.ID, false)
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}
	inboundRepo.DeleteAll()
}

func TestFilterInbounds(t *testing.T) {
	inboundRepo.Create(inbound1)
	inboundRepo.Create(inbound2)
	inboundRepo.Create(inbound3)
	inboundRepo.Create(inbound4)

	inbounds, err := inboundRepo.Filter(&inboundDto.FilterInbound{Port: port2})
	if err != nil {
		t.Errorf("1. the problem has occured that is %v", err)
	}
	if len(inbounds) != 2 {
		t.Errorf("1. the number of items must be %v but got %v items", 2, len(inbounds))
	}

	inbounds, err = inboundRepo.Filter(&inboundDto.FilterInbound{Protocol: "vmess"})
	if err != nil {
		t.Errorf("2. the problem has occured that is %v", err)
	}
	if len(inbounds) != 2 {
		t.Errorf("2. the number of items must be %v but got %v items", 2, len(inbounds))
	}

	inbounds, err = inboundRepo.Filter(&inboundDto.FilterInbound{Domain: "google.com"})
	if err != nil {
		t.Errorf("3. the problem has occured that is %v", err)
	}
	if len(inbounds) != 1 {
		t.Errorf("3. the number of items must be %v but got %v items", 1, len(inbounds))
	}

	inbounds, err = inboundRepo.Filter(&inboundDto.FilterInbound{UserID: userID2, Port: port2})
	if err != nil {
		t.Errorf("4. the problem has occured that is %v", err)
	}
	if len(inbounds) != 2 {
		t.Errorf("4. the number of items must be %v but got %v items", 2, len(inbounds))
	}

	isAvailableT := true
	inbounds, err = inboundRepo.Filter(&inboundDto.FilterInbound{IsActive: &isAvailableT})
	if err != nil {
		t.Errorf("5. the problem has occured that is %v", err)
	}
	if len(inbounds) != 1 {
		t.Errorf("5. the number of items must be %v but got %v items", 1, len(inbounds))
	}

	inboundRepo.DeleteAll()
}

func TestRertriveFaultyInbounds(t *testing.T) {
	inboundRepo.Create(inbound4)
	inboundRepo.Create(inbound5)
	inboundRepo.Create(inbound6)
	inboundRepo.Create(inbound7)

	inbounds, err := inboundRepo.RetriveFaultyInbounds()
	if err != nil {
		t.Errorf("the problem has happend that was %v", err)
	}
	if len(inbounds) != 2 {
		t.Errorf("the number of inbouns could be 2 but system got %v", len(inbounds))
	}
	userID4Status := false
	userID6Status := false
	for _, inbound := range inbounds {
		switch inbound.UserID {
		case userID4:
			userID4Status = true
		case userID6:
			userID6Status = true
		default:
			t.Errorf("un expected vpn with userID %v", inbound.UserID)
		}
	}
	if !userID4Status {
		t.Error("we didn't get the userID4")
	}
	if !userID6Status {
		t.Error("we didn't get the userID6")
	}

	inboundRepo.DeleteAll()
}

func TestCouningUsedDataEachDomian(t *testing.T) {
	inboundRepo.Create(inbound3)
	inboundRepo.Create(inbound4)
	inboundRepo.Create(inbound5)
	inboundRepo.Create(inbound6)
	inboundRepo.Create(inbound7)
	inboundRepo.Create(inbound8)
	inboundRepo.Create(inbound9)

	sumery, err := inboundRepo.CountingUsedPortEachHost()
	if err != nil {
		t.Error(err.Error())
	}
	if len(sumery) != 2 {
		t.Errorf("we expeted the len of mapSumery be 2 bug we got %v", len(sumery))
	}

	mapSumery := map[string]uint16{}

	for _, data := range sumery {
		mapSumery[data.domain] = data.count
	}

	if count, ok := mapSumery["twitter.com"]; ok && count != 1 {
		t.Errorf("the count of twitter.com must be 1 but we got %v", count)
	}

	if count, ok := mapSumery["googoo.com"]; ok && count != 3 {
		t.Errorf("the count of googoo.com must be 3 but we got %v", count)
	}

	inboundRepo.DeleteAll()
}
