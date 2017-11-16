// +build go1.8,go1.9

package httputil

import (
	"context"
)

// Close closes the underline server.
// It will gracefully close and shutdown the server.
func (s serverItem) Close(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
