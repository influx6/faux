package httputil

import (
	"net/tls"

  "golang.org/x/crypto/acme/autocert"
)

// LetsEncryptTLS returns a tls.Config instance which retrieves its
	manager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}

	var tlsConfig tls.Config
	tlsConfig.Addr = address
	s.TLSConfig.GetCertificate = manager.GetCertificate

	if http2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	}
	return manager, tlsConfig
	return manager, tlsConfig
    return manager, tlsConfig
}
