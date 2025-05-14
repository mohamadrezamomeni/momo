package core

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func GetID(update *tgbotapi.Update) (string, error) {
	scope := "telegram.core.GetID"
	var id int64

	if update.Message != nil {
		id = update.Message.From.ID
	} else if update.CallbackQuery != nil {
		id = update.CallbackQuery.From.ID
	} else if update.MyChatMember != nil {
		id = update.MyChatMember.From.ID
	} else if update.EditedMessage != nil {
		id = update.EditedMessage.From.ID
	} else if update.PollAnswer != nil {
		id = update.PollAnswer.User.ID
	} else if update.ChannelPost != nil {
		id = update.ChannelPost.From.ID
	} else if update.InlineQuery != nil {
		id = update.InlineQuery.From.ID
	}

	if id == 0 {
		return "", momoError.Scope(scope).Input(update).DebuggingError()
	}
	return strconv.Itoa(int(id)), nil
}
