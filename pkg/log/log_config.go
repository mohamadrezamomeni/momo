package log

type LogConfig struct {
	AccessFile string `koanf:"accessFile"`
	ErrorFile  string `koanf:"errorFile"`
}
