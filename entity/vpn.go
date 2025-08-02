package entity

type VPN struct {
	ID        int
	Domain    string
	IsActive  bool
	ApiPort   string
	VPNType   VPNType
	UserCount int
	Country   string
	StartPort int
	EndPort   int
}

type VPNType = int

const (
	XRAY_VPN VPNType = iota + 1
	UknownVPNType
)

func VPNTypeString(vpnType int) string {
	switch vpnType {
	case XRAY_VPN:
		return "xray"
	default:
		return "unkhown"
	}
}

func ConvertStringVPNTypeToEnum(key string) VPNType {
	switch key {
	case "xray":
		return XRAY_VPN
	case "unkhown":
		return UknownVPNType
	default:
		return 0
	}
}
