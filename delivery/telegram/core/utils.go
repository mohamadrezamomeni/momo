package core

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func GetID(update *Update) (string, error) {
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

func GetTelegramUser(update *Update) (*tgbotapi.User, error) {
	scope := "telegram.core.GetTelegramUser"
	var user *tgbotapi.User

	if update.Message != nil {
		user = update.Message.From
	} else if update.CallbackQuery != nil {
		user = update.CallbackQuery.From
	} else if update.MyChatMember != nil {
		user = &update.MyChatMember.From
	} else if update.EditedMessage != nil {
		user = update.EditedMessage.From
	} else if update.PollAnswer != nil {
		user = &update.PollAnswer.User
	} else if update.ChannelPost != nil {
		user = update.ChannelPost.From
	} else if update.InlineQuery != nil {
		user = update.InlineQuery.From
	}
	if user == nil {
		return nil, momoError.Scope(scope).Input(update).DebuggingError()
	}
	return user, nil
}
