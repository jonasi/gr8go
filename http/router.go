package gr8http

import (
	"context"
	"net/http"
	"sort"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
)

func NewRouter() *Router {
	r := httprouter.New()
	r.HandleMethodNotAllowed = false

	return &Router{
		router: r,
	}
}

type Router struct {
	router             *httprouter.Router
	routes             []*Route
	middleware         []Middleware
	notFound           Handler
	notFoundMiddleware []Middleware
	started            int32
}

// Handle registers the provided route with the router
func (r *Router) Handle(rts ...*Route) {
	if atomic.LoadInt32(&r.started) == 1 {
		panic("Attempting to register routes after the server has started")
	}

	r.routes = append(r.routes, rts...)
}

// Lookup finds the associated route that was registered with the provided
// method and path
func (r *Router) Lookup(method, path string) *Route {
	if atomic.LoadInt32(&r.started) == 1 {
		panic("Attempting to lookup routes after the server has started")
	}

	for _, rt := range r.routes {
		if rt.Method == method && rt.Path == path {
			return rt
		}
	}

	return nil
}

// HandleNotFound registers a handler to run when a method is not found
func (r *Router) HandleNotFound(h Handler, mw ...Middleware) {
	if atomic.LoadInt32(&r.started) == 1 {
		panic("Attempting to register routes after the server has started")
	}

	r.notFound = h
	r.notFoundMiddleware = mw
}

// NotFoundHandler returns the registered not found handler
func (r *Router) NotFoundHandler() Handler {
	if r.notFound != nil {
		return r.notFound
	}

	return FromHttpHandler(http.NotFoundHandler())
}

// AddMiddleware adds global middleware to the server
func (r *Router) AddMiddleware(mw ...Middleware) {
	if atomic.LoadInt32(&r.started) == 1 {
		panic("Attempting to register global middleware after the server has started")
	}

	r.middleware = append(r.middleware, mw...)
}

func (r *Router) MakeHttpHandler(ctx context.Context) http.Handler {
	atomic.StoreInt32(&r.started, 1)

	r.initRoutes(ctx)
	return r.router
}

var allMethods = []string{http.MethodPut, http.MethodGet, http.MethodPost, http.MethodHead, http.MethodTrace, http.MethodPatch, http.MethodDelete, http.MethodOptions, http.MethodConnect}

func (r *Router) initRoutes(ctx context.Context) {
	sort.Sort(routes(r.routes))
	for _, rt := range r.routes {
		h, mws := r.applyMiddleware(rt.Method, rt.Path, rt.Handler, rt.Middleware)
		logRouterHandlingRoute(ctx, rt.Method, rt.Path, mws)

		if rt.Method == "*" {
			for _, method := range allMethods {
				r.router.Handler(method, rt.Path, h)
			}
		} else {
			r.router.Handler(rt.Method, rt.Path, h)
		}
	}

	nf, mws := r.applyMiddleware("", "", r.NotFoundHandler(), r.notFoundMiddleware)
	logRouterHandlingNotFound(ctx, mws)
	r.router.NotFound = nf
}

func (r *Router) applyMiddleware(method string, path string, h Handler, mw []Middleware) (http.Handler, []string) {
	var (
		ids = []string{}
		all = append(append([]Middleware{}, r.middleware...), mw...)
	)

	// reverse
	for i := len(all)/2 - 1; i >= 0; i-- {
		opp := len(all) - 1 - i
		all[i], all[opp] = all[opp], all[i]
	}

	for len(all) > 0 {
		mw := all[0]
		all = all[1:]

		// check if the middleware wants to modify the
		// rest of the middleware list
		if filter, ok := mw.(MiddlewareFilterer); ok {
			all = filter.Filter(all)
		}

		h = mw.Handler(method, path, h)
		ids = append(ids, mw.ID())
	}

	return ToHttpHandler(h), ids
}
