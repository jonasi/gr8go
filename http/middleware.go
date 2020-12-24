package gr8http

func NopMiddleware(method, path string, h Handler) Handler {
	return h
}

// Middleware takes an http route and spits out a new handler
type Middleware interface {
	ID() string
	Handler(method, path string, h Handler) Handler
}

type MiddlewareFilterer interface {
	Middleware
	Filter([]Middleware) []Middleware
}

type mw struct {
	id string
	mk func(string, string, Handler) Handler
}

func (m *mw) ID() string                                     { return m.id }
func (m *mw) Handler(method, path string, h Handler) Handler { return m.mk(method, path, h) }

type mwf struct {
	mw
	filter func([]Middleware) []Middleware
}

func (m *mwf) Filter(mw []Middleware) []Middleware {
	return m.filter(mw)
}

// MiddlewareFunc is a helper for defining middleware
func MiddlewareFunc(id string, fn func(method, path string, h Handler) Handler) Middleware {
	return &mw{id: id, mk: fn}
}

func MiddlewareFilterFunc(id string, fn func([]Middleware) []Middleware) Middleware {
	return &mwf{mw: mw{id: id, mk: NopMiddleware}, filter: fn}
}

// SkipMiddleware returns a MiddlewareFilter to skip specific middleware
func SkipMiddleware(smws ...Middleware) Middleware {
	id := "__skip__"
	for i, mw := range smws {
		id += mw.ID()
		if i != len(smws)-1 {
			id += "_"
		}
	}

	return MiddlewareFilterFunc(id, func(mws []Middleware) []Middleware {
		filtered := []Middleware{}

		for _, m := range mws {
			add := true
			for _, sm := range smws {
				if sm.ID() == m.ID() {
					add = false
					break
				}
			}

			if add {
				filtered = append(filtered, m)
			}
		}

		return filtered
	})
}
