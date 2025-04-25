package entity

import (
	momoError "momo/pkg/error"
)

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
	UknownVPNTYPe
)

func VPNTypeString(vpnType int) string {
	switch vpnType {
	case XRAY_VPN:
		return "xray"
	default:
		return "unkhown"
	}
}

func ConvertStringVPNTypeToEnum(key string) (int, error) {
	switch key {
	case "xray":
		return XRAY_VPN, nil
	default:
		return UknownVPNTYPe, momoError.Error("unexpected vpn_type's value")
	}
}
