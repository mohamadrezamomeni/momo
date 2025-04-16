package xray

import (
	"fmt"

	"momo/proxy/xray/dto"
	"momo/proxy/xray/serializer"

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
	HsClient   handlerService.HandlerServiceClient
	SsClient   statsService.StatsServiceClient
	LsClient   loggerService.LoggerServiceClient
}

type IXray interface {
	AddInbound(*dto.AddInbound) (*serializer.AddInboundSerializer, error)
	RemoveInbound(*dto.RemoveInbound) (*serializer.RemoveInbound, error)
	ReceiveInboundTraffic(*dto.ReceiveInboundTraffic) (*serializer.ReceiveInboundTraffic, error)
	AddUser(*dto.AddUser) error
	RemoveUser(*dto.RemoveUser) error
}

func New(cfg XrayConfig) (IXray, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Address, cfg.ApiPort), grpc.WithInsecure())
	if err != nil {
		return &Xray{}, momoError.Error("xray isnt accessable please check configuration")
	}

	return &Xray{
		HsClient: handlerService.NewHandlerServiceClient(conn),
		SsClient: statsService.NewStatsServiceClient(conn),
		LsClient: loggerService.NewLoggerServiceClient(conn),

		address:    cfg.Address,
		apiPort:    cfg.ApiPort,
		configPath: cfg.configPath,
	}, nil
}
