package config

import (
	workerServer "github.com/mohamadrezamomeni/momo/delivery/worker"
	"github.com/mohamadrezamomeni/momo/pkg/log"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

type Config struct {
	Log            log.LogConfig               `koanf:"log"`
	DB             sqllite.DBConfig            `koanf:"db"`
	Metric         workerServer.MetricConfig   `koanf:"metric"`
	PortAssignment workerServer.PortAssignment `koanf:"port_assignment"`
	WorkerServer   workerServer.WorkerConfig   `koanf:"worker_server"`
}
