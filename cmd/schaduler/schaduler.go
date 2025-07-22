package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/mohamadrezamomeni/momo/pkg/config"
	momoLog "github.com/mohamadrezamomeni/momo/pkg/log"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/scheduler"

	_ "github.com/mattn/go-sqlite3"

	serviceInitializer "github.com/mohamadrezamomeni/momo/pkg/service"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ERROR: somthing went wrong with loding config \n - you can check existance of config \n - you can see content of config")
	}

	momoLog.Init(cfg.Log)
	telegrammessages.Load()
	migration := migrate.New(&cfg.DB)

	migration.UP()

	hostSvc, vpnSvc, _, _, _, _, _, _, _, healingUpInbound, hostInboundSvc, inboundTrafficSvc, _ := serviceInitializer.GetServices(&cfg)

	done := make(chan struct{})

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		scheduler := scheduler.New(healingUpInbound, inboundTrafficSvc, hostInboundSvc, vpnSvc, hostSvc)
		scheduler.Start(done, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	done <- struct{}{}
	wg.Wait()
}
