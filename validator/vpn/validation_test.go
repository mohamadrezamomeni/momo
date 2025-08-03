package vpn

import (
	"os"
	"testing"

	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	mockVPNSource "github.com/mohamadrezamomeni/momo/mocks/service/vpn_source"
)

var validator *Validator

func TestMain(m *testing.M) {
	validator = New(mockVPNSource.New([]*entity.VPNSource{
		{
			Country: "uk",
			English: "uk",
		},
	}))
	code := m.Run()
	os.Exit(code)
}

func TestCreateVPN(t *testing.T) {
	err := validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:        "xray",
		Port:           "322",
		Domain:         "mohamadreza.com",
		Country:        "uk",
		UserCount:      2,
		StartPort:      1000,
		VPNStatusLabel: "drain",
		EndPort:        2000,
	})
	if err != nil {
		t.Errorf("someting went wrong that was %v", err)
	}

	err = validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:        "xray",
		Port:           "322",
		Domain:         "mohamadreza.com",
		Country:        "uk",
		UserCount:      2,
		StartPort:      1000,
		VPNStatusLabel: "draind",
		EndPort:        2000,
	})
	if err == nil {
		t.Errorf("someting went wrong that was %v", err)
	}

	err = validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:        "xray",
		Port:           "322",
		Domain:         "mohamadreza.com",
		Country:        "uk",
		UserCount:      2,
		StartPort:      1000,
		VPNStatusLabel: "",
		EndPort:        2000,
	})
	if err == nil {
		t.Errorf("someting went wrong that was %v", err)
	}

	err = validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:   "xray",
		Port:      "3m22",
		Domain:    "mohamadreza.com",
		Country:   "uk",
		UserCount: 2,
		StartPort: 1000,
		EndPort:   2000,
	})

	if err == nil {
		t.Error("we expected we would get err but we got nothing")
	}

	err = validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:   "xray",
		Port:      "322",
		Domain:    "mohamadreza.com",
		UserCount: 2,
		StartPort: 1000,
		EndPort:   2000,
	})

	if err == nil {
		t.Error("we expected we would get error but we got nothing")
	}

	err = validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:   "xray",
		Port:      "322",
		Domain:    "mohamadreza.com",
		UserCount: 2,
		Country:   "china",
		StartPort: 1000,
		EndPort:   2000,
	})

	if err == nil {
		t.Error("we expected we would get error but we got nothing")
	}

	err = validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:   "xray",
		Port:      "322",
		Domain:    "mohamadreza.com",
		UserCount: 2,
		Country:   "china",
	})
	if err == nil {
		t.Error("we expect the error beacuse endport is 0")
	}
}

func TestValidateFilter(t *testing.T) {
	err := validator.ValidateFilterVPNs(vpnControllerDto.FilterVPNs{
		Domain:  "twitter.com",
		VPNType: "xray",
	})
	if err != nil {
		t.Errorf("something went wrong that was %v", err)
	}
	err = validator.ValidateFilterVPNs(vpnControllerDto.FilterVPNs{
		VPNType: "xray",
	})
	if err != nil {
		t.Errorf("something went wrong that was %v", err)
	}
	err = validator.ValidateFilterVPNs(vpnControllerDto.FilterVPNs{})
	if err != nil {
		t.Errorf("something went wrong that was %v", err)
	}
	err = validator.ValidateFilterVPNs(vpnControllerDto.FilterVPNs{
		Domain: "twitter.com",
	})
	if err != nil {
		t.Errorf("something went wrong that was %v", err)
	}
	err = validator.ValidateFilterVPNs(vpnControllerDto.FilterVPNs{
		Domain: "fggggg",
	})
	if err == nil {
		t.Errorf("we expected an error")
	}
}

func TestValidateUpdate(t *testing.T) {
	err := validator.ValidateUpdatingVPN(vpnControllerDto.UpdateVPN{
		VPNStatusLabel: "cordon",
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	err = validator.ValidateUpdatingVPN(vpnControllerDto.UpdateVPN{
		VPNStatusLabel: "cordonn",
	})
	if err == nil {
		t.Fatal("we expected error would be occured")
	}
}
