package httputil

import "net/http"

type Route struct {
	Method     string
	Path       string
	Handler    APIHandler
	Middleware []func(http.Handler) http.Handler
}

// Register mounts each route on mux. A route's middleware is applied in
// slice order, so the last entry is the outermost and runs first.
func Register(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		var handler http.Handler = route.Handler
		for _, mw := range route.Middleware {
			handler = mw(handler)
		}
		mux.Handle(route.Method+" "+route.Path, handler)
	}
}
