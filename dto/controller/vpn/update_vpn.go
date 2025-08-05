package vpn

type UpdateVPN struct {
	IdentifyVPNDto
	VPNStatusLabel string `json:"status"`
	Domain         string `json:"domain"`
	ApiPort        string `json:"api_port"`
}
