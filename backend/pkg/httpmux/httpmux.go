// Package httpmux is a thin wrapper around http.ServeMux with support for
// middleware.
package httpmux

import "net/http"

// Middleware is a handler function designed to run code before and/or after
// another Handler.
type Middleware func(http.Handler) http.Handler

type Mux struct {
	mw  []Middleware
	mux *http.ServeMux
}

// New creates a new Mux. Middleware are executed in the order they are
// provided and before any handler-specific middleware.
func New(mw ...Middleware) *Mux {
	return &Mux{
		mw:  mw,
		mux: http.NewServeMux(),
	}
}

// ServeHTTP implements the http.Handler interface.
func (m *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.mux.ServeHTTP(w, req)
}

// HandleFunc registers the handler function for the given pattern. Middleware
// are executed in the order they are provided and after any global middleware.
func (m *Mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), mw ...Middleware) {
	m.Handle(pattern, http.HandlerFunc(handler), mw...)
}

// Handle registers a new handler with given path pattern. Middleware are
// executed in the order they are provided and after any global middleware.
func (m *Mux) Handle(pattern string, h http.Handler, mw ...Middleware) {
	h = wrapMiddleware(mw, h)
	h = wrapMiddleware(m.mw, h)
	m.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	})
}

func wrapMiddleware(mw []Middleware, h http.Handler) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		wrap := mw[i]
		if wrap != nil {
			h = wrap(h)
		}
	}
	return h
}
