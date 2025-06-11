package entity

type Charge struct {
	ID           string
	Status       ChargeStatus
	AdminComment string
	Detail       string
	InboundID    string
	UserID       string
	PackageID    string
}

type ChargeStatus = int

const (
	UnkhownStatusCharge = iota + 1
	ApprovedStatusCharge
	PendingStatusCharge
	RegejectedStatusCharge
	AssignedCharged
)

func TranslateChargeStatus(enum ChargeStatus) string {
	switch enum {
	case ApprovedStatusCharge:
		return "approved"
	case PendingStatusCharge:
		return "pending"
	case RegejectedStatusCharge:
		return "regected"
	case AssignedCharged:
		return "assiggned"
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
	case "assiggned":
		return AssignedCharged
	case "unkhown":
		return UnkhownStatusCharge
	default:
		return 0
	}
}
