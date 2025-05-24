package charge

import (
	"strconv"

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
				"ask_selecting_inbound",
				"answer_selecting_inbound",
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

func (h *Handler) CreateCharge(update *core.Update) (*core.ResponseHandlerFunc, error) {
	inbound, pkg, detail, err := getDataCreateCharge(update)
	if err != nil {
		return nil, err
	}
	_, err = h.chargeSvc.Create(&chargeServiceDto.CreateChargeDto{
		UserID:    update.UserSystem.ID,
		PackageID: pkg.ID,
		InboundID: strconv.Itoa(inbound.ID),
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

func getDataCreateCharge(update *core.Update) (*entity.Inbound, *entity.VPNPackage, string, error) {
	scope := "telegram.charge.getDataCreateCharge"
	state, isExist := telegramState.FindState(update.UserSystem.TelegramID)
	if !isExist {
		return nil, nil, "", momoError.Scope(scope).ErrorWrite()
	}

	val, isExist := state.GetData("vpn_package")
	if !isExist {
		return nil, nil, "", momoError.Scope(scope).DebuggingError()
	}

	vpnPackage, ok := val.(*entity.VPNPackage)
	if !ok {
		return nil, nil, "", momoError.Scope(scope).DebuggingError()
	}

	val, isExist = state.GetData("inbound")

	if !isExist {
		return nil, nil, "", momoError.Scope(scope).DebuggingError()
	}

	inbound, ok := val.(*entity.Inbound)
	if !ok {
		return nil, nil, "", momoError.Scope(scope).DebuggingError()
	}

	val, isExist = state.GetData("detail")
	if !isExist {
		return nil, nil, "", momoError.Scope(scope).DebuggingError()
	}

	detail, ok := val.(string)
	if !ok {
		return nil, nil, "", momoError.Scope(scope).DebuggingError()
	}

	return inbound, vpnPackage, detail, nil
}
