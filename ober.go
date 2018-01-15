package ober

import (
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
		Middleware *Middleware
		Server     *http.Server
		Logger     *logrus.Logger
		Listener   net.Listener

		DisableHTTP2 bool

		router *mux.Router
	}
)

// Functional Options ===========================
// Options to disable HTTP2 protocol since using TLS will
// transparently set the http2 protocol on by default
func DisableHTTP2(disableHTTP2 bool) func(*Ober) {
	return func(o *Ober) {
		o.DisableHTTP2 = disableHTTP2
	}
}

// Router? ======================================
// Adds a new route...will change this
func (o *Ober) Add(path string, handler http.HandlerFunc) {
	o.router.HandleFunc(path, handler)
}

// Server Methods ===============================
// New creates an instance of Ober.
func New(options ...func(*Ober)) (o *Ober) {
	o = &Ober{
		Middleware:   new(Middleware),
		Server:       new(http.Server),
		Logger:       logrus.New(),
		DisableHTTP2: false,
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
	// o.Logger.Println("hello")
	for _, h := range o.Middleware.handlers {
		h.ServeHTTP(w, r)
	}

	o.router.ServeHTTP(w, r)
}
