package main

import (
	"log"

	"github.com/mohamadrezamomeni/momo/pkg/config"
	momoLog "github.com/mohamadrezamomeni/momo/pkg/log"
	"github.com/mohamadrezamomeni/momo/repository/migrate"

	_ "github.com/mattn/go-sqlite3"

	serviceInitializer "github.com/mohamadrezamomeni/momo/pkg/service"
	authValidator "github.com/mohamadrezamomeni/momo/validator/auth"
	userValidator "github.com/mohamadrezamomeni/momo/validator/user"

	httpserver "github.com/mohamadrezamomeni/momo/delivery/http_server"
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

	_, _, userSvc, _, authSvc, cryptSvc := serviceInitializer.GetServices(&cfg)

	server := httpserver.New(&cfg.HTTP, authSvc, userSvc, cryptSvc, userValidator.New(), authValidator.New())

	server.Serve()
}
