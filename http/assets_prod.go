// +build !dev

package gr8http

import (
	"net/http"
)

func (c SPAConf) indexHandler(assets http.FileSystem) (Handler, error) {
	return c.mkIndexHandler(assets)
}
