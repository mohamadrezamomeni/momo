package xray

import (
	"fmt"

	momoError "momo/pkg/error"

	loggerService "github.com/xtls/xray-core/app/log/command"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
)

type Xray struct {
	address    string
	apiPort    string
	configPath string
	hsClient   handlerService.HandlerServiceClient
	ssClient   statsService.StatsServiceClient
	lsClient   loggerService.LoggerServiceClient
}

func New(cfg *XrayConfig) *Xray {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Address, cfg.ApiPort), grpc.WithInsecure())
	if err != nil {
		momoError.Fatalf("xray isnt accessable please check configuration")
	}

	return &Xray{
		hsClient: handlerService.NewHandlerServiceClient(conn),
		ssClient: statsService.NewStatsServiceClient(conn),
		lsClient: loggerService.NewLoggerServiceClient(conn),

		address: cfg.Address,
		apiPort: cfg.ApiPort,
	}
}

func (x *Xray) GetAddress() string {
	return x.address
}

func (x *Xray) Add() error {
	return nil
}

func (x *Xray) Disable() error {
	return nil
}

func (x *Xray) GetTraffic() {}
