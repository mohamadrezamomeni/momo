package main

import (
	"log"

	workerServer "github.com/mohamadrezamomeni/momo/delivery/worker"
	"github.com/mohamadrezamomeni/momo/pkg/config"
	momoLogger "github.com/mohamadrezamomeni/momo/pkg/log"
	metricService "github.com/mohamadrezamomeni/momo/service/metric"
	portService "github.com/mohamadrezamomeni/momo/service/port"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding config \n - the problem was %v", err)
	}
	momoLogger.Init(cfg.Log)

	metricSvc := metricService.New(&cfg.Metric)

	portSvc := portService.New(&cfg.PortAssignment)

	server := workerServer.New(metricSvc, portSvc, cfg.WorkerServer)

	server.Start()
}
