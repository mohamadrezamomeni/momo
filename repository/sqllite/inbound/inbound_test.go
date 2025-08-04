package inbound

import (
	"os"
	"testing"
	"time"

	inboundDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
	timeTransformer "github.com/mohamadrezamomeni/momo/transformer/time"
)

var inboundRepo *Inbound

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-inbbound.db",
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
	defer inboundRepo.DeleteAll()
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
}

func TestChangeStatus(t *testing.T) {
	defer inboundRepo.DeleteAll()
	ret, err := inboundRepo.Create(inbound1)
	if err != nil {
		t.Errorf("something wrong has happend the problem was %v", err)
	}

	err = inboundRepo.changeStatus(ret.ID, false)
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}
}

func TestFilterInbounds(t *testing.T) {
	defer inboundRepo.DeleteAll()
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

	inbounds, err = inboundRepo.Filter(&inboundDto.FilterInbound{Port: port2})
	if err != nil {
		t.Errorf("1. the problem has occured that is %v", err)
	}
}

func TestRetriveActiveInboundBlocked(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundRepo.Create(inbound4)
	inboundRepo.Create(inbound5)
	inboundRepo.Create(inbound6)
	inboundRepo.Create(inbound7)
	inboundRepo.Create(inbound15)
	inboundRepo.Create(inbound16)

	inbounds, err := inboundRepo.RetriveActiveInboundBlocked()
	if err != nil {
		t.Errorf("something went wrong the problem was %v", err)
	}

	if len(inbounds) != 1 {
		t.Errorf("we expected the lengh of result be 1 but we got %d", len(inbounds))
	}

	if inbounds[0].UserID != userID4 ||
		inbounds[0].Domain != "twitter.com" {
		t.Error("error to compare data")
	}
}

func TestRetriveActiveInboundExpired(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundRepo.Create(inbound4)
	inboundRepo.Create(inbound5)
	inboundRepo.Create(inbound6)
	inboundCreated1, _ := inboundRepo.Create(inbound7)
	inboundRepo.Create(inbound15)
	inboundCreated2, _ := inboundRepo.Create(inbound16)

	inbounds, err := inboundRepo.RetriveActiveInboundExpired()
	if err != nil {
		t.Errorf("something went wrong the problem was %v", err)
	}

	if len(inbounds) != 2 {
		t.Errorf("we expected the lengh of result be 2 but we got %d", len(inbounds))
	}
	seen := map[string]struct{}{}
	for _, inbound := range inbounds {
		seen[inbound.ID] = struct{}{}
	}
	if _, ok := seen[inboundCreated1.ID]; !ok {
		t.Error("error to compare data")
	}
	if _, ok := seen[inboundCreated2.ID]; !ok {
		t.Error("error to compare data")
	}
}

func TestRetriveActiveInboundsOverQuota(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundRepo.Create(inbound4)
	inboundRepo.Create(inbound5)
	inboundRepo.Create(inbound6)
	inboundRepo.Create(inbound7)
	inboundRepo.Create(inbound15)
	inboundCreated, _ := inboundRepo.Create(inbound16)

	inbounds, err := inboundRepo.RetriveActiveInboundsOverQuota()
	if err != nil {
		t.Errorf("something went wrong the problem was %v", err)
	}

	if len(inbounds) != 1 {
		t.Errorf("we expected the lengh of result be 1 but we got %d", len(inbounds))
	}

	if inbounds[0].ID != inboundCreated.ID {
		t.Error("error to compare data")
	}
}

func TestRetriveDeactiveInboundsCharged(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundRepo.Create(inbound4)
	inboundRepo.Create(inbound5)
	inboundRepo.Create(inbound6)
	inboundRepo.Create(inbound7)
	inboundRepo.Create(inbound15)
	inboundRepo.Create(inbound16)
	inboundCreated, _ := inboundRepo.Create(inbound18)

	inbounds, err := inboundRepo.RetriveDeactiveInboundsCharged()
	if err != nil {
		t.Errorf("something went wrong the problem was %v", err)
	}

	if len(inbounds) != 1 {
		t.Errorf("we expected the lengh of result be 1 but we got %d", len(inbounds))
	}

	if inbounds[0].ID != inboundCreated.ID {
		t.Error("error to compare data")
	}
}

