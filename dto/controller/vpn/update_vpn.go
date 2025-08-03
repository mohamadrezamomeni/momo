package vpn

type UpdateVPN struct {
	IdentifyVPNDto
	VPNStatusLabel string `json:"status"`
}
