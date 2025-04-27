package main

import (
	"log"

	workerServer "momo/delivery/worker"
	"momo/pkg/config"
	momoLogger "momo/pkg/log"
	metricService "momo/service/metric"
	portService "momo/service/port"
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
