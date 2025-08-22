package charge

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	"github.com/mohamadrezamomeni/momo/dto/service/charge"
	"github.com/mohamadrezamomeni/momo/entity"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) FilterCharges(update *core.Update) (*core.ResponseHandlerFunc, error) {
	idStr, err := core.GetID(update)
	if err != nil {
		return nil, err
	}

	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, err
	}

	charges, err := h.chargeSvc.FilterCharges(&charge.FilterCharges{})
	if err != nil {
		return nil, err
	}

	if len(charges) == 0 {
		return h.sendNotFoundCharges(update.UserSystem)
	}

	var sb strings.Builder

	title, err := telegrammessages.GetMessage("charge.list.title", map[string]string{})
	if err != nil {
		return nil, err
	}
	sb.WriteString(title)

	for i, charge := range charges {
		item, err := h.writeItem(i, charge)
		if err != nil {
			return nil, err
		}
		sb.WriteString(item)
	}

	message := sb.String()
	msgConfig := tgbotapi.NewMessage(id, message)
	msgConfig.ParseMode = "HTML"
	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		ReleaseState:  true,
		RedirectRoot:  true,
	}, nil
}

func (h *Handler) writeItem(i int, charge *entity.Charge) (string, error) {
	var sb strings.Builder

	idReport, err := telegrammessages.GetMessage("charge.list.id", map[string]string{
		"counter": strconv.Itoa(i + 1),
		"id":      charge.ID,
	})
	if err != nil {
		return "", err
	}

	statusReport, err := telegrammessages.GetMessage("charge.list.status", map[string]string{
		"counter": strconv.Itoa(i + 1),
		"status":  entity.TranslateChargeStatus(charge.Status),
	})
	if err != nil {
		return "", err
	}

	detailReport, err := telegrammessages.GetMessage("charge.list.detail", map[string]string{
		"counter": strconv.Itoa(i + 1),
		"detail":  charge.Detail,
	})
	if err != nil {
		return "", err
	}

	sb.WriteString(idReport)
	sb.WriteString(statusReport)
	sb.WriteString(detailReport)

	sb.WriteString("\n\n")
	return sb.String(), nil
}

func (h *Handler) sendNotFoundCharges(user *entity.User) (*core.ResponseHandlerFunc, error) {
	id, err := utils.ConvertToInt64(user.TelegramID)
	if err != nil {
		return nil, err
	}
	inboundNotFoundText, _ := telegrammessages.GetMessage("charge.list.not_found", map[string]string{})

	msgConfig := tgbotapi.NewMessage(id, inboundNotFoundText)
	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		ReleaseState:  true,
		RedirectRoot:  true,
	}, nil
}
