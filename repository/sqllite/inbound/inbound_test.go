package inbound

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"momo/pkg/config"
	"momo/pkg/utils"
	"momo/repository/migrate"
	"momo/repository/sqllite"
	"momo/repository/sqllite/inbound/dto"

	"momo/proxy/vpn"

	"github.com/google/uuid"
)

var inboundReop *Inbound

func TestMain(m *testing.M) {
	cfg, err := config.Load("config_test.yaml")
	if err != nil {
		os.Exit(1)
	}
	db := sqllite.New(&cfg.DB)

	migrate := migrate.New(&cfg.DB)

	migrate.UP()

	inboundReop = New(db)

	code := m.Run()
	os.Exit(code)
}

var (
	port           = strconv.Itoa(utils.GenerateRandomInRange(1080, 1089))
	domain         = "google.com"
	id             = uuid.New()
	inboundExample = &dto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port),
		Protocol: "vmess",
		Port:     port,
		Domain:   domain,
		UserID:   id.String(),
		VPNType:  vpn.XRAY_VPN,
	}
)

func TestCreateInbound(t *testing.T) {
	ret, err := inboundReop.Create(inboundExample)
	if err != nil {
		t.Errorf("something wrong has happend the problem was %v", err)
	}
	if ret.Domain != domain || ret.Port != port || ret.IsAvailable != false || ret.UserID != id.String() || ret.VPNType != vpn.XRAY_VPN {
		t.Error("data wasn't saved currectly")
	}

	inboundReop.Delete(ret.ID)
}

func TestChangeStatus(t *testing.T) {
	ret, err := inboundReop.Create(inboundExample)
	if err != nil {
		t.Errorf("something wrong has happend the problem was %v", err)
	}

	err = inboundReop.changeStatus(ret.ID, false)
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}
	inboundReop.Delete(ret.ID)
}
