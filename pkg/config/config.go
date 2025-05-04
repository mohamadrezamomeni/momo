package config

import (
	workerServer "github.com/mohamadrezamomeni/momo/delivery/grpc_worker"
	httpserver "github.com/mohamadrezamomeni/momo/delivery/http_server"
	"github.com/mohamadrezamomeni/momo/pkg/log"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
	auth "github.com/mohamadrezamomeni/momo/service/auth"
	encrypt "github.com/mohamadrezamomeni/momo/service/crypt"
)

type Config struct {
	CryptConfig    encrypt.CryptConfig         `koanf:"encrypt"`
	AuthConfig     auth.Auth                   `koanf:"auth"`
	HTTP           httpserver.HTTPConfig       `koanf:"http"`
	Log            log.LogConfig               `koanf:"log"`
	DB             sqllite.DBConfig            `koanf:"db"`
	Metric         workerServer.MetricConfig   `koanf:"metric"`
	PortAssignment workerServer.PortAssignment `koanf:"port_assignment"`
	WorkerServer   workerServer.WorkerConfig   `koanf:"worker_server"`
}
