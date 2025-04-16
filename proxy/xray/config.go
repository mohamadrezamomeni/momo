package xray

type XrayConfig struct {
	Address    string `koanf:"address"`
	ApiPort    string `koanf:"api_port"`
	configPath string `koanf:"config"`
}
