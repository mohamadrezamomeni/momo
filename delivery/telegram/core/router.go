package core

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/pkg/cache"
)

type Router struct {
	routing     map[string]HandlerFunc
	rootHandler HandlerFunc
}

func New(rootHandler HandlerFunc) *Router {
	return &Router{
		routing:     make(map[string]HandlerFunc),
		rootHandler: rootHandler,
	}
}

func (r *Router) Register(path string, h HandlerFunc, ms ...Middleware) {
	finalHandler := applyMiddleware(h, ms...)
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
	key := r.getKey(update.FromChat().ID)

	if update.CallbackQuery != nil {
		res, path = r.callbackQuery(update)
	}

	if update.Message != nil {
		res, path = r.message(update)
	}

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
	path := update.CallbackQuery.Data
	handler := r.getHandler(path)

	res, err := handler(update)
	if err != nil {
		res, _ := r.rootHandler(update)
		return res, ""
	}
	return res, path
}

func (r *Router) message(update *tgbotapi.Update) (*ResponseHandlerFunc, string) {
	text := update.Message.Text
	key := r.getKey(update.FromChat().ID)

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

func (r *Router) getKey(id int64) string {
	idStr := strconv.Itoa(int(id))
	return idStr
}
