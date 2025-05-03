package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/pkg/config"
	momoLog "github.com/mohamadrezamomeni/momo/pkg/log"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
	"github.com/mohamadrezamomeni/momo/scheduler"

	_ "github.com/mattn/go-sqlite3"

	hostManagerSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/host_manager"
	inboundSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/inbound"
	userSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/user"
	vpnSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/vpn_manager"

	hostService "github.com/mohamadrezamomeni/momo/service/host"
	inboundService "github.com/mohamadrezamomeni/momo/service/inbound"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	vpnService "github.com/mohamadrezamomeni/momo/service/vpn_manager"
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
	_, vpnSvc, _, inboundSvc := getServices(&cfg.DB)

	done := make(chan struct{})

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		scheduler := scheduler.New(inboundSvc, vpnSvc)
		scheduler.Start(done, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	done <- struct{}{}
	wg.Wait()
}

func getServices(cfg *sqllite.DBConfig) (
	*hostService.Host,
	*vpnService.VPNService,
	*userService.User,
	*inboundService.Inbound,
) {
	db := sqllite.New(cfg)
	userRepo := userSqlite.New(db)
	inboundRepo := inboundSqlite.New(db)
	hostRepo := hostManagerSqlite.New(db)
	vpnRepo := vpnSqlite.New(db)

	hostSvc := hostService.New(hostRepo, adapter.AdaptedWorkerFactory)
	vpnSvc := vpnService.New(vpnRepo, adapter.AdaptedVPNProxyFactory)
	userSvc := userService.New(userRepo)

	inbouncSvc := inboundService.New(inboundRepo, vpnSvc, userSvc, hostSvc)
	return hostSvc, vpnSvc, userSvc, inbouncSvc
}
