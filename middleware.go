package ober

import (
	"net/http"
)

// MW middleware struct used to configure and access functionality
type MW struct {
	handlers []http.Handler
}

// New creates a new middleware chain.
func (m *MW) New() *MW {
	return &MW{handlers: make([]http.Handler, 0, 0)}
}

// Add middleware to handlers slice of Handler
func (m *MW) Add(h ...http.Handler) {
	m.handlers = append(m.handlers, h...)
}

// So we can satisfy the http.Handler interface.
// func (m *MW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	for _, h := range m.handlers {
// 		h.ServeHTTP(w, r)
// 	}
//
// 	m.mux.ServeHTTP(w, r)
// }
