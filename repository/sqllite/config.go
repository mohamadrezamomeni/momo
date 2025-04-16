package sqllite

type DBConfig struct {
	Dialect    string `koanf:"dialect"`
	Path       string `koanf:"path"`
	Migrations string `koanf:"migrations"`
}
