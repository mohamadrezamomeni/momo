package charge

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	chargeServiceDto "github.com/mohamadrezamomeni/momo/dto/service/charge"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) SetStateChargeInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		idStr, err := core.GetID(update)
		if err != nil {
			return nil, err
		}
		_, isExist := telegramState.FindState(idStr)

		if !isExist {
			telegramState.NewState(idStr,
				"ask_selecting_inbound",
				"answer_selecting_inbound",
				"ask_selecting_VPNPackage",
				"answer_selecting_VPNPackage",
				"ask_detail_charge",
				"answer_detail_charge",
				"charge_inbound",
			)
			return nil, nil
		}

		return next(update)
	}
}

func (h *Handler) ChargeInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	inboundID, packageID, detail, err := h.getDataChargeInbound(update)
	if err != nil {
		return nil, err
	}
	_, err = h.chargeSvc.Create(&chargeServiceDto.CreateChargeDto{
		UserID:    update.UserSystem.ID,
		PackageID: packageID,
		InboundID: inboundID,
		Detail:    detail,
	})
	if err != nil {
		return nil, err
	}

	title, err := telegrammessages.GetMessage("charge.successfully_extending", map[string]string{})
	if err != nil {
		return nil, err
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(id, title)
	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		ReleaseState:  true,
		RedirectRoot:  true,
	}, nil
}

func (h *Handler) getDataChargeInbound(update *core.Update) (string, string, string, error) {
	scope := "telegram.charge.getDataChargeCharge"
	state, isExist := telegramState.FindState(update.UserSystem.TelegramID)
	if !isExist {
		return "", "", "", momoError.Scope(scope).ErrorWrite()
	}

	vpnPackageID, isExist := state.GetData("vpn_package_id")
	if !isExist {
		return "", vpnPackageID, "", momoError.Scope(scope).DebuggingError()
	}

	inboundID, isExist := state.GetData("inbound_id")
	if !isExist {
		return "", "", "", momoError.Scope(scope).DebuggingError()
	}

	detail, isExist := state.GetData("detail")
	if !isExist {
		return "", "", "", momoError.Scope(scope).DebuggingError()
	}

	return inboundID, vpnPackageID, detail, nil
}
