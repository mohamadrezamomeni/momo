package worker

type WorkerConfig struct {
	Address string `koanf:"address"`
	Port    string `koanf:"port"`
}
