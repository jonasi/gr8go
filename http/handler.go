package gr8http

import (
	"context"
	"net/http"
)

type Handler interface {
	Handle(http.ResponseWriter, *http.Request) Handler
	WithContext(context.Context) Handler
}

func HandlerFunc(fn func(http.ResponseWriter, *http.Request) Handler) Handler {
	return handlerFunc{nil, fn}
}

type handlerFunc struct {
	ctx context.Context
	fn  func(http.ResponseWriter, *http.Request) Handler
}

func (h handlerFunc) WithContext(ctx context.Context) Handler {
	return handlerFunc{ctx, h.fn}
}

func (h handlerFunc) Handle(w http.ResponseWriter, r *http.Request) Handler {
	if h.ctx != nil {
		r = r.WithContext(h.ctx)
	}

	return h.fn(w, r)
}

type httpHandler struct {
	ctx context.Context
	h   http.Handler
}

func (h httpHandler) WithContext(ctx context.Context) Handler {
	return httpHandler{ctx, h.h}
}

func (h httpHandler) Handle(w http.ResponseWriter, r *http.Request) Handler {
	if h.ctx != nil {
		r = r.WithContext(h.ctx)
	}

	h.h.ServeHTTP(w, r)
	return nil
}

func FromHttpHandler(h http.Handler) Handler {
	return httpHandler{nil, h}
}

func ToHttpHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeHttp(h, w, r)
	})
}

func ServeHttp(h Handler, w http.ResponseWriter, r *http.Request) {
	for h != nil {
		h = h.Handle(w, r)
	}
}

func CodeHandler(statusCode int) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) Handler {
		w.WriteHeader(statusCode)
		return nil
	})
}

func ErrorCodeHandler(statusCode int) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) Handler {
		http.Error(w, http.StatusText(statusCode), statusCode)
		return nil
	})
}
