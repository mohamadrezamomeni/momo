package host

type CreateHostDto struct {
	Domain string `json:"domain"`
	Port   string `json:"port"`
}
