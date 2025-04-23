package entity

import momoError "momo/pkg/error"

type Host struct {
	Domain string
	ID     int
	Port   string
	Status HostStatus
}

type HostStatus = int

func HostStatusString(status HostStatus) string {
	switch status {
	case High:
		return "high"
	case Medium:
		return "medium"
	case Low:
		return "low"
	default:
		return "deactive"
	}
}

func MapTuStatus(statusString string) (HostStatus, error) {
	switch statusString {
	case HighStr:
		return High, nil
	case MediumStr:
		return Medium, nil
	case LowStr:
		return Low, nil
	case DeactiveStr:
		return Deactive, nil
	default:
		return uknown, momoError.Errorf("the status of \"%s\" doesn't exist ", statusString)
	}
}

const (
	High HostStatus = iota
	Medium
	Low
	Deactive
	uknown
)

const (
	HighStr     = "high"
	MediumStr   = "medium"
	LowStr      = "low"
	DeactiveStr = "deactive"
)
