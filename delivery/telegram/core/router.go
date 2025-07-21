package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
)

type Router struct {
	routing      map[string]HandlerFunc
	defaultRoute string
}

func New(defaultRoute string) *Router {
	return &Router{
		routing:      make(map[string]HandlerFunc),
		defaultRoute: defaultRoute,
	}
}

func (r *Router) Register(path string, h HandlerFunc, ms ...Middleware) {
	scope := "telegram.core.registerroutes"

	finalHandler := applyMiddleware(h, ms...)
	if _, isExist := r.routing[path]; isExist {
		momoError.Scope(scope).UnExpected().Fatalf("you can't set duplicated route %s is set before", path)
	}
	r.routing[path] = finalHandler
}

func (r *Router) getHandler(path string) HandlerFunc {
	if handler, isExist := r.routing[path]; isExist {
		return handler
	}

	return r.RootHandler
}

func (r *Router) Route(update *Update) (*ResponseHandlerFunc, error) {
	var res *ResponseHandlerFunc
	var err error

	if update.CallbackQuery != nil {
		res, err = r.callbackQuery(update)
	}

	if update.Message != nil {
		res, err = r.message(update)
	}

	if update.MyChatMember != nil {
		return nil, nil
	}
	if res == nil {
		res, _ = r.RootHandler(update)
	}
	if res.MenuTab {
		res.MessageConfig.ReplyMarkup = r.enrichKeyboardMarkup(res.MessageConfig.ReplyMarkup)
	}

	return res, err
}

func (r *Router) enrichKeyboardMarkup(replyMarkup interface{}) interface{} {
	inlineKeyboardMarkup, ok := replyMarkup.(tgbotapi.InlineKeyboardMarkup)
	if !ok {
		return replyMarkup
	}

	menuButtonText, err := telegrammessages.GetMessage("root.menu_button", map[string]string{})
	if err != nil {
		menuButtonText = "menu"
	}

	menu := tgbotapi.NewInlineKeyboardButtonData(
		menuButtonText, "/menu",
	)
	row := tgbotapi.NewInlineKeyboardRow(menu)
	inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, row)
	return inlineKeyboardMarkup
}

func (r *Router) callbackQuery(update *Update) (*ResponseHandlerFunc, error) {
	text := update.CallbackQuery.Data
	return r.getResponse(text, update)
}

func (r *Router) message(update *Update) (*ResponseHandlerFunc, error) {
	text := update.Message.Text
	return r.getResponse(text, update)
}

func (r *Router) getResponse(text string, update *Update) (*ResponseHandlerFunc, error) {
	var res *ResponseHandlerFunc
	var err error
	id, err := GetID(update)
	if err != nil {
		return nil, err
	}

	if r.isPath(text) {
		telegramState.ResetState(id)
		path := r.getPathFromText(text)
		res, err = r.routeFromText(path, update)
	}

	maxLoop := 25
	for i := 0; err == nil && res == nil && i < maxLoop; i++ {
		res, err = r.getResponseFromState(update)
	}
	return res, err
}

func (r *Router) getResponseFromState(update *Update) (*ResponseHandlerFunc, error) {
	scope := "telegram.router.getResponseFromState"
	id, err := GetID(update)
	if err != nil {
		return nil, err
	}
	state, isExist := telegramState.FindState(id)
	if !isExist {
		return nil, momoError.Scope(scope).ErrorWrite()
	}

	if state.IsRequestCopeleted() {
		res, _ := r.RootHandler(update)
		return res, nil
	}

	path := state.GetPath()
	handler := r.getHandler(path)
	res, err := handler(update)
	if err != nil {
		res, _ := r.RootHandler(update)
		return res, err
	}

	if err != nil || state.IsRequestCopeleted() {
		state.ReleaseState()
	} else {
		state.Save()
	}

	state.Next()
	return res, nil
}

func (r *Router) isPath(text string) bool {
	action := byte('/')

	if text[0] == action {
		return true
	}

	return false
}

func (r *Router) getPathFromText(path string) string {
	return path[1:]
}

func (r *Router) routeFromText(path string, update *Update) (*ResponseHandlerFunc, error) {
	handler := r.getHandler(path)
	res, err := handler(update)
	if err != nil {
		res, _ = r.RootHandler(update)
		return res, err
	}
	return res, nil
}

func (r *Router) RootHandler(update *Update) (*ResponseHandlerFunc, error) {
	handler := r.getHandler(r.defaultRoute)
	res, err := handler(update)
	return res, err
}
