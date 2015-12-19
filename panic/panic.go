package panic

import (
	"fmt"
	"runtime"
	"strings"
)

// RecoverHandler provides a recovery handler functions for use to automate the recovery processes
func RecoverHandler(tag string, opFunc func() error, recoFunc func(interface{})) error {
	defer func() {
		if err := recover(); err != nil {
			if recoFunc != nil {
				recoFunc(err)
			}
		}
	}()

	if err := opFunc(); err != nil {
		return err
	}

	return nil
}

// LogRecoverHandler provides a recovery handler functions for use to automate
// the recovery processes and logs out any panic that occurs.
func LogRecoverHandler(tag string, opFunc func() error) error {
	return RecoverHandler(tag, opFunc, func(err interface{}) {
		trace := make([]byte, 10000)
		count := runtime.Stack(trace, true)
		fmt.Printf("---------%s-Panic----------------:", strings.ToUpper(tag))
		fmt.Printf("Error: %+s", err)
		fmt.Printf("Stack of %d bytes: \n%+s\n", count, trace)
		fmt.Printf("---------%s--END-----------------:", strings.ToUpper(tag))
	})
}

// GoDefer lets you run a function inside a goroutine that gets a defer recovery,
// and reports to stdout if a panic occurs.
func GoDefer(title string, fx func()) {
	go LogRecoverHandler(title, func() error {
		fx()
		return nil
	})
}

// GoDeferQuietly lets you run a function inside a goroutine that gets a defer
// recovery but silently ignores the panic without logging to stdout.
func GoDeferQuietly(title string, fx func()) {
	go RecoverHandler(title, func() error {
		fx()
		return nil
	}, nil)
}
