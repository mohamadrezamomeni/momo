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
	VPNStatus VPNStatus
}

type (
	VPNType   = int
	VPNStatus = int
)

const (
	Cordon = iota + 1
	Drain
	Ready
	UnkhownVPNStatus
)

const (
	XRAY_VPN VPNType = iota + 1
	UknownVPNType
)

func VPNStatusString(VPNStatus VPNType) string {
	switch VPNStatus {
	case Cordon:
		return "cordon"
	case Drain:
		return "drain"
	case Ready:
		return "ready"
	default:
		return "unkhown"
	}
}

func ConvertVPNStatusLabelToVPNStatus(label string) VPNStatus {
	switch label {
	case "cordon":
		return Cordon
	case "drain":
		return Drain
	case "ready":
		return Ready
	default:
		return UnkhownVPNStatus
	}
}

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
