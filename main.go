package main

import (
	"log"

	"momo/pkg/config"
	momoLog "momo/pkg/log"
	"momo/proxy/xray"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding error \n - you can follow the problem in error log")
	}
	momoLogger := momoLog.New(cfg.Log)

	_, _ = xray.New(cfg.XrayConfig, momoLogger)
}
