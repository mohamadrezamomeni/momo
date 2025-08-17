package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mohamadrezamomeni/momo/pkg/config"
	momoLog "github.com/mohamadrezamomeni/momo/pkg/log"
	"github.com/mohamadrezamomeni/momo/repository/migrate"

	_ "github.com/mattn/go-sqlite3"

	serviceInitializer "github.com/mohamadrezamomeni/momo/pkg/service"
	authValidator "github.com/mohamadrezamomeni/momo/validator/auth"
	chargeValidator "github.com/mohamadrezamomeni/momo/validator/charge"
	hostValidator "github.com/mohamadrezamomeni/momo/validator/host"
	inboundValidator "github.com/mohamadrezamomeni/momo/validator/inbound"
	userValidator "github.com/mohamadrezamomeni/momo/validator/user"
	vpnValidator "github.com/mohamadrezamomeni/momo/validator/vpn"
	vpnPackageValidator "github.com/mohamadrezamomeni/momo/validator/vpn_package"
	vpnSourceValidator "github.com/mohamadrezamomeni/momo/validator/vpn_source"

	httpserver "github.com/mohamadrezamomeni/momo/delivery/http_server"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	userService "github.com/mohamadrezamomeni/momo/service/user"

	userControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/user"
	"github.com/mohamadrezamomeni/momo/dto/service/user"
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
	hostSvc, vpnSvc, userSvc, inboundSvc, authSvc, cryptSvc, vpnPackageSvc, _, chargeSvc, _, _, _, vpnSourceSvc, inboundChargeSvc, tierSvc := serviceInitializer.GetServices(&cfg)

	initializer(userSvc, &cfg)

	server := httpserver.New(
		&cfg.HTTP,
		authSvc,
		userSvc,
		cryptSvc,
		hostSvc,
		vpnSvc,
		inboundSvc,
		inboundChargeSvc,
		chargeSvc,
		vpnPackageSvc,
		vpnSourceSvc,
		tierSvc,
		userValidator.New(),
		authValidator.New(),
		hostValidator.New(),
		vpnValidator.New(vpnSourceSvc),
		inboundValidator.New(userSvc, inboundSvc),
		vpnSourceValidator.New(),
		vpnPackageValidator.New(),
		chargeValidator.New(),
	)

	go func() {
		server.Serve()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	shutdown(userSvc, &cfg)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxWithTimeout); err != nil {
		momoLog.Info(err.Error())
	}
}

func initializer(userSvc *userService.User, cfg *config.Config) {
	createSuperAdmin(userSvc, &cfg.Admin)
}

func createSuperAdmin(userSvc *userService.User, cfg *config.Admin) {
	scope := "main.createUser"
	userValidator := userValidator.New()
	err := userValidator.ValidateAddUserRequest(userControllerDto.AddUser{
		IsAdmin:   true,
		Username:  cfg.Username,
		Password:  cfg.Password,
		FirstName: cfg.FirstName,
		LastName:  cfg.LastName,
	})
	if err != nil {
		momoError.Wrap(err).Scope(scope).Fatalf("please check admin config")
	}

	_, err = userSvc.CreateUserAdmin(&user.AddUser{
		IsAdmin:   true,
		Username:  cfg.Username,
		Password:  cfg.Password,
		FirstName: cfg.FirstName,
		LastName:  cfg.LastName,
	})
	if err != nil {
		momoError.Wrap(err).Scope(scope).Fatalf("please check admin config")
	}
}

func shutdown(userSvc *userService.User, cfg *config.Config) {
	deleteSuperAdmin(userSvc, &cfg.Admin)
}

func deleteSuperAdmin(userSvc *userService.User, cfg *config.Admin) {
	userSvc.DeleteByUsername(cfg.Username)
}
