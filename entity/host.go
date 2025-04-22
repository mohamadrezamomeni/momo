package entity

type Host struct {
	Domain string
	ID     int
	Port   string
	Status HostStatus
}

type HostStatus = int

func (h *Host) HostStatusString() string {
	switch h.Status {
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

const (
	High HostStatus = iota
	Medium
	Low
	Deactive
)
