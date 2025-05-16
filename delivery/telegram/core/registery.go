package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ResponseHandlerFunc struct {
	Result       tgbotapi.MessageConfig
	ReleaseState bool
	RedirectRoot bool
}

type HandlerFunc = func(*Update) (*ResponseHandlerFunc, error)

type Middleware = func(HandlerFunc) HandlerFunc

func applyMiddleware(handler HandlerFunc, ms ...Middleware) HandlerFunc {
	for i := len(ms) - 1; i >= 0; i-- {
		handler = ms[i](handler)
	}
	return handler
}
