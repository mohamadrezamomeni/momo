package xray

import (
	"fmt"

	"momo/proxy/xray/adapter"
	"momo/proxy/xray/dto"
	"momo/proxy/xray/serializer"

	"momo/pkg/log"
)

type Xray struct {
	address    string
	apiPort    string
	configPath string
	xSDK       *adapter.XraySDK
}

type IXray interface {
	AddInbound(*dto.AddInbound) (*serializer.AddInboundSerializer, error)
	RemoveInbound()
	QueryUser()
	QueryInbound()
	AddUser()
	RemoveUser()
}

func New(cfg XrayConfig, logger log.ILog) (IXray, error) {
	xSDK, err := adapter.New(&adapter.BaseConfig{
		APIAddress: cfg.Address,
		APIPort:    cfg.ApiPort,
	})
	if err != nil {
		logger.WriteWarrning("xray isn't accessable \n - check configuration")
		return &Xray{}, fmt.Errorf("xray isnt accessable")
	}

	return &Xray{
		xSDK:       xSDK,
		address:    cfg.Address,
		apiPort:    cfg.ApiPort,
		configPath: cfg.configPath,
	}, nil
}
