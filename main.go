package main

import (
	"log"

	"momo/pkg/config"
	momoLog "momo/pkg/log"
	"momo/proxy/xray"
	"momo/proxy/xray/dto"
	"momo/repository/migrate"

	_ "github.com/mattn/go-sqlite3"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding config \n - you can check existance of config \n - you can see content of config")
	}
	momoLog.Init(cfg.Log)

	migration := migrate.New(&cfg.DB)

	migration.UP()

	d, _ := xray.New(cfg.XrayConfig)
	d.AddInbound(&dto.AddInbound{
		Tag:      "prrrr",
		Port:     "54321",
		Protocol: "vmess",
	})
}
