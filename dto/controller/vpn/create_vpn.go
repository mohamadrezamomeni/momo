package vpn

type CreateVPN struct {
	VpnType   string `json:"vpn_type"`
	Port      string `json:"port"`
	Domain    string `json:"domain"`
	UserCount int    `json:"user_count"`
}
