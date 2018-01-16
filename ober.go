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
	// Main type for package
	Ober struct {
		Middleware *MW
		Server     *http.Server
		Logger     *logrus.Logger
		Listener   net.Listener

		DisableHTTP2 bool

		ctx        context.Context
		router     *mux.Router
		middleware []Middleware
	}

	Middleware func(http.Handler) http.Handler
)

// Functional Options ===========================
// Options to disable HTTP2 protocol since using TLS will
// transparently set the http2 protocol on by default
func DisableHTTP2(disableHTTP2 bool) func(*Ober) {
	return func(o *Ober) {
		o.DisableHTTP2 = disableHTTP2
	}
}

// Middleware ===================================
func (o *Ober) Use(middleware ...Middleware) {
	o.middleware = append(o.middleware, middleware...)
}

// Router? ======================================
// Adds a new route...will change this
func (o *Ober) Add(path string, handler http.Handler) {
	o.router.Handle(path, handler)
}

// Server Methods ===============================
// New creates an instance of Ober.
func New(options ...func(*Ober)) (o *Ober) {
	o = &Ober{
		Middleware:   new(MW),
		Server:       new(http.Server),
		Logger:       logrus.New(),
		DisableHTTP2: false,
		ctx:          context.Background(),
	}

	for _, option := range options {
		option(o)
	}

	o.Logger.Formatter = &logrus.JSONFormatter{}
	o.Server.Handler = o
	o.router = mux.NewRouter()

	return o
}

// Start starts an HTTPs server.
func (o *Ober) Start(address string, certFile, keyFile string) (err error) {

	if certFile == "" || keyFile == "" {
		return errors.New("invalid tls configuration")
	}

	s := o.Server
	s.Addr = address

	// TLS configurations
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.Certificates = make([]tls.Certificate, 1)
	s.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		o.Logger.Printf("%v", err)
		return errors.New("tls certificates error")
	}

	if o.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "http/1.1")
	}

	o.Logger.Println("⇨ https server started on", s.Addr)

	return s.ListenAndServeTLS(certFile, keyFile)
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (o *Ober) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// o.Logger.Printf("%v", o.ctx)
	r = r.WithContext(o.ctx)
	for _, h := range o.middleware {
		h(r).ServeHTTP(w, r)
	}

	o.router.ServeHTTP(w, r)
}
