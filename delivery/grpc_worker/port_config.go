package worker

type PortAssignment struct {
	StartPort int `koanf:"start_port"`
	EndPort   int `koanf:"end_port"`
}
