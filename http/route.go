package gr8http

// A Route is the tuple of (method, path, handler) that is registered
// with the router
type Route struct {
	Method     string
	Path       string
	Handler    Handler
	Middleware []Middleware
}

type routes []*Route

func (r routes) Len() int      { return len(r) }
func (r routes) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r routes) Less(i, j int) bool {
	if r[i].Path == r[j].Path {
		return r[i].Method < r[j].Method
	}
	return r[i].Path < r[j].Path
}
