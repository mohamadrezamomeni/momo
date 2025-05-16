package inbound

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

type TrafficItem struct {
	ID           string
	TrafficUsage int
	TrafficLimit int
}

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
		notFoundTitle, _ := telegrammessages.GetMessage("inbound.list.not_found", map[string]string{})
		return &core.ResponseHandlerFunc{
			Result:       tgbotapi.NewMessage(id, notFoundTitle),
			ReleaseState: true,
			RedirectRoot: true,
		}, nil
	}

	var sb strings.Builder

	title, err := telegrammessages.GetMessage("inbound.list.inbounds_title", map[string]string{})
	if err != nil {
		return nil, err
	}

	active, err := telegrammessages.GetMessage("inbound.list.active", map[string]string{})
	if err != nil {
		return nil, err
	}

	deactive, err := telegrammessages.GetMessage("inbound.list.deactive", map[string]string{})
	if err != nil {
		return nil, err
	}

	block, err := telegrammessages.GetMessage("inbound.list.block", map[string]string{})
	if err != nil {
		return nil, err
	}

	unblock, err := telegrammessages.GetMessage("inbound.list.unblock", map[string]string{})
	if err != nil {
		return nil, err
	}

	sb.WriteString(title)

	for i, inbound := range inbounds {
		idReport, err := telegrammessages.GetMessage("inbound.list.id_report", map[string]string{
			"counter": strconv.Itoa(i + 1),
			"id":      strconv.Itoa(inbound.ID),
		})
		if err != nil {
			return nil, err
		}
		trafficReport, err := telegrammessages.GetMessage("inbound.list.inbound_traffic_used", map[string]string{
			"traffic_used":  strconv.Itoa(int(inbound.TrafficUsage)),
			"traffic_limit": strconv.Itoa(int(inbound.TrafficLimit)),
		})
		if err != nil {
			return nil, err
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

	}
	message := sb.String()
	msg := tgbotapi.NewMessage(id, message)
	msg.ParseMode = "HTML"
	return &core.ResponseHandlerFunc{
		Result:       msg,
		ReleaseState: true,
		RedirectRoot: true,
	}, nil
}
