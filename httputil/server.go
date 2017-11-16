package httputil

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/influx6/faux/netutils"
	"golang.org/x/crypto/acme/autocert"
)

// Server defines a type which closes a underline server and
// returns any error associated with the call.
type Server interface {
	Wait(...func())
	Close(context.Context) error
	TLSManager() *autocert.Manager
}

type serverItem struct {
	server   *http.Server
	listener net.Listener
	man      *autocert.Manager
}

// WaitAndShutdown attempts to wait till a interrupt is received.
func (s *serverItem) Wait(after ...func()) {
	defer s.Close(context.Background())
	WaitOnInterrupt(after...)
}

// TLSManager returns the autocert.Manager associated with the giving server
// for its tls certificates.
func (s *serverItem) TLSManager() *autocert.Manager {
	return s.man
}

// ListenWith will start a server and returns a ServerCloser which will allow
// closing of the server.
func ListenWith(tlsconfig *tls.Config, addr string, handler http.Handler) (Server, error) {
	listener, err := netutils.MakeListener("tcp", addr, tlsconfig)
	if err != nil {
		return nil, err
	}

	server, _, err := netutils.NewHTTPServer(listener, handler, tlsconfig)
	if err != nil {
		return nil, err
	}

	return &serverItem{
		server:   server,
		listener: listener,
	}, nil
}

// Listen will start a server and returns a ServerCloser which will allow
// closing of the server.
func Listen(tlsOK bool, addr string, handler http.Handler) (Server, error) {
	var tlsconfig *tls.Config
	var man *autocert.Manager

	if tlsOK {
		man, tlsconfig = LetsEncryptTLS(true)
	}

	listener, err := netutils.MakeListener("tcp", addr, tlsconfig)
	if err != nil {
		return nil, err
	}

	server, _, err := netutils.NewHTTPServer(listener, handler, tlsconfig)
	if err != nil {
		return nil, err
	}

	return &serverItem{
		server:   server,
		listener: listener,
		man:      man,
	}, nil
}
