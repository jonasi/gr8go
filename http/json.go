package gr8http

import (
	"encoding/json"
	"net/http"

	"github.com/munnerz/goautoneg"
)

type jsonFormatter struct {
}

func (h *jsonFormatter) Format(w http.ResponseWriter, r *http.Request, value interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}

func (h *jsonFormatter) FormatError(w http.ResponseWriter, r *http.Request, err error) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func JSONMiddleware(method, path string, h Handler) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) Handler {
		acc := r.Header.Get("Accept")
		if acc != "" {
			ct := goautoneg.Negotiate(acc, []string{"application/json"})
			if ct == "" {
				return CodeHandler(http.StatusNotAcceptable)
			}
		}

		r = r.WithContext(WithDataFormatter(r.Context(), &jsonFormatter{}))
		return h.Handle(w, r)
	})
}
