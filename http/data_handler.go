package gr8http

import (
	"context"
	"net/http"
)

type DataFormatter interface {
	Format(http.ResponseWriter, *http.Request, interface{}) error
	FormatError(http.ResponseWriter, *http.Request, error) error
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

type DataHandler func(http.ResponseWriter, *http.Request) (interface{}, error)

func (h DataHandler) Handle(w http.ResponseWriter, r *http.Request) Handler {
	formatter := GetDataFormatter(r.Context())
	v, err := h(w, r)
	if err != nil {
		formatter.FormatError(w, r, err)
	} else {
		formatter.Format(w, r, v)
	}

	return nil
}
