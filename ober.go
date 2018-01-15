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
	Ober struct {
		Middleware *Middleware
		Server     *http.Server
		Logger     *logrus.Logger
		Listener   net.Listener

		DisableHTTP2 bool
		Debug        bool
		HideBanner   bool
		HidePort     bool

		router *mux.Router
		routes map[string]*Route
	}

	HandlerFunc func() error

	Route struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Name   string `json:"name"`
	}
)

// New creates an instance of Ober.
func New() (o *Ober) {
	o = &Ober{
		Middleware:   new(Middleware),
		Server:       new(http.Server),
		Logger:       logrus.New(),
		DisableHTTP2: false,
	}

	o.Logger.Formatter = &logrus.JSONFormatter{}
	o.Server.Handler = o
	o.router = mux.NewRouter()

	return
}

// Adds a new route...will change this
func (o *Ober) Add(path string, handler http.HandlerFunc) {
	o.router.HandleFunc(path, handler)
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
