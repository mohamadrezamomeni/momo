package vpnpackage

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	vpnPackageServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_package"

	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) AskSelectingVPNPackage(update *core.Update) (*core.ResponseHandlerFunc, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	vpnPackages, err := h.vpnPackageSvc.Filter(&vpnPackageServiceDto.FilterVPNPackage{})
	if err != nil {
		return nil, err
	}
	if len(vpnPackages) == 0 {
		return h.notFoundAnyVPNPackages(update)
	}

	for _, pkg := range vpnPackages {
		button, err := h.getPackageButton(pkg)
		if err != nil {
			return nil, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(*button))
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}

	askingPackageTitle, err := telegrammessages.GetMessage("vpn_package.ask_selecting_package", map[string]string{})
	if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(id, askingPackageTitle)

	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		MenuTab:       true,
	}, nil
}

func (h *Handler) getPackageButton(pkg *entity.VPNPackage) (*tgbotapi.InlineKeyboardButton, error) {
	var titleDuration string
	var err error

	if pkg.Days > 1 {
		titleDuration, err = telegrammessages.GetMessage("vpn_package.many_days", map[string]string{
			"count": strconv.Itoa(int(pkg.Days)),
		})
	} else if pkg.Days == 0 {
		titleDuration, err = telegrammessages.GetMessage("vpn_package.one_day", map[string]string{})
	} else if pkg.Months > 1 {
		titleDuration, err = telegrammessages.GetMessage("vpn_package.many_months", map[string]string{
			"count": strconv.Itoa(int(pkg.Months)),
		})
	} else {
		titleDuration, err = telegrammessages.GetMessage("vpn_package.one_month", map[string]string{})
	}
	if err != nil {
		return nil, err
	}

	titlePkg, err := telegrammessages.GetMessage("vpn_package.package_buttom", map[string]string{
		"timeDuration": titleDuration,
		"price":        pkg.PriceTitle,
	})
	if err != nil {
		return nil, err
	}
	button := tgbotapi.NewInlineKeyboardButtonData(titlePkg, pkg.ID)
	return &button, nil
}

func (h *Handler) SelectVPNPackage(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.selectVPNPackage"
	packageID := update.CallbackQuery.Data
	pkg, err := h.vpnPackageSvc.FindVPNPackageByID(packageID)
	if err != nil {
		return nil, err
	}

	state, isExist := telegramState.FindState(update.UserSystem.TelegramID)

	if !isExist {
		return nil, momoError.Scope(scope).DebuggingError()
	}
	state.SetData("vpn_package", pkg)
	return nil, nil
}
