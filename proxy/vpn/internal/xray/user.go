package xray

import (
	"context"

	"momo/pkg/utils"
	"momo/proxy/vpn/internal/xray/dto"

	momoError "momo/pkg/error"

	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vmess"
)

func (x *Xray) addUser(inpt *dto.AddUser) error {
	level, err := utils.ConvertToUint32(inpt.Level)
	if err != nil {
		return momoError.Error("user's level is wrong.")
	}
	_, err = x.hsClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
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

func (x *Xray) removeUser(inpt *dto.RemoveUser) error {
	_, err := x.hsClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inpt.Tag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: inpt.Email,
		}),
	})
	return err
}
