package ober

import (
	"net/http"
)

// Main strcut used to configure and access functionality
type Middleware struct {
	mux      http.Handler
	handlers []http.HandlerFunc
}

// Inits a new middleware chain.
func (m *Middleware) New() *Middleware {
	return &Middleware{handlers: make([]http.HandlerFunc, 0, 0)}
}

// Add middleware to handlers slice of HandlerFunc
func (m *Middleware) Add(h ...http.HandlerFunc) {
	m.handlers = append(m.handlers, h...)
}

// AddMux adds our mux to run.
func (m *Middleware) AddMux(mux http.Handler) {
	m.mux = mux
}

// So we can satisfy the http.Handler interface.
// func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	for _, h := range m.handlers {
// 		h.ServeHTTP(w, r)
// 	}
//
// 	m.mux.ServeHTTP(w, r)
// }
