// +build go1.2,go1.3,go1.4,go1.5,go1.6,go1.7,!go1.8,!go1.9

package httputil

import (
	"context"
)

// Close closes the underline server.
// It will forcefully close the server.
func (s serverItem) Close(ctx context.Context) error {
	return s.listener.Close()
}
