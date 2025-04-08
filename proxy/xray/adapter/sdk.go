package adapter

import (
	"fmt"

	loggerService "github.com/xtls/xray-core/app/log/command"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
)

type BaseConfig struct {
	APIAddress string
	APIPort    string
}

type UserInfo struct {
	Uuid       string
	Level      uint32
	InTag      string
	Email      string
	CipherType string
	Password   string
}

type XraySDK struct {
	HsClient handlerService.HandlerServiceClient
	SsClient statsService.StatsServiceClient
	LsClient loggerService.LoggerServiceClient
}

func New(cfg *BaseConfig) (sdk *XraySDK, err error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.APIAddress, cfg.APIPort), grpc.WithInsecure())
	if err != nil {
		return &XraySDK{}, err
	}

	xraySDK := &XraySDK{
		HsClient: handlerService.NewHandlerServiceClient(conn),
		SsClient: statsService.NewStatsServiceClient(conn),
		LsClient: loggerService.NewLoggerServiceClient(conn),
	}

	return xraySDK, nil
}
