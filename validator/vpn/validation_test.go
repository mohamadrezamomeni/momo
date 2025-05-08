package vpn

import (
	"os"
	"testing"

	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
)

var validator *Validator

func TestMain(m *testing.M) {
	validator = New()
	code := m.Run()
	os.Exit(code)
}

func TestCreateVPN(t *testing.T) {
	err := validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:   "xray",
		Port:      "322",
		Domain:    "mohamadreza.com",
		UserCount: 2,
	})
	if err != nil {
		t.Errorf("someting went wrong that was %v", err)
	}
	err = validator.ValidateCreatingVPN(vpnControllerDto.CreateVPN{
		VpnType:   "xray",
		Port:      "3m22",
		Domain:    "mohamadreza.com",
		UserCount: 2,
	})

	if err == nil {
		t.Error("we expected we would get err but we got nothing")
	}
}
