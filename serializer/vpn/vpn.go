package vpn

type VPNSerializer struct {
	ID        int    `json:"ID"`
	Domain    string `json:"domain"`
	IsActive  bool   `json:"is_active"`
	ApiPort   string `json:"api_port"`
	VPNType   string `json:"vpn_type"`
	UserCount int    `json:"user_count"`
	Country   string `json:"country"`
	StartPort int    `json:"start_port"`
	EndPort   int    `json:"end_port"`
}
