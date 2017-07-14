package httputil

import (
	"crypto/tls"

	"golang.org/x/crypto/acme/autocert"
)

// LetsEncryptTLS returns a tls.Config instance which retrieves its
// its tls certificate from LetsEncrypt service.
func LetsEncryptTLS(address string, http2 bool) (*autocert.Manager, *tls.Config) {
	manager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}

	var tlsConfig tls.Config
	tlsConfig.Addr = address
	s.TLSConfig.GetCertificate = manager.GetCertificate

	if http2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}

	return manager, &tlsConfig
}
