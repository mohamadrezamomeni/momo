package config

import "momo/pkg/log"

type Config struct {
	Log log.LogConfig `koanf:"log"`
}
