package inbound

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) ListInbounds(update *core.Update) (*core.ResponseHandlerFunc, error) {
	idStr, err := core.GetID(update)
	if err != nil {
		return nil, err
	}

	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, err
	}

	inbounds, err := h.inboundSvc.Filter(&inboundServiceDto.FilterInbounds{
		UserID: update.UserSystem.ID,
	})
	if err != nil {
		return nil, err
	}

	if len(inbounds) == 0 {
		return h.sendNotFoundInbounds(update.UserSystem)
	}

	var sb strings.Builder

	title, err := telegrammessages.GetMessage("inbound.list.inbounds_title", map[string]string{})
	if err != nil {
		return nil, err
	}

	sb.WriteString(title)

	for i, inbound := range inbounds {
		item, err := h.writeItem(i, inbound)
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

func (h *Handler) getStatus() (string, string, string, string, error) {
	active, err := telegrammessages.GetMessage("inbound.list.active", map[string]string{})
	if err != nil {
		return "", "", "", "", err
	}

	deactive, err := telegrammessages.GetMessage("inbound.list.deactive", map[string]string{})
	if err != nil {
		return "", "", "", "", err
	}

	block, err := telegrammessages.GetMessage("inbound.list.block", map[string]string{})
	if err != nil {
		return "", "", "", "", err
	}

	unblock, err := telegrammessages.GetMessage("inbound.list.unblock", map[string]string{})
	if err != nil {
		return "", "", "", "", err
	}

	return active, deactive, block, unblock, nil
}

func (h *Handler) writeItem(i int, inbound *entity.Inbound) (string, error) {
	var sb strings.Builder

	active, deactive, block, unblock, err := h.getStatus()

	idReport, err := telegrammessages.GetMessage("inbound.list.id_report", map[string]string{
		"counter": strconv.Itoa(i + 1),
		"id":      strconv.Itoa(inbound.ID),
	})
	if err != nil {
		return "", err
	}
	trafficReport, err := telegrammessages.GetMessage("inbound.list.inbound_traffic_used", map[string]string{
		"traffic_used":  strconv.Itoa(int(inbound.TrafficUsage)),
		"traffic_limit": strconv.Itoa(int(inbound.TrafficLimit)),
	})
	if err != nil {
		return "", err
	}

	sb.WriteString(idReport)
	sb.WriteString(trafficReport)
	if inbound.IsActive {
		sb.WriteString(active)
	} else {
		sb.WriteString(deactive)
	}
	sb.WriteString(" | ")
	if inbound.IsBlock {
		sb.WriteString(block)
	} else {
		sb.WriteString(unblock)
	}
	sb.WriteString("\n\n")
	return sb.String(), nil
}
