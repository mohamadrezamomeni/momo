package xray

import (
	"fmt"

	"momo/proxy/xray/dto"
	"momo/proxy/xray/serializer"

	"momo/pkg/log"

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
	RemoveUser()
}

func New(cfg XrayConfig, logger log.ILog) (IXray, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Address, cfg.ApiPort), grpc.WithInsecure())
	if err != nil {
		return &Xray{}, err
	}

	if err != nil {
		logger.WriteWarrning("xray isn't accessable \n - check configuration")
		return &Xray{}, fmt.Errorf("xray isnt accessable")
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
