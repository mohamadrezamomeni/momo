package entity

type VPN struct {
	ID        int
	Domain    string
	IsActive  bool
	ApiPort   string
	VPNType   VPNType
	UserCount int
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

func ConvertStringVPNTypeToEnum(key string) int {
	switch key {
	case "xray":
		return XRAY_VPN
	case "":
		return 0
	default:
		return UknownVPNType
	}
}
