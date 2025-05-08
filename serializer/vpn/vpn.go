package vpn

type VPNSerializer struct {
	ID        int
	Domain    string
	IsActive  bool
	ApiPort   string
	VPNType   string
	UserCount int
}
