package vpn

import (
	momoError "momo/pkg/error"
	"momo/proxy/vpn/internal/xray"
)

type VPN struct {
	xray xray.IXray
}

func New(cfg Config) *VPN {
	x, err := xray.New(cfg.Xray)
	if err != nil {
		momoError.Fatalf("something went wrong to initialize xray error: %v", err)
	}

	return &VPN{
		xray: x,
	}
}
