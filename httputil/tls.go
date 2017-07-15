package httputil

import (
	"crypto/tls"

	"golang.org/x/crypto/acme/autocert"
)

// LetsEncryptTLS returns a tls.Config instance which retrieves its
// its tls certificate from LetsEncrypt service.
func LetsEncryptTLS(http2 bool) (*autocert.Manager, *tls.Config) {
	manager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}

	var tlsConfig tls.Config
	tlsConfig.GetCertificate = manager.GetCertificate

	if http2 {
		tlsConfig.NextProtos = append(tlsConfig.NextProtos, "h2")
	}

	return manager, &tlsConfig
}
