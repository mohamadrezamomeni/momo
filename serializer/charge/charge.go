package charge

type Charge struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	AdminComment string `json:"admin_coment"`
	Detail       string `json:"detail"`
	InboundID    string `json:"inbound_id"`
	UserID       string `json:"user_id"`
	PackageID    string `json:"package_id"`
	VPNType      string `json:"vpn_type"`
	Country      string `json:"country"`
}
