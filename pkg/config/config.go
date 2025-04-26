package config

import (
	workerServer "momo/delivery/worker"
	"momo/pkg/log"
	"momo/repository/sqllite"
)

type Config struct {
	Log            log.LogConfig               `koanf:"log"`
	DB             sqllite.DBConfig            `koanf:"db"`
	Metric         workerServer.MetricConfig   `koanf:"metric"`
	PortAssignment workerServer.PortAssignment `koanf:"port_assignment"`
	WorkerServer   workerServer.WorkerConfig   `koanf:"worker_server"`
}
