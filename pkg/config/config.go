package config

import (
	"momo/pkg/log"
	"momo/proxy/xray"

	"momo/repository/sqllite"
)

type Config struct {
	Log        log.LogConfig    `koanf:"log"`
	XrayConfig xray.XrayConfig  `koanf:"xray"`
	DB         sqllite.DBConfig `koanf:"db"`
}
