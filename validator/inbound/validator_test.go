package inbound

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	inboundServiceMock "github.com/mohamadrezamomeni/momo/mocks/service/inbound"
	userServiceMock "github.com/mohamadrezamomeni/momo/mocks/service/user"
)

var (
	validator      *Validator
	inboundSvcMock = inboundServiceMock.New()
)

func TestMain(m *testing.M) {
	validator = New(
		userServiceMock.New(),
		inboundSvcMock,
	)

	code := m.Run()
	os.Exit(code)
}

func TestCreatingInbound(t *testing.T) {
	err := validator.ValidateCreatingInbound(inboundControllerDto.CreateInbound{
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

	err = validator.ValidateCreatingInbound(inboundControllerDto.CreateInbound{
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
	err = validator.ValidateCreatingInbound(inboundControllerDto.CreateInbound{
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

	err = validator.ValidateCreatingInbound(inboundControllerDto.CreateInbound{
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

	err = validator.ValidateCreatingInbound(inboundControllerDto.CreateInbound{
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

	err = validator.ValidateCreatingInbound(inboundControllerDto.CreateInbound{
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

func TestValidateFilteringInbounds(t *testing.T) {
	err := validator.ValidateFilteringInbounds(inboundControllerDto.FilterInboundsDto{})
	if err != nil {
		t.Errorf("we exptected not errors that  was %v", err)
	}
	err = validator.ValidateFilteringInbounds(inboundControllerDto.FilterInboundsDto{Domain: "tww"})
	if err == nil {
		t.Errorf("we exptected errors")
	}

	err = validator.ValidateFilteringInbounds(inboundControllerDto.FilterInboundsDto{Port: "21f"})
	if err == nil {
		t.Errorf("we exptected errors")
	}

	err = validator.ValidateFilteringInbounds(inboundControllerDto.FilterInboundsDto{UserID: "asdfad-234"})
	if err == nil {
		t.Errorf("we exptected errors")
	}

	err = validator.ValidateFilteringInbounds(inboundControllerDto.FilterInboundsDto{VPNType: "xrayy"})
	if err == nil {
		t.Errorf("we exptected errors")
	}
}

func TestValidateExtendingInbound(t *testing.T) {
	now := time.Now()
	now = now.Truncate(time.Second)
	inboundCreated, _ := inboundSvcMock.Create(&inboundServiceDto.CreateInbound{
		ServerType: entity.High,
		Start:      now.AddDate(0, -1, 0),
		End:        now.AddDate(0, 1, 0),
	})

	defer inboundSvcMock.DeletedAll()

	err := validator.ValidateExtendingInbound(inbound.ExtendInboundDto{
		IdentifyInbounbdDto: inbound.IdentifyInbounbdDto{
			ID: strconv.Itoa(inboundCreated.ID),
		},
		End: now.AddDate(0, 2, 0).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		t.Errorf("something went wrong the problem was %v", err)
	}

	err = validator.ValidateExtendingInbound(inbound.ExtendInboundDto{
		IdentifyInbounbdDto: inbound.IdentifyInbounbdDto{
			ID: strconv.Itoa(inboundCreated.ID),
		},
		End: now.AddDate(0, 1, 0).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		t.Error("we expected input end will be greater than current end")
	}

	err = validator.ValidateExtendingInbound(inbound.ExtendInboundDto{
		IdentifyInbounbdDto: inbound.IdentifyInbounbdDto{
			ID: strconv.Itoa(inboundCreated.ID),
		},
		End: now.AddDate(0, 0, 2).Format("2006-01-02 15:04:05"),
	})

	if err == nil {
		t.Error("we expected input end will be greater than current end")
	}

	inboundCreated, _ = inboundSvcMock.Create(&inboundServiceDto.CreateInbound{
		ServerType: entity.High,
		Start:      now.AddDate(-1, -1, 0),
		End:        now.AddDate(-1, 0, 0),
	})

	err = validator.ValidateExtendingInbound(inbound.ExtendInboundDto{
		IdentifyInbounbdDto: inbound.IdentifyInbounbdDto{
			ID: strconv.Itoa(inboundCreated.ID),
		},
		End: now.AddDate(1, 0, 2).Format("2006-01-02 15:04:05"),
	})

	if err == nil {
		t.Error("current time must be lower than current end time")
	}
}
