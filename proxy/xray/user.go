package xray

import (
	"context"
	"fmt"

	"momo/pkg/utils"
	"momo/proxy/xray/dto"

	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vmess"
)

func (x *Xray) AddUser(inpt *dto.AddUser) error {
	level, err := utils.ConvertToUint32(inpt.Level)
	if err != nil {
		return fmt.Errorf("user's level is wrong.")
	}
	_, err = x.HsClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inpt.Tag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Email: inpt.Email,
				Level: level,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id: inpt.UUID,
				}),
			},
		}),
	})
	return nil
}

func (x *Xray) RemoveUser() {
}
