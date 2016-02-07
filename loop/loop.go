package loop

// Mux defines the callback type signature for each runner.
type Mux func(float64)

// Looper defines a loop subscriber for a specific loop subscriber.
type Looper interface {
	End()
}
