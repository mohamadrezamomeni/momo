package vpnsource

import (
	"testing"

	vpnsource "github.com/mohamadrezamomeni/momo/dto/controller/vpn_source"
)

var VPNSourceValidator Validator = Validator{}

func TestUpsertingVPNSource(t *testing.T) {
	err := VPNSourceValidator.ValidateUpsert(vpnsource.CreateVPNSourceDto{
		IDentifyVPNSource: vpnsource.IDentifyVPNSource{
			Country: "england",
		},
		English: "uk",
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	err = VPNSourceValidator.ValidateUpsert(vpnsource.CreateVPNSourceDto{
		English: "uk",
	})
	if err == nil {
		t.Fatal("we expected an error beacuse we didn't fill english")
	}
	err = VPNSourceValidator.ValidateUpsert(vpnsource.CreateVPNSourceDto{})
	if err == nil {
		t.Fatal("we expected an error beacuse we didn't fill english")
	}

	err = VPNSourceValidator.ValidateUpsert(vpnsource.CreateVPNSourceDto{
		IDentifyVPNSource: vpnsource.IDentifyVPNSource{
			Country: "england",
		},
	})
	if err == nil {
		t.Fatal("we expected an error beacuse we didn't fill english")
	}
}
