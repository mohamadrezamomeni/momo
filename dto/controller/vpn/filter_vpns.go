package vpn

type FilterVPNs struct {
	VPNType string `query:"vpn_type"`
	Domain  string `query:"domain"`
}
