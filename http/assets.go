package gr8http

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

func fixPrefix(prefix string) string {
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	if prefix[len(prefix)-1] != '/' {
		prefix = prefix + "/"
	}

	return prefix
}

// AssetsRoute returns a Route to handle static assets
func AssetsRoute(prefix string, fs http.FileSystem, mws ...Middleware) *Route {
	prefix = fixPrefix(prefix)

	return &Route{
		Method:     "GET",
		Path:       prefix + "*splat",
		Middleware: mws,
		Handler:    FromHttpHandler(http.StripPrefix(prefix, http.FileServer(fs))),
	}
}

// TemplateHandler returns an http.Handler that renders the provided template
func TemplateHandler(t *template.Template, fn func(*http.Request) interface{}) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) Handler {
		data := fn(r)
		if err := t.Execute(w, data); err != nil {
			logTemplateHandlerError(r.Context(), err)
			return ErrorCodeHandler(http.StatusInternalServerError)
		}

		return nil
	})
}

// SPAConf is a helper for defining routes and
// asset handling for single page apps
type SPAConf struct {
	Root              string
	IndexTemplate     *template.Template
	IndexTemplateData func(*http.Request, map[string]map[string]interface{}) interface{}
	IndexFilter       func(*http.Request) bool
	IndexMiddleware   []Middleware
	Assets            http.FileSystem
	AssetFile         string
	AssetPrefix       string
	AssetMiddleware   []Middleware
}

// Init initializes all the routes and confs for SPAConf
func SPA(r Router, conf SPAConf) error {
	indexHandler, err := conf.indexHandler(conf.Assets)
	if err != nil {
		return err
	}

	// use (abuse?) the not found mechanism to load the client
	// only do it for GET requests that pass the provided IndexFilter method
	oldNF := r.NotFoundHandler()
	r.HandleNotFound(HandlerFunc(func(w http.ResponseWriter, r *http.Request) Handler {
		if conf.IndexFilter != nil && r.Method == http.MethodGet && conf.IndexFilter(r) {
			var (
				h                  = oldNF
				isGet              = r.Method == http.MethodGet
				matchesRoot        = conf.Root == "" || strings.HasPrefix(r.URL.Path, conf.Root)
				matchesIndexFilter = conf.IndexFilter == nil || conf.IndexFilter(r)
			)

			if isGet && matchesRoot && matchesIndexFilter {
				h = indexHandler
			}

			return h
		}

		return nil
	}), conf.IndexMiddleware...)

	pref := fixPrefix(conf.AssetPrefix)
	if rt := r.Lookup("GET", pref+"*splat"); rt == nil {
		ar := AssetsRoute(pref, conf.Assets, conf.AssetMiddleware...)
		r.Handle(ar)
	}

	return nil
}

func (c SPAConf) mkIndexHandler(assets http.FileSystem) (Handler, error) {
	b, err := ReadFile(assets, c.AssetFile)
	if err != nil {
		return nil, err
	}
	var js map[string]map[string]interface{}
	if err := json.Unmarshal(b, &js); err != nil {
		return nil, err
	}

	return TemplateHandler(c.IndexTemplate, func(r *http.Request) interface{} {
		return c.IndexTemplateData(r, js)
	}), nil
}
