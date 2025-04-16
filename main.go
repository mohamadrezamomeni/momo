package main

import (
	"log"

	"momo/pkg/config"
	momoLog "momo/pkg/log"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding error \n - you can follow the problem in error log")
	}
	momoLog.Init(cfg.Log)
}
