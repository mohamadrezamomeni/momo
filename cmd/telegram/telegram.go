package main

import (
	"log"

	telegram "github.com/mohamadrezamomeni/momo/delivery/telegram"
	config "github.com/mohamadrezamomeni/momo/pkg/config"
	momoLog "github.com/mohamadrezamomeni/momo/pkg/log"
	"github.com/mohamadrezamomeni/momo/repository/migrate"

	serviceInitializer "github.com/mohamadrezamomeni/momo/pkg/service"
	telegramMessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding config \n - you can check existance of config \n - you can see content of config")
	}

	momoLog.Init(cfg.Log)

	telegramMessages.Load()

	migration := migrate.New(&cfg.DB)

	migration.UP()

	_, _, userSvc, inboundSvc, authSvc, _, vpnPackage, _, chargeSvc, _, _, _, vpnSourceSvc, _ := serviceInitializer.GetServices(&cfg)

	t := telegram.New(&cfg.TelegramConfig,
		userSvc,
		authSvc,
		inboundSvc,
		vpnPackage,
		chargeSvc,
		vpnSourceSvc,
	)
	t.Serve()
}
