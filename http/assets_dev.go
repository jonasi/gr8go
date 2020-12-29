// +build dev

package gr8http

import (
	"net/http"
)

func (c SPAConf) indexHandler(assets http.FileSystem) (Handler, error) {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) Handler {
		h, err := c.mkIndexHandler(assets)
		if err != nil {
			logSPAIndexHandlerError(r.Context(), err)
			return ErrorCodeHandler(http.StatusInternalServerError)
		}

		return h
	}), nil
}
