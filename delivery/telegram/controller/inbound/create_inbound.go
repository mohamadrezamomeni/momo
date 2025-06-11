package inbound

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) SetStateCreateInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		idStr, err := core.GetID(update)
		if err != nil {
			return nil, err
		}
		_, isExist := telegramState.FindState(idStr)
		if !isExist {
			telegramState.NewState(idStr,
				"ask_selecting_VPNSource",
				"answer_VPNSource",
				"ask_selecting_VPN",
				"answer_selecting_VPN",
				"ask_selecting_VPNPackage",
				"answer_selecting_VPNPackage",
				"create_inbound",
			)
			return nil, nil
		}

		return next(update)
	}
}

func (h *Handler) CreateInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	vpnType, vpnPackage, vpnSource, err := getCreatingInboundData(update)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_, err = h.inboundSvc.Create(&inboundServiceDto.CreateInbound{
		ServerType:   entity.High,
		UserID:       update.UserSystem.ID,
		VPNType:      vpnType,
		TrafficLimit: vpnPackage.TrafficLimit,
		Start:        now,
		Country:      vpnSource.Country,
		End:          now.AddDate(0, int(vpnPackage.Months), int(vpnPackage.Days)),
	})
	if err != nil {
		return nil, err
	}
	text, err := telegrammessages.GetMessage("inbound.create.successfully_creation", map[string]string{})
	if err != nil {
		return nil, err
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}
	msgConfig := tgbotapi.NewMessage(id, text)

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		ReleaseState:  true,
		RedirectRoot:  true,
	}, nil
}

func getCreatingInboundData(update *core.Update) (entity.VPNType, *entity.VPNPackage, *entity.VPNSource, error) {
	scope := "controller.createVPN.getData"
	idStr, err := core.GetID(update)
	if err != nil {
		return 0, nil, nil, err
	}

	state, isExist := telegramState.FindState(idStr)
	if !isExist {
		return 0, nil, nil, momoError.Scope(scope).ErrorWrite()
	}

	vpnTypeVal, isExist := state.GetData("vpn_type")
	if !isExist {
		return 0, nil, nil, momoError.Scope(scope).ErrorWrite()
	}

	vpnType, ok := vpnTypeVal.(entity.VPNType)
	if !ok {
		return 0, nil, nil, momoError.Scope(scope).ErrorWrite()
	}

	vpnPackageVal, isExist := state.GetData("vpn_package")
	if !isExist {
		return 0, nil, nil, momoError.Scope(scope).ErrorWrite()
	}

	vpnPackage, ok := vpnPackageVal.(*entity.VPNPackage)
	if !ok {
		return 0, nil, nil, momoError.Scope(scope).ErrorWrite()
	}

	vpnSourceVal, isExist := state.GetData("vpn_source")
	if !isExist {
		return 0, nil, nil, momoError.Scope(scope).ErrorWrite()
	}

	vpnSource, ok := vpnSourceVal.(*entity.VPNSource)
	if !ok {
		return 0, nil, nil, momoError.Scope(scope).ErrorWrite()
	}
	return vpnType, vpnPackage, vpnSource, nil
}
