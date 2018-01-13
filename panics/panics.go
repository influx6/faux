package panics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Guard provides a function for handling panics safely, returning the
// error recieved after running a provided function.
func Guard(fx func() error) error {
	var err error

	func() {
		defer func() {
			if ex := recover(); ex != nil {
				if exx, ok := ex.(error); ok {
					err = exx
				} else {
					err = fmt.Errorf("%v", ex)
				}
			}
		}()
		err = fx()
	}()

	return err
}

// GuardWith provides a function guard for a function in need of an in input
func GuardWith(fx func(interface{}) error, d interface{}) error {
	var err error

	func() {
		defer func() {
			if eux := recover(); eux != nil {
				if ex, ok := eux.(error); ok {
					err = ex
				} else {
					err = fmt.Errorf("%s", ToString(eux, false))
				}
			}
		}()
		err = fx(d)
	}()

	return err
}

// GuardDefer returns a function which guarded by a defer..recover() for handling
// panic recover.
func GuardDefer(fx func(interface{}) error) func(interface{}) error {
	return func(d interface{}) error {
		var err error

		func() {
			defer func() {
				if eux := recover(); eux != nil {
					if ex, ok := eux.(error); ok {
						err = ex
					} else {
						err = fmt.Errorf("%s", ToString(eux, false))
					}
				}
			}()
			err = fx(d)
		}()

		return err
	}
}

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

// Defer provides a recovery handler functions for use to automate
// the recovery processes and logs out any panic that occurs.
func Defer(op func(), logfn func(*bytes.Buffer)) {
	defer func() {
		if err := recover(); err != nil {
			if logfn != nil {
				var data bytes.Buffer
				trace := make([]byte, 10000)
				runtime.Stack(trace, true)
				data.Write([]byte("----------------------------------------------------------------"))
				data.Write([]byte("\n"))
				data.Write([]byte(fmt.Sprintf("Error: %+v\n", err)))
				data.Write([]byte("\n"))
				data.Write([]byte("----------------------------------------------------------------"))
				data.Write([]byte("\n"))
				data.Write(trace)
				data.Write([]byte("\n"))
				data.Write([]byte("----------------------------------------------------------------"))
				data.Write([]byte("\n"))
				logfn(&data)
			}
		}
	}()

	op()
}

// DeferReport provides a recovery handler functions for use to automate
// the recovery processes and logs out any panic that occurs.
func DeferReport(op func(), logfn func(*bytes.Buffer)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if logfn != nil {
					var data bytes.Buffer
					trace := make([]byte, 10000)
					runtime.Stack(trace, true)
					data.Write([]byte("----------------------------------------------------------------"))
					data.Write([]byte("\n"))
					data.Write([]byte(fmt.Sprintf("Error: %+v\n", err)))
					data.Write([]byte("\n"))
					data.Write([]byte("----------------------------------------------------------------"))
					data.Write([]byte("\n"))
					data.Write(trace)
					data.Write([]byte("\n"))
					data.Write([]byte("----------------------------------------------------------------"))
					data.Write([]byte("\n"))
					logfn(&data)
				}
			}
		}()

		op()
	}()
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

// ToString provides a string version of the value using json.Marshal.
func ToString(value interface{}, indent bool) string {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(value, "", "\n")
	} else {
		data, err = json.Marshal(value)
	}

	if err != nil {
		return ""
	}

	return string(data)
}
