package config

import (
	"momo/pkg/log"
	"momo/proxy/xray"
)

type Config struct {
	Log        log.LogConfig   `koanf:"log"`
	XrayConfig xray.XrayConfig `koanf:"xray"`
}
