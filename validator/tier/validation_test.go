package charge

import (
	"os"
	"testing"

	"github.com/mohamadrezamomeni/momo/dto/controller/tier"
)

var validator *Validator

func TestMain(m *testing.M) {
	validator = New()
	code := m.Run()
	os.Exit(code)
}

func TestCreatingTier(t *testing.T) {
	isDefault := false
	err := validator.ValidateCreatingTier(tier.CreateTier{
		IdentifyTierDto: tier.IdentifyTierDto{
			Name: "silver",
		},
		IsDefault: &isDefault,
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	err = validator.ValidateCreatingTier(tier.CreateTier{
		IdentifyTierDto: tier.IdentifyTierDto{
			Name: "s",
		},
		IsDefault: &isDefault,
	})
	if err == nil {
		t.Fatalf("something went wrong that was %v", err)
	}
}
