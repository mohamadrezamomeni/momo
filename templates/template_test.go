package templates

import (
	"testing"

	"github.com/mohamadrezamomeni/momo/entity"
)

func TestClientConfig(t *testing.T) {
	_, err := LoadClientConfig(entity.XRAY_VPN, "instagram.com", "1234", "1234")
	if err != nil {
		t.Fatalf("someting went wrong that was %v", err)
	}
}
