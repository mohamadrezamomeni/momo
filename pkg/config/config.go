package config

import (
	"momo/pkg/log"

	"momo/repository/sqllite"
)

type Config struct {
	Log log.LogConfig    `koanf:"log"`
	DB  sqllite.DBConfig `koanf:"db"`
}
