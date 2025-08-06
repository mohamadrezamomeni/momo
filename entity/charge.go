package entity

type Charge struct {
	ID           string
	Status       ChargeStatus
	AdminComment string
	Detail       string
	InboundID    string
	UserID       string
	PackageID    string
	VPNType      VPNType
	Country      string
}

type ChargeStatus = int

const (
	UnkhownStatusCharge ChargeStatus = iota + 1
	ApprovedStatusCharge
	PendingStatusCharge
	RejectedStatusCharge
	AssignedCharged
)

func TranslateChargeStatus(enum ChargeStatus) string {
	switch enum {
	case ApprovedStatusCharge:
		return "approved"
	case PendingStatusCharge:
		return "pending"
	case RejectedStatusCharge:
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
		return RejectedStatusCharge
	case "assiggned":
		return AssignedCharged
	case "unkhown":
		return UnkhownStatusCharge
	default:
		return 0
	}
}

func ConvertLabelsToChargeStatuses(labels []string) []ChargeStatus {
	statuses := make([]ChargeStatus, 0)
	for _, label := range labels {
		status := ConvertStringToChargeStatus(label)
		if status != 0 {
			statuses = append(statuses, status)
		}
	}
	return statuses
}

func ConvertStatusesToStatusLabels(statuses []ChargeStatus) []string {
	lablels := make([]string, 0)
	for _, status := range statuses {
		label := TranslateChargeStatus(status)
		lablels = append(lablels, label)
	}
	return lablels
}
