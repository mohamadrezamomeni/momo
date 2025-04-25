package config

import (
	metricServer "momo/delivery/metric_server"
	"momo/pkg/log"
	"momo/repository/sqllite"
)

type Config struct {
	Log          log.LogConfig             `koanf:"log"`
	DB           sqllite.DBConfig          `koanf:"db"`
	MetricServer metricServer.MetricConfig `koanf:"metric_server"`
}
