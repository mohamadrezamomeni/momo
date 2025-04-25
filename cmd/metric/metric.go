package main

import (
	"log"

	metricServer "momo/delivery/metric_server"
	"momo/pkg/config"
	momoLogger "momo/pkg/log"
	metricService "momo/service/metric"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding config \n - the problem was %v", err)
	}
	momoLogger.Init(cfg.Log)

	metricSrv := metricService.New(&cfg.MetricServer)

	server := metricServer.New(metricSrv, cfg.MetricServer)

	server.Start()
}
