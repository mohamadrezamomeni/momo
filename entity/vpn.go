package entity

import momoError "momo/pkg/error"

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
	XRAY_VPN VPNType = iota
)

func (v *VPN) VPNTypeString() string {
	switch v.VPNType {
	case XRAY_VPN:
		return "xray"
	default:
		return "UNKHOWN"
	}
}

func ConvertStringVPNTypeToEnum(key string) (int, error) {
	switch key {
	case "xray":
		return XRAY_VPN, nil
	default:
		return 0, momoError.Error("unexpected vpn_type's value")
	}
}
