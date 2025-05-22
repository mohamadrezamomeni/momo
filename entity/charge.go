package entity

type Charge struct {
	ID           string
	Status       ChargeStatus
	AdminComment string
	Detail       string
	InboundID    string
}

type ChargeStatus = int

const (
	UnkhownStatusCharge = iota + 1
	ApprovedStatusCharge
	PendingStatusCharge
	RegejectedStatusCharge
)

func TranslateChargeStatus(enum ChargeStatus) string {
	switch enum {
	case ApprovedStatusCharge:
		return "approved"
	case PendingStatusCharge:
		return "pending"
	case RegejectedStatusCharge:
		return "regected"
	default:
		return "unkhown"
	}
}

func ConvertStringToChargeStatus(s string) ChargeStatus {
	switch s {
	case "approved":
		return ApprovedStatusCharge
	case "pending":
		return PendingStatusCharge
	case "regected":
		return RegejectedStatusCharge
	default:
		return UnkhownStatusCharge
	}
}
