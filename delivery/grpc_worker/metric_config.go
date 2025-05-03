package worker

type MetricConfig struct {
	Address       string `koanf:"address"`
	Port          string `koanf:"port"`
	CPUWeight     int    `koanf:"cpu_weight"`
	MemoryWeight  int    `koanf:"memory_weight"`
	NetworkWeight int    `koanf:"network_weight"`
	HighStatus    int    `koanf:"high_status"`
	MediumStatus  int    `koanf:"medium_status"`
	LowStatus     int    `koanf:"low_status"`
}
