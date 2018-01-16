package ober

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type (
	// Ober is the main type which will be the fw basically
	Ober struct {
		Middleware *MW
		Server     *http.Server
		Logger     *logrus.Logger
		Listener   net.Listener

		DisableHTTP2 bool
		CertFile     string
		KeyFile      string
		Address      string

		ctx        context.Context
		router     *mux.Router
		middleware []Middleware
	}

	// Middleware type for handling middlewares
	Middleware func(http.Handler) http.Handler
)

// Functional Options ===========================

// DisableHTTP2 option to disable HTTP2 protocol since using TLS
// will transparently set the http2 protocol on by default
func DisableHTTP2(disableHTTP2 bool) func(*Ober) {
	return func(o *Ober) {
		o.DisableHTTP2 = disableHTTP2
	}
}

// CertFile option adds the TLS certificate
func CertFile(cert string) func(*Ober) {
	return func(o *Ober) {
		o.CertFile = cert
	}
}

// KeyFile option adds the TLS key
func KeyFile(key string) func(*Ober) {
	return func(o *Ober) {
		o.KeyFile = key
	}
}

// Address adds the address to the server
func Address(addr string) func(*Ober) {
	return func(o *Ober) {
		o.Address = addr
	}
}

// Middleware ===================================

// Use ...
func (o *Ober) Use(middleware ...Middleware) {
	o.middleware = append(o.middleware, middleware...)
}

// Router? ======================================

// Add a new route...will change this
func (o *Ober) Add(path string, handler http.HandlerFunc) {
	o.router.HandleFunc(path, handler)
}

// Server Methods ===============================

// New creates an instance of Ober.
func New(options ...func(*Ober)) (o *Ober) {
	o = &Ober{
		Middleware:   new(MW),
		Server:       new(http.Server),
		Logger:       logrus.New(),
		DisableHTTP2: false,
		Address:      ":9999",
		ctx:          context.Background(),
	}

	for _, option := range options {
		option(o)
	}

	o.Logger.Formatter = &logrus.JSONFormatter{}
	o.Server.Handler = o
	o.Server.Addr = o.Address
	o.router = mux.NewRouter()

	return o
}

// Router returns router.
func (o *Ober) Router() *mux.Router {
	return o.router
}

// Start starts an HTTPs server.
func (o *Ober) Start() (err error) {

	if o.CertFile == "" || o.KeyFile == "" {
		return errors.New("invalid tls configuration")
	}

	// Our server
	s := o.Server

	// TLS configurations
	if s.TLSConfig == nil {
		s.TLSConfig = new(tls.Config)
		s.TLSConfig.Certificates = make([]tls.Certificate, 1)
		s.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(o.CertFile, o.KeyFile)
		if err != nil {
			return errors.New("tls certificates error")
		}
	}

	if o.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "http/1.1")
	}

	o.Logger.Println("⇨ https server started on", s.Addr)

	return s.ListenAndServeTLS(o.CertFile, o.KeyFile)
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (o *Ober) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// o.Logger.Printf("%v", o.ctx)
	// r = r.WithContext(o.ctx)
	// for _, h := range o.middleware {
	// 	h(r).ServeHTTP(w, r)
	// }

	o.router.ServeHTTP(w, r)
}