func TestInboundsIsNotAssigned(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundRepo.Create(inbound10)
	inboundRepo.Create(inbound11)
	inboundRepo.Create(inbound17)

	inbounds, err := inboundRepo.FindInboundIsNotAssigned()
	if err != nil {
		t.Fatalf("someting went wront the problem was %v", err)
	}
	if len(inbounds) != 1 {
		t.Fatalf("we expected we got 1 items but we got %v", len(inbounds))
	}
	if inbounds[0].UserID != userID9 {
		t.Fatalf("the answer was wrong")
	}
}

func TestFindInboundByUserID(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundCreated, _ := inboundRepo.Create(inbound10)
	inbound, err := inboundRepo.FindInboundByID(inboundCreated.ID)
	if err != nil {
		t.Fatalf("the query was field the problem was %v", err)
	}
	if inbound.UserID != inbound10.UserID || inboundCreated.ID != inbound.ID {
		t.Fatalf("the query has answerd wrong")
	}
}

func TestUpdateInbound(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundCreated, _ := inboundRepo.Create(inbound10)
	newDomain := "facebook.com"
	newPort := "2020"
	err := inboundRepo.UpdateDomainPort(inboundCreated.ID, newDomain, newPort, "3")
	if err != nil {
		t.Fatalf("update has field the problem was %v", err)
	}

	inbound, _ := inboundRepo.FindInboundByID(inboundCreated.ID)

	if inbound.VPNID != "3" ||
		inbound.Domain != newDomain ||
		inbound.Port != newPort ||
		inbound.IsAssigned != true {
		t.Fatalf("update hasn't worked carefuly")
	}
}

func TestGetListOfPortsByDomain(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundRepo.Create(inbound12)
	inboundRepo.Create(inbound13)
	inboundRepo.Create(inbound14)

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
	defer inboundRepo.DeleteAll()
	inboundCreated, _ := inboundRepo.Create(inbound1)

	err := inboundRepo.ChangeBlockState(inboundCreated.ID, true)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	inboundFound, _ := inboundRepo.FindInboundByID(inboundCreated.ID)

	if !inboundFound.IsBlock {
		t.Fatal("inbound that is founded must be blocked")
	}
}

func TestUnBlock(t *testing.T) {
	inboundCreated, _ := inboundRepo.Create(inbound5)
	defer inboundRepo.DeleteAll()
	err := inboundRepo.ChangeBlockState(inboundCreated.ID, false)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	inboundFound, _ := inboundRepo.FindInboundByID(inboundCreated.ID)

	if inboundFound.IsBlock {
		t.Fatal("inbound that is founded must be blocked")
	}
}

func TestUpdate(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundCreated, _ := inboundRepo.Create(inbound1)

	err := inboundRepo.Update(inboundCreated.ID, &inboundDto.UpdateInboundDto{
		Start:        utils.GetDateTime("2026-04-21 14:30:00"),
		End:          utils.GetDateTime("2026-04-22 14:30:00"),
		TrafficLimit: 200,
	})
	if err != nil {
		t.Fatalf("the error was %v", err)
	}

	inboundFound, _ := inboundRepo.FindInboundByID(inboundCreated.ID)

	if inboundFound.Start != utils.GetDateTime("2026-04-21 14:30:00") {
		t.Errorf("we expected start be %s", inboundFound.Start.Format(time.DateTime))
	}

	if inboundFound.TrafficLimit != 200 {
		t.Errorf("we expected traffic limit was 200 but we got %d", inboundFound.TrafficLimit)
	}

	if inboundFound.End != utils.GetDateTime("2026-04-22 14:30:00") {
		t.Errorf("we expected start be %s", inboundFound.End.Format(time.DateTime))
	}

	inboundCreated, _ = inboundRepo.Create(inbound2)

	err = inboundRepo.Update(inboundCreated.ID, &inboundDto.UpdateInboundDto{
		Start: utils.GetDateTime("2026-04-21 14:30:00"),
	})
	if err != nil {
		t.Fatalf("the error was %v", err)
	}

	inboundFound, _ = inboundRepo.FindInboundByID(inboundCreated.ID)

	if inboundFound.Start != utils.GetDateTime("2026-04-21 14:30:00") {
		t.Errorf("we expected start be %s", inboundFound.Start.Format(time.DateTime))
	}
}

