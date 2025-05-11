package inbound

import (
	"os"
	"strconv"
	"testing"

	inboundDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

var inboundRepo *Inbound

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-inbound.db",
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

	defer inboundRepo.DeleteAll()
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

	inbounds, err = inboundRepo.Filter(&inboundDto.FilterInbound{Port: port2})
	if err != nil {
		t.Errorf("1. the problem has occured that is %v", err)
	}
}

func TestRertriveFaultyInbounds(t *testing.T) {
	inboundRepo.Create(inbound4)
	inboundRepo.Create(inbound5)
	inboundRepo.Create(inbound6)
	inboundRepo.Create(inbound7)
	inboundRepo.Create(inbound15)

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

func TestInboundsIsNotAssigned(t *testing.T) {
	inboundRepo.Create(inbound10)
	inboundRepo.Create(inbound11)

	inbounds, err := inboundRepo.FindInboundIsNotAssigned()
	if err != nil {
		t.Fatalf("someting went wront the problem was %v", err)
	}

	if len(inbounds) != 1 {
		t.Fatalf("we expected we got 1 items but we got %v", len(inbounds))
	}
	if inbounds[0].UserID != userID8 {
		t.Fatalf("the answer was wrong")
	}
	inboundRepo.DeleteAll()
}

func TestFindInboundByUserID(t *testing.T) {
	inboundCreated, _ := inboundRepo.Create(inbound10)
	defer inboundRepo.DeleteAll()
	inbound, err := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated.ID))
	if err != nil {
		t.Fatalf("the query was field the problem was %v", err)
	}
	if inbound.UserID != inbound10.UserID || inboundCreated.ID != inbound.ID {
		t.Fatalf("the query has answerd wrong")
	}
}

func TestUpdateInbound(t *testing.T) {
	inboundCreated, _ := inboundRepo.Create(inbound10)
	defer inboundRepo.DeleteAll()
	newDomain := "facebook.com"
	newPort := "2020"
	err := inboundRepo.UpdateDomainPort(inboundCreated.ID, newDomain, newPort)
	if err != nil {
		t.Fatalf("update has field the problem was %v", err)
	}

	inbound, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated.ID))

	if inbound.Domain != newDomain || inbound.Port != newPort {
		t.Fatalf("update hasn't worked carefuly")
	}
}

func TestGetListOfPortsByDomain(t *testing.T) {
	inboundRepo.Create(inbound12)
	inboundRepo.Create(inbound13)
	inboundRepo.Create(inbound14)

	defer inboundRepo.DeleteAll()

	summery, err := inboundRepo.GetListOfPortsByDomain()
	if err != nil {
		t.Fatalf("the problem has happend that was %v", err)
	}
	if len(summery) != 2 {
		t.Fatalf("result was wrong we expected the number of items are 2 but we got %d", len(summery))
	}

	mapSummery := map[string][]string{}
	for _, item := range summery {
		mapSummery[item.Domain] = item.Ports
	}

	if ports, ok := mapSummery["twitter.com"]; !ok || len(ports) != 2 {
		t.Fatalf("output was wrong.")
	}

	if ports, ok := mapSummery["google.com"]; !ok || len(ports) != 1 {
		t.Fatalf("output was wrong.")
	}
}

func TestBlock(t *testing.T) {
	inboundCreated, _ := inboundRepo.Create(inbound1)
	defer inboundRepo.DeleteAll()

	err := inboundRepo.ChangeBlockState(strconv.Itoa(inboundCreated.ID), true)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	inboundFound, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated.ID))

	if !inboundFound.IsBlock {
		t.Fatal("inbound that is founded must be blocked")
	}
}

func TestUnBlock(t *testing.T) {
	inboundCreated, _ := inboundRepo.Create(inbound5)
	defer inboundRepo.DeleteAll()
	err := inboundRepo.ChangeBlockState(strconv.Itoa(inboundCreated.ID), false)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	inboundFound, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated.ID))

	if inboundFound.IsBlock {
		t.Fatal("inbound that is founded must be blocked")
	}
}

func TestUpdate(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundCreated, _ := inboundRepo.Create(inbound1)

	err := inboundRepo.Update(strconv.Itoa(inboundCreated.ID), &inboundDto.UpdateInboundDto{
		Start: utils.GetDateTime("2026-04-21 14:30:00"),
		End:   utils.GetDateTime("2026-04-22 14:30:00"),
	})
	if err != nil {
		t.Fatalf("the error was %v", err)
	}

	inboundFound, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated.ID))

	if inboundFound.Start != utils.GetDateTime("2026-04-21 14:30:00") {
		t.Errorf("we expected start be %s", inboundFound.Start.Format("2006-01-02 15:04:05"))
	}

	if inboundFound.End != utils.GetDateTime("2026-04-22 14:30:00") {
		t.Errorf("we expected start be %s", inboundFound.End.Format("2006-01-02 15:04:05"))
	}

	inboundCreated, _ = inboundRepo.Create(inbound2)

	err = inboundRepo.Update(strconv.Itoa(inboundCreated.ID), &inboundDto.UpdateInboundDto{
		Start: utils.GetDateTime("2026-04-21 14:30:00"),
	})
	if err != nil {
		t.Fatalf("the error was %v", err)
	}

	inboundFound, _ = inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated.ID))

	if inboundFound.Start != utils.GetDateTime("2026-04-21 14:30:00") {
		t.Errorf("we expected start be %s", inboundFound.Start.Format("2006-01-02 15:04:05"))
	}
}
