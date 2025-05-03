package vpn

import "github.com/mohamadrezamomeni/momo/entity"

type VPNConfig struct {
	Port    string
	VPNType entity.VPNType
	Domain  string
}
