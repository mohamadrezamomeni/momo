package host

type Host struct {
	Domain string `json:"domain"`
	Port   string `json:"port"`
	Rank   uint32 `json:"rank"`
	Status string `json:"status"`
}
