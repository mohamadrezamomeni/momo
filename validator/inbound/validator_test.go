package inbound

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	userServiceMock "github.com/mohamadrezamomeni/momo/mocks/service/user"
)

var validator *Validator

func TestMain(m *testing.M) {
	validator = New(userServiceMock.New())

	code := m.Run()
	os.Exit(code)
}

func TestCreatingInbound(t *testing.T) {
	err := validator.ValidateCreatingInbound(inbound.CreateInbound{
		Protocol: "vmess",
		VPNType:  "xray",
		Domain:   "twitter.com",
		Port:     "234",
		UserID:   uuid.New().String(),
		Start:    "2025-11-01T15:00:00Z",
		End:      "2025-11-01T16:00:00Z",
	})
	if err != nil {
		t.Errorf("someting went wrong err was %v", err)
	}

	err = validator.ValidateCreatingInbound(inbound.CreateInbound{
		Protocol: "vmess",
		VPNType:  "xray",
		Domain:   "twitter.com",
		Port:     "234",
		UserID:   uuid.New().String(),
		Start:    "2025-11-01T15:00:00Z",
		End:      "2025-11-01T14:00:00Z",
	})
	if err == nil {
		t.Errorf("this validation must validate start be before end")
	}
	err = validator.ValidateCreatingInbound(inbound.CreateInbound{
		Protocol: "vmess",
		VPNType:  "xrayy",
		Domain:   "twitter.com",
		Port:     "234",
		UserID:   uuid.New().String(),
		Start:    "2025-11-01T15:00:00Z",
		End:      "2025-11-01T16:00:00Z",
	})
	if err == nil {
		t.Errorf("vpnType could be validated")
	}

	err = validator.ValidateCreatingInbound(inbound.CreateInbound{
		Protocol: "vmess",
		VPNType:  "xrayy",
		Domain:   "twitter.com",
		Port:     "2343s",
		UserID:   uuid.New().String(),
		Start:    "2025-11-01T15:00:00Z",
		End:      "2025-11-01T16:00:00Z",
	})
	if err == nil {
		t.Errorf("port could be validated")
	}

	err = validator.ValidateCreatingInbound(inbound.CreateInbound{
		Protocol: "vmess",
		VPNType:  "xray",
		Domain:   "twitte",
		Port:     "2343",
		UserID:   uuid.New().String(),
		Start:    "2025-11-01T15:00:00Z",
		End:      "2025-11-01T16:00:00Z",
	})
	if err == nil {
		t.Errorf("domain could be validated")
	}

	err = validator.ValidateCreatingInbound(inbound.CreateInbound{
		Protocol: "vmess",
		VPNType:  "xrayy",
		Domain:   "twitter.com",
		Port:     "2343",
		UserID:   uuid.New().String() + "3",
		Start:    "2025-11-01T15:00:00Z",
		End:      "2025-11-01T16:00:00Z",
	})
	if err == nil {
		t.Errorf("domain could be validated")
	}
}
