package inbound

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

type PackageOperation = int

const (
	askPackage PackageOperation = iota + 100
	answerPackage
)

var packages []*entity.VPNPackage = []*entity.VPNPackage{
	{
		ID:           "1",
		TrafficLimit: 50000,
		Months:       1,
		Days:         0,
		Price:        20000,
	},
}

func (h *Handler) getResponseAskPackage(user *entity.User) (*core.ResponseHandlerFunc, error) {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, pkg := range packages {
		button, err := h.getPackageButton(pkg)
		if err != nil {
			return nil, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(*button))
	}

	id, err := utils.ConvertToInt64(user.TelegramID)
	if err != nil {
		return nil, err
	}

	askingPackageTitle, err := telegrammessages.GetMessage("inbound.extend.ask_package", map[string]string{})
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(id, askingPackageTitle)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &core.ResponseHandlerFunc{
		Result: msg,
	}, nil
}

func (h *Handler) getPackageButton(pkg *entity.VPNPackage) (*tgbotapi.InlineKeyboardButton, error) {
	var titleDuration string
	var err error

	if pkg.Days > 1 {
		titleDuration, err = telegrammessages.GetMessage("inbound.extend.many_days", map[string]string{
			"count": strconv.Itoa(int(pkg.Days)),
		})
	} else if pkg.Days == 0 {
		titleDuration, err = telegrammessages.GetMessage("inbound.extend.one_day", map[string]string{})
	} else if pkg.Months > 1 {
		titleDuration, err = telegrammessages.GetMessage("inbound.extend.many_months", map[string]string{
			"count": strconv.Itoa(int(pkg.Months)),
		})
	} else {
		titleDuration, err = telegrammessages.GetMessage("inbound.extend.one_month", map[string]string{})
	}

	if err != nil {
		return nil, err
	}

	titlePkg, err := telegrammessages.GetMessage("inbound.extend.package_buttom", map[string]string{
		"timeDuration": titleDuration,
		"price":        strconv.Itoa(int(pkg.Price)),
	})
	if err != nil {
		return nil, err
	}
	button := tgbotapi.NewInlineKeyboardButtonData(titlePkg, pkg.ID)
	return &button, nil
}

func (h *Handler) answerPackage(packageID string) (*entity.VPNPackage, error) {
	scope := "telegram.controller.answerPackage"
	var pkgSelected *entity.VPNPackage

	for _, pkg := range packages {
		if pkg.ID == packageID {
			pkgSelected = pkg
			break
		}
	}

	if pkgSelected == nil {
		return nil, momoError.Scope(scope).NotFound().ErrorWrite()
	}

	return pkgSelected, nil
}
