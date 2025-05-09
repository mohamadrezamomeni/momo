package inbound

type FilterInboundsDto struct {
	Protocol string `query:"protocol"`
	Domain   string `query:"domain"`
	VPNType  string `query:"vpn_type"`
	Port     string `query:"port"`
	UserID   string `query:"user_id"`
}
