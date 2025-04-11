package main

import (
	"log"

	"momo/pkg/config"
)

var configPath = "config.yaml"

func main() {
	_, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding config \n - you can check existance of config \n - you can see content of config")
	}
}
