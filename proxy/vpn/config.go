package vpn

import "momo/entity"

type VPNConfig struct {
	Port    string
	VPNType entity.VPNType
	Domain  string
}
