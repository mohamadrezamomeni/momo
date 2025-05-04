package auth

type AuthConfig struct {
	SecretKey  string `koanf:"secret_key"`
	ExpireTime uint64 `koanf:"expire_time_by_minutes"`
}
