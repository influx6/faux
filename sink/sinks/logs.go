package sinks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/influx6/faux/sink"
)

// contains the set of filter values used by the Stdout sinks
const (
	INFO sink.Filter = iota + 1
	DEBUG
	ERROR
	NOTICE
)

// contains different color types for printing.
var (
	blue  = color.New(color.FgBlue)
	cyan  = color.New(color.FgCyan)
	red   = color.New(color.FgRed)
	white = color.New(color.FgWhite)
	black = color.New(color.FgBlack)
)

//==============================================================================

// Info returns a sink.Entry based on the provided message.
func Info(message string, m ...interface{}) sink.Entry {
	return sink.Entry{
		Message: fmt.Sprintf(message, m...),
		ID:      INFO,
		Pair:    (*sink.Pair)(nil),
	}
}

// Error returns a sink.Entry based on the provided message.
func Error(mi interface{}, m ...interface{}) sink.Entry {
	var message string

	switch mo := mi.(type) {
	case string:
		message = fmt.Sprintf(mo, m...)
		break
	case error:
		message = mo.Error()
		break
	}

	return sink.Entry{
		Message: message,
		ID:      ERROR,
		Pair:    (*sink.Pair)(nil),
	}
}

// Notice returns a sink.Entry based on the provided message.
func Notice(message string, m ...interface{}) sink.Entry {
	return sink.Entry{
		Message: fmt.Sprintf(message, m...),
		ID:      NOTICE,
		Pair:    (*sink.Pair)(nil),
	}
}

// Debug returns a sink.Entry based on the provided message.
func Debug(message string, m ...interface{}) sink.Entry {
	return sink.Entry{
		Message: fmt.Sprintf(message, m...),
		ID:      DEBUG,
		Pair:    (*sink.Pair)(nil),
	}
}

//==============================================================================

// Stdout emits all entries into the systems stdout.
type Stdout struct{}

// Emit implements the sink.Sink interface and does nothing with the
// provided entry.
func (Stdout) Emit(e sink.Entry) error {
	var bu bytes.Buffer

	switch e.ID {
	case INFO:
		blue.Fprint(&bu, "INFO")
		break
	case DEBUG:
		cyan.Fprint(&bu, "DEBUG")
		break
	case ERROR:
		red.Fprint(&bu, "ERROR")
		break
	case NOTICE:
		white.Fprint(&bu, "NOTICE")
		break
	}

	black.Fprint(&bu, "[opening]")
	bu.Write([]byte(":"))

	if e.Message != "" {
		bu.Write([]byte("\t\t"))
		bu.Write([]byte(e.Message))
	}

	printEntryParams(&bu, e)
	bu.Write([]byte("\n"))

	bu.WriteTo(os.Stdout)

	return nil
}

//==============================================================================

// Stderr emits all entries into the systems stderr.
type Stderr struct{}

// Emit implements the sink.Sink interface and does nothing with the
// provided entry.
func (Stderr) Emit(e sink.Entry) error {
	var bu bytes.Buffer

	switch e.ID {
	case ERROR:
		red.Fprint(&bu, "ERROR")
		break
	default:
		return errors.New("Only Error ID allowed")
	}

	black.Fprint(&bu, "[opening]")
	bu.Write([]byte(":"))

	if e.Message != "" {
		bu.Write([]byte("\t\t"))
		bu.Write([]byte(e.Message))
	}

	printEntryParams(&bu, e)
	bu.Write([]byte("\n"))

	bu.WriteTo(os.Stdout)
	return nil
}

func printEntryParams(bu io.Writer, e sink.Entry) {
	bu.Write([]byte("\t\t"))

	fields := e.Fields()

	for key, val := range fields {
		switch e.ID {
		case INFO:
			blue.Fprint(bu, key)
			black.Fprint(bu, "=")
			black.Fprint(bu, printValue(val))
			bu.Write([]byte(" "))
			break
		case DEBUG:
			cyan.Fprint(bu, key)
			black.Fprint(bu, "=")
			black.Fprint(bu, printValue(val))
			bu.Write([]byte(" "))
			break
		case ERROR:
			red.Fprint(bu, key)
			black.Fprint(bu, "=")
			black.Fprint(bu, printValue(val))
			bu.Write([]byte(" "))
			break
		case NOTICE:
			white.Fprint(bu, key)
			black.Fprint(bu, "=")
			black.Fprint(bu, printValue(val))
			bu.Write([]byte(" "))
			break
		}
	}

}

type stringer interface {
	String() string
}

func printValue(val interface{}) string {
	switch mo := val.(type) {
	case stringer:
		return mo.String()
	default:
		data, err := json.Marshal(val)
		if err != nil {
			return "-"
		}

		return string(data)
	}
}
