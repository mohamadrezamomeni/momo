package main

import (
	"log"

	"momo/pkg/config"
)

var configPath = "config.yaml"

func main() {
	_, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding error \n - you can follow problem in error log")
	}
}
