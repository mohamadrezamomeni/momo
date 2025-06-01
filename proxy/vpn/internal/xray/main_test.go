package xray

import (
	"os"
	"testing"
)

var xrayU *Xray

func TestMain(m *testing.M) {
	xrayU, _ = New(&XrayConfig{
		Address: "192.168.116.129",
		ApiPort: "62789",
	})

	code := m.Run()
	os.Exit(code)
}
