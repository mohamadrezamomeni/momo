package vpn

type CreateVPN struct {
	VpnType   string `json:"vpn_type"`
	Port      string `json:"port"`
	Domain    string `json:"domain"`
	UserCount int    `json:"user_count"`
	Country   string `json:"country"`
	StartPort int    `json:"start_port"`
	EndPort   int    `json:"end_port"`
}
