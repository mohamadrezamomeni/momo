package main

import (
	"log"

	telegram "github.com/mohamadrezamomeni/momo/delivery/telegram"
	config "github.com/mohamadrezamomeni/momo/pkg/config"
	momoLog "github.com/mohamadrezamomeni/momo/pkg/log"
	"github.com/mohamadrezamomeni/momo/repository/migrate"

	serviceInitializer "github.com/mohamadrezamomeni/momo/pkg/service"
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

	_, _, userSvc, _, authSvc, _ := serviceInitializer.GetServices(&cfg)

	t := telegram.New(&cfg.TelegramConfig, userSvc, authSvc)
	t.Serve()
}
