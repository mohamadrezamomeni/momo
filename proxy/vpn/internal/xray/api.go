package xray

import (
	"fmt"

	momoError "momo/pkg/error"

	vpnDto "momo/proxy/vpn/dto"
	"momo/proxy/vpn/internal/xray/dto"
	vpnSerializer "momo/proxy/vpn/serializer"

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

func (x *Xray) Add(inpt *vpnDto.Inbound) error {
	_, err := x.addInbound(&dto.AddInbound{
		Port:     inpt.Port,
		Tag:      inpt.Tag,
		Protocol: inpt.Protocol,
		User: &dto.InboundUser{
			Username: inpt.User.Username,
			Level:    inpt.User.Level,
			UUID:     inpt.User.ID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (x *Xray) Disable(inpt *vpnDto.Inbound) error {
	_, err := x.removeInbound(&dto.RemoveInbound{
		Tag: inpt.Tag,
	})
	if err != nil {
		return err
	}
	err = x.resetTraffic(inpt.Tag)
	if err != nil {
		return err
	}
	return nil
}

func (x *Xray) GetTraffic(inpt *vpnDto.Inbound) (*vpnSerializer.Traffic, error) {
	data, err := x.getInboundTrafficWithoutBeigReseted(inpt.Tag)
	if err != nil {
		return &vpnSerializer.Traffic{}, err
	}
	return &vpnSerializer.Traffic{
		Download: int(data.DownLink),
		Upload:   int(data.UpLink),
	}, nil
}

func (x *Xray) DoesExist(inpt *vpnDto.Inbound) (bool, error) {
	data, err := x.getUsers(inpt.Tag)
	if err == nil && len(data.Usernames) > 0 {
		return true, nil
	}
	return false, err
}

func (x *Xray) Test() error {
	return x.fakeReceiveInboundTraffic()
}
