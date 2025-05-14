package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/pkg/cache"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
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

func (r *Router) getHandler(p string) HandlerFunc {
	for path, handler := range r.routing {
		if path == p {
			return handler
		}
	}
	return r.rootHandler
}

func (r *Router) Route(update *tgbotapi.Update) *ResponseHandlerFunc {
	var res *ResponseHandlerFunc
	var path string

	if update.CallbackQuery != nil {
		res, path = r.callbackQuery(update)
	}

	if update.Message != nil {
		res, path = r.message(update)
	}
	if update.MyChatMember != nil {
		return nil
	}

	key := r.getKey(update)
	if res != nil && !res.ReleaseState && len(path) > 0 {
		cache.Set(key, path)
	} else if res != nil && (res.ReleaseState || len(path) == 0) {
		cache.Delete(key)
	}

	if res == nil {
		res, _ = r.rootHandler(update)
	}
	return res
}

func (r *Router) callbackQuery(update *tgbotapi.Update) (*ResponseHandlerFunc, string) {
	text := update.CallbackQuery.Data
	return r.getResponse(text, update)
}

func (r *Router) message(update *tgbotapi.Update) (*ResponseHandlerFunc, string) {
	text := update.Message.Text
	return r.getResponse(text, update)
}

func (r *Router) getResponse(text string, update *tgbotapi.Update) (*ResponseHandlerFunc, string) {
	key := r.getKey(update)

	if r.isPath(text) {
		path := r.getPathFromText(text)
		return r.routeFromText(path, update), path
	}

	value, isExist := cache.Get(key)
	if !isExist {
		res, _ := r.rootHandler(update)
		return res, ""
	}

	path, ok := value.(string)

	if !ok {
		res, _ := r.rootHandler(update)
		return res, ""
	}

	handler := r.getHandler(path)

	res, err := handler(update)
	if err != nil {
		res, _ := r.rootHandler(update)
		return res, ""
	}

	return res, path
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

func (r *Router) routeFromText(path string, update *tgbotapi.Update) *ResponseHandlerFunc {
	handler := r.getHandler(path)
	res, err := handler(update)
	if err != nil {
		res, _ = r.rootHandler(update)
		return res
	}
	return res
}

func (r *Router) getKey(update *tgbotapi.Update) string {
	id, err := GetID(update)
	if err != nil {
		momoError.Wrap(err).Input(update).Fatal()
	}
	return id
}

func (r *Router) rootHandler(update *tgbotapi.Update) (*ResponseHandlerFunc, error) {
	handler := r.getHandler(r.defaultRoute)
	res, err := handler(update)
	return res, err
}
