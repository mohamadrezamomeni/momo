package inbound

import (
	"fmt"
	"os"
	"testing"
	"time"

	"momo/pkg/config"
	"momo/repository/migrate"
	"momo/repository/sqllite"
	"momo/repository/sqllite/inbound/dto"

	"momo/proxy/vpn"

	"github.com/google/uuid"
)

var inboundRepo *Inbound

var (
	layout          = "2006-01-02 15:04:05"
	ts, _           = time.Parse(layout, "2025-04-21 14:30:00")
	te, _           = time.Parse(layout, "2025-04-22 14:30:00")
	port1           = "1081"
	port2           = "1082"
	port3           = "1083"
	userID1         = uuid.New().String()
	userID2         = uuid.New().String()
	inboundExample1 = &dto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port1),
		Protocol: "vmess",
		Port:     port1,
		Domain:   "google.com",
		UserID:   userID1,
		VPNType:  vpn.XRAY_VPN,
		Start:    ts,
		End:      te,
	}

	inboundExample2 = &dto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port2),
		Protocol: "vmess",
		Port:     port2,
		Domain:   "twitter.com",
		UserID:   userID2,
		VPNType:  vpn.XRAY_VPN,
		Start:    ts,
		End:      te,
	}

	inboundExample3 = &dto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port2),
		Protocol: "http",
		Port:     port2,
		Domain:   "googoo.com",
		UserID:   userID2,
		VPNType:  vpn.XRAY_VPN,
		Start:    ts,
		End:      te,
	}

	inboundExample4 = &dto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port3),
		Protocol: "http",
		Port:     port3,
		Domain:   "googoo.com",
		UserID:   userID2,
		VPNType:  vpn.XRAY_VPN,
		IsActive: true,
		Start:    ts,
		End:      te,
	}
)

func TestMain(m *testing.M) {
	cfg, err := config.Load("config_test.yaml")
	if err != nil {
		os.Exit(1)
	}
	db := sqllite.New(&cfg.DB)

	migrate := migrate.New(&cfg.DB)

	migrate.UP()

	inboundRepo = New(db)

	code := m.Run()
	os.Exit(code)
}

func TestCreateInbound(t *testing.T) {
	ret, err := inboundRepo.Create(inboundExample1)
	if err != nil {
		t.Errorf("something wrong has happend the problem was %v", err)
	}
	if ret.Domain != inboundExample1.Domain ||
		ret.Port != inboundExample1.Port ||
		ret.IsActive != false ||
		ret.UserID != inboundExample1.UserID ||
		ret.VPNType != inboundExample1.VPNType {
		t.Error("data wasn't saved currectly")
	}

	inboundRepo.Delete(ret.ID)
}

func TestChangeStatus(t *testing.T) {
	ret, err := inboundRepo.Create(inboundExample1)
	if err != nil {
		t.Errorf("something wrong has happend the problem was %v", err)
	}

	err = inboundRepo.changeStatus(ret.ID, false)
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}
	inboundRepo.Delete(ret.ID)
}

func TestFilterInbounds(t *testing.T) {
	i1, _ := inboundRepo.Create(inboundExample1)
	i2, _ := inboundRepo.Create(inboundExample2)
	i3, _ := inboundRepo.Create(inboundExample3)
	i4, _ := inboundRepo.Create(inboundExample4)
	inbounds1, err := inboundRepo.Filter(&dto.FilterInbound{Port: port2})
	if err != nil {
		t.Errorf("1. the problem has occured that is %v", err)
	}
	if len(inbounds1) != 2 {
		t.Errorf("1. the number of items must be %v but got %v items", 2, len(inbounds1))
	}

	inbounds2, err := inboundRepo.Filter(&dto.FilterInbound{Protocol: "vmess"})
	if err != nil {
		t.Errorf("2. the problem has occured that is %v", err)
	}
	if len(inbounds2) != 2 {
		t.Errorf("2. the number of items must be %v but got %v items", 2, len(inbounds2))
	}

	inbounds3, err := inboundRepo.Filter(&dto.FilterInbound{Domain: "google.com"})
	if err != nil {
		t.Errorf("3. the problem has occured that is %v", err)
	}
	if len(inbounds3) != 1 {
		t.Errorf("3. the number of items must be %v but got %v items", 1, len(inbounds3))
	}

	inbounds4, err := inboundRepo.Filter(&dto.FilterInbound{UserID: userID2, Port: port2})
	if err != nil {
		t.Errorf("4. the problem has occured that is %v", err)
	}
	if len(inbounds4) != 2 {
		t.Errorf("4. the number of items must be %v but got %v items", 2, len(inbounds4))
	}

	isAvailableT := true
	inbounds5, err := inboundRepo.Filter(&dto.FilterInbound{IsActice: &isAvailableT})
	if err != nil {
		t.Errorf("5. the problem has occured that is %v", err)
	}
	if len(inbounds5) != 1 {
		t.Errorf("5. the number of items must be %v but got %v items", 1, len(inbounds5))
	}

	deleteInbounds(i1.ID, i2.ID, i3.ID, i4.ID)
}

func deleteInbounds(ids ...int) {
	for _, id := range ids {
		inboundRepo.Delete(id)
	}
}
