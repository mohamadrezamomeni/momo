package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"momo/pkg/config"
	momoLog "momo/pkg/log"
	"momo/repository/migrate"
	"momo/repository/sqllite"
	"momo/scheduler"

	_ "github.com/mattn/go-sqlite3"

	hostManagerSqlite "momo/repository/sqllite/host_manager"
	inboundSqlite "momo/repository/sqllite/inbound"
	userSqlite "momo/repository/sqllite/user"
	vpnSqlite "momo/repository/sqllite/vpn_manager"

	hostService "momo/service/host"
	inboundService "momo/service/inbound"
	userService "momo/service/user"
	vpnService "momo/service/vpn_manager"
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
	_, _, _, inboundSvc := getServices(&cfg.DB)

	done := make(chan struct{})

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		scheduler := scheduler.New(inboundSvc)
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

	hostSvc := hostService.New(hostRepo)
	vpnSvc := vpnService.New(vpnRepo)
	userSvc := userService.New(userRepo)

	inbouncSvc := inboundService.New(inboundRepo, vpnSvc, userSvc, hostSvc)
	return hostSvc, vpnSvc, userSvc, inbouncSvc
}