func TestExtendInbound(t *testing.T) {
	defer inboundRepo.DeleteAll()

	inboundCreated, _ := inboundRepo.Create(inbound1)
	endTime, _ := timeTransformer.ConvertStrToTime("2026-04-21 14:30:00")
	endTime.Truncate(time.Second)
	err := inboundRepo.ExtendInbound(inboundCreated.ID, &inboundDto.ExtendInboundDto{
		TrafficExtended: 200,
		End:             endTime,
	})
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	err = inboundRepo.ExtendInbound(inboundCreated.ID, &inboundDto.ExtendInboundDto{
		TrafficExtended: 400,
		End:             endTime,
	})
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	inboundFound, _ := inboundRepo.FindInboundByID(inboundCreated.ID)
	if inboundFound.TrafficLimit != 600 {
		t.Fatalf("the field of trafficlimit must be 600 but we got %d", inboundFound.TrafficLimit)
	}

	if !inboundFound.End.Equal(endTime) {
		t.Fatalf("the end field must be %s but we got %s", inboundFound.End.Format(time.DateTime), endTime.Format(time.DateTime))
	}
}

func TestIncreateTrafficUsage(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundCreated, _ := inboundRepo.Create(inbound1)

	var trafficUsage uint32 = 50000
	err := inboundRepo.IncreaseTrafficUsage(inboundCreated.ID, trafficUsage)
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	inboundFound, _ := inboundRepo.FindInboundByID(inboundCreated.ID)
	if inboundFound.TrafficUsage != trafficUsage {
		t.Fatalf("error to compare data")
	}
}

func TestRetriveFinishedInbounds(t *testing.T) {
	defer inboundRepo.DeleteAll()

	inboundCreated1, _ := inboundRepo.Create(inbound19)
	inboundCreated2, _ := inboundRepo.Create(inbound20)
	inboundRepo.Create(inbound21)
	inboundRepo.Create(inbound22)

	inbounds, err := inboundRepo.RetriveFinishedInbounds()
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(inbounds) != 2 {
		t.Fatalf("we expected the lengh of inbound be 2 but we got %d", len(inbounds))
	}

	if !(inboundCreated1.ID != inbounds[0].ID || inboundCreated1.ID != inbounds[1].ID) {
		t.Fatalf("we expected the inbound that are returned contain %s", inboundCreated1.ID)
	}

	if !(inboundCreated2.ID != inbounds[0].ID || inboundCreated2.ID != inbounds[1].ID) {
		t.Fatalf("we expected the inbound that is returned contain %s", inboundCreated2.ID)
	}
}

func TestActiveInbounds(t *testing.T) {
	defer inboundRepo.DeleteAll()
	inboundCreated1, _ := inboundRepo.Create(inbound23)
	inboundCreated2, _ := inboundRepo.Create(inbound24)
	inboundRepo.Create(inbound25)

	inbounds, err := inboundRepo.RetriveActiveInbounds()
	if err != nil {
		t.Fatalf("something went worng that was %v", err)
	}
	if len(inbounds) != 2 {
		t.Fatalf("we expeted the lentgh of inbounds be %d but we got %d", 2, len(inbounds))
	}
	inboundIDMap := make(map[string]struct{})
	inboundIDMap[inboundCreated1.ID] = struct{}{}
	inboundIDMap[inboundCreated2.ID] = struct{}{}
	if _, isExist := inboundIDMap[inboundCreated1.ID]; !isExist {
		t.Fatalf("error to compare data inbound1 is missed")
	}

	if _, isExist := inboundIDMap[inboundCreated2.ID]; !isExist {
		t.Fatalf("error to compare data inbound2 is missed")
	}
}
