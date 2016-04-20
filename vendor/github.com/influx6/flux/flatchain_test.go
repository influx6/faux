package flux

import "testing"

func TestFlatChain(t *testing.T) {
	fo := NewFlatChain(func(err error, d interface{}, next NextHandler) {
		if d != 40 {
			FatalFailed(t, "Received incorrect value %d", d)
		}
		LogPassed(t, "Received correct value %d", d)
		next(err, d)
	})

	fo.Call(nil, 40)
}
