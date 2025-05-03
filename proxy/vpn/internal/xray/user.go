package xray

import (
	"context"

	"github.com/mohamadrezamomeni/momo/pkg/utils"
	"github.com/mohamadrezamomeni/momo/proxy/vpn/internal/xray/dto"
	"github.com/mohamadrezamomeni/momo/proxy/vpn/internal/xray/serializer"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"

	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vmess"
)

func (x *Xray) addUser(inpt *dto.AddUser) error {
	scope := "xrayProxy.addUser"

	level, err := utils.ConvertToUint32(inpt.Level)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", *inpt)
	}
	_, err = x.hsClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inpt.Tag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Email: inpt.Username,
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
	scope := "xrayProxy.removeUser"

	_, err := x.hsClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inpt.Tag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: inpt.Username,
		}),
	})
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", *inpt)
	}
	return err
}

func (x *Xray) getUsers(tag string) (*serializer.GetUsers, error) {
	scope := "xrayProxy.getUsers"

	res, err := x.hsClient.GetInboundUsers(context.Background(), &command.GetInboundUserRequest{
		Tag: tag,
	})
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the tag is %s", tag)
	}
	usernames := make([]string, 0)
	for _, user := range res.Users {
		usernames = append(usernames, user.Email)
	}
	return &serializer.GetUsers{
		Usernames: usernames,
	}, nil
}
