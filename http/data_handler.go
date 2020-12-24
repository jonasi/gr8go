package gr8http

import (
	"context"
	"net/http"
)

type DataFormatter interface {
	Format(http.ResponseWriter, *http.Request, interface{}) error
}

type dataFormatterCtx struct{}

func WithDataFormatter(ctx context.Context, handler DataFormatter) context.Context {
	return context.WithValue(ctx, dataFormatterCtx{}, handler)
}

func GetDataFormatter(ctx context.Context) DataFormatter {
	v := ctx.Value(dataFormatterCtx{})
	dh, _ := v.(DataFormatter)
	return dh
}

func Success(value interface{}) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) Handler {
		w.WriteHeader(200)
		GetDataFormatter(r.Context()).Format(w, r, value)
		return nil
	})
}
