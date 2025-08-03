package charge

import (
	"os"
	"testing"

	"github.com/mohamadrezamomeni/momo/dto/controller/charge"
)

var validator *Validator

func TestMain(m *testing.M) {
	validator = New()
	code := m.Run()
	os.Exit(code)
}

func TestFilterCharges(t *testing.T) {
	err := validator.ValidateFilterCharges(charge.FilterCharges{})
	if err != nil {
		t.Fatalf("we expected no error but we got %v", err)
	}

	err = validator.ValidateFilterCharges(charge.FilterCharges{
		UserID: "0393ed06-29bb-41c2-b3f4-6382a6729c3e",
	})
	if err != nil {
		t.Fatalf("we expected no error but we got %v", err)
	}

	err = validator.ValidateFilterCharges(charge.FilterCharges{
		InboundID: "123",
	})
	if err != nil {
		t.Fatalf("we expected no error but we got %v", err)
	}

	err = validator.ValidateFilterCharges(charge.FilterCharges{
		Statuses:  "approved,pending",
		UserID:    "0393ed06-29bb-41c2-b3f4-6382a6729c3e",
		InboundID: "123",
	})
	if err != nil {
		t.Fatalf("we expected no error but we got %v", err)
	}

	err = validator.ValidateFilterCharges(charge.FilterCharges{
		Statuses: "approved,pending",
	})
	if err != nil {
		t.Fatalf("we expected no error but we got %v", err)
	}

	err = validator.ValidateFilterCharges(charge.FilterCharges{
		Statuses: "approvedd,pending",
	})
	if err == nil {
		t.Fatalf("we expected no error but we got %v", err)
	}

	err = validator.ValidateFilterCharges(charge.FilterCharges{
		Statuses: "assignedd",
	})
	if err == nil {
		t.Fatalf("we expected no error but we got %v", err)
	}
}
