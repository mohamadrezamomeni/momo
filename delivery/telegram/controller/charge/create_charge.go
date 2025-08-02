package charge

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	chargeServiceDto "github.com/mohamadrezamomeni/momo/dto/service/charge"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) SetStateCreateCharge(next core.HandlerFunc) core.HandlerFunc {
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
				"ask_detail_charge",
				"answer_detail_charge",
				"create_charge",
			)
			return nil, nil
		}

		return next(update)
	}
}

func (h *Handler) CreateInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	VPNSource, packageID, detail, VPNType, err := getDataCreateCharge(update)
	if err != nil {
		return nil, err
	}
	_, err = h.chargeSvc.Create(&chargeServiceDto.CreateChargeDto{
		UserID:    update.UserSystem.ID,
		PackageID: packageID,
		VPNSource: VPNSource,
		Detail:    detail,
		VPNType:   VPNType,
	})
	if err != nil {
		return nil, err
	}

	title, err := telegrammessages.GetMessage("charge.successfully_creating_vpn_request", map[string]string{})
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

func getDataCreateCharge(update *core.Update) (string, string, string, entity.VPNType, error) {
	scope := "telegram.charge.getDataCreateCharge"
	state, isExist := telegramState.FindState(update.UserSystem.TelegramID)
	if !isExist {
		return "", "", "", 0, momoError.Scope(scope).Input(update.UserSystem.TelegramID).Errorf("error to gat state")
	}

	vpnPackageID, isExist := state.GetData("vpn_package_id")
	if !isExist {
		return "", "", "", 0, momoError.Scope(scope).Input(state).DebuggingErrorf("error to get package")
	}

	VPNSource, isExist := state.GetData("vpn_source")
	if !isExist {
		return "", "", "", 0, momoError.Scope(scope).Input(state).DebuggingErrorf("error to get vpn source id")
	}

	vpnTypeStr, isExist := state.GetData("vpn_type")
	if !isExist {
		return "", "", "", 0, momoError.Scope(scope).DebuggingError()
	}
	VPNType := entity.ConvertStringVPNTypeToEnum(vpnTypeStr)
	if VPNType == entity.UknownVPNType || VPNType == 0 {
		return "", "", "", 0, momoError.Scope(scope).Input(state).DebuggingErrorf("error to get vpn type")
	}

	detail, isExist := state.GetData("detail")
	if !isExist {
		return "", "", "", 0, momoError.Scope(scope).Input(state).DebuggingErrorf("error to get detail")
	}

	return VPNSource, vpnPackageID, detail, VPNType, nil
}
