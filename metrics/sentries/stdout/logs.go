package stdout

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/influx6/faux/metrics"
)

// contains different color types for printing.
var (
	blue  = color.New(color.FgBlue)
	cyan  = color.New(color.FgCyan)
	red   = color.New(color.FgRed)
	white = color.New(color.FgWhite)
	black = color.New(color.FgBlack)
)

// sets of const used in package.
const (
	logTypeKey = "LogKEY"
	INFO       = "INFO"
	DEBUG      = "DEBUG"
	ERROR      = "ERROR"
	NOTICE     = "NOTICE"
	UNKOWN     = "Unknown"
)

//==============================================================================

// Info returns a metrics.Entry based on the provided message.
func Info(message string, m ...interface{}) metrics.Entry {
	return metrics.Entry{
		Message: fmt.Sprintf(message, m...),
		Pair:    (new(metrics.Pair)).Append(logTypeKey, INFO),
	}
}

// Error returns a metrics.Entry based on the provided message.
func Error(mi interface{}, m ...interface{}) metrics.Entry {
	var message string

	switch mo := mi.(type) {
	case string:
		message = fmt.Sprintf(mo, m...)
		break
	case error:
		message = mo.Error()
		break
	}

	return metrics.Entry{
		Message: message,
		Pair:    (new(metrics.Pair)).Append(logTypeKey, ERROR),
	}
}

// Notice returns a metrics.Entry based on the provided message.
func Notice(message string, m ...interface{}) metrics.Entry {
	return metrics.Entry{
		Message: fmt.Sprintf(message, m...),
		Pair:    (new(metrics.Pair)).Append(logTypeKey, NOTICE),
	}
}

// Debug returns a metrics.Entry based on the provided message.
func Debug(message string, m ...interface{}) metrics.Entry {
	return metrics.Entry{
		Message: fmt.Sprintf(message, m...),
		Pair:    (new(metrics.Pair)).Append(logTypeKey, DEBUG),
	}
}

//==============================================================================

// Stdout emits all entries into the systems stdout.
type Stdout struct{}

// Emit implements the metrics.metrics interface and does nothing with the
// provided entry.
func (Stdout) Emit(e metrics.Entry) error {
	var bu bytes.Buffer

	var id string

	if cid, ok := e.Get(logTypeKey); ok {
		if sid, ok := cid.(string); ok {
			id = sid
		}
	}

	switch id {
	case INFO:
		blue.Fprint(&bu, INFO)
		break
	case DEBUG:
		cyan.Fprint(&bu, DEBUG)
		break
	case ERROR:
		red.Fprint(&bu, ERROR)
		break
	case NOTICE:
		white.Fprint(&bu, NOTICE)
		break
	default:
		white.Fprint(&bu, UNKOWN)
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

// Emit implements the metrics.metrics interface and does nothing with the
// provided entry.
func (Stderr) Emit(e metrics.Entry) error {
	var bu bytes.Buffer

	var id string

	if cid, ok := e.Get(logTypeKey); ok {
		if sid, ok := cid.(string); ok {
			id = sid
		}
	}

	switch id {
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

func printEntryParams(bu io.Writer, e metrics.Entry) {
	bu.Write([]byte("\t\t"))

	fields := e.Fields()

	var id string

	if cid, ok := e.Get(logTypeKey); ok {
		if sid, ok := cid.(string); ok {
			id = sid
		}
	}

	for key, val := range fields {

		// We don't want keyless or value-less items.
		if key == "" || val == nil {
			continue
		}

		switch id {
		case INFO:
			blue.Fprint(bu, key)
			blue.Fprint(bu, "=")
			white.Fprint(bu, printValue(val))
			bu.Write([]byte(" "))
			break
		case DEBUG:
			cyan.Fprint(bu, key)
			blue.Fprint(bu, "=")
			white.Fprint(bu, printValue(val))
			bu.Write([]byte(" "))
			break
		case ERROR:
			red.Fprint(bu, key)
			blue.Fprint(bu, "=")
			white.Fprint(bu, printValue(val))
			bu.Write([]byte(" "))
			break
		case NOTICE:
			white.Fprint(bu, key)
			blue.Fprint(bu, "=")
			white.Fprint(bu, printValue(val))
			bu.Write([]byte(" "))
			break
		}
	}

}

type stringer interface {
	String() string
}

func printValue(val interface{}) string {
	switch bo := val.(type) {
	case string:
		return bo
	case stringer:
		return bo.String()
	case int:
		return strconv.Itoa(bo)
	case int64:
		return strconv.Itoa(int(bo))
	case rune:
		return strconv.QuoteRune(bo)
	case bool:
		return strconv.FormatBool(bo)
	case byte:
		return strconv.QuoteRune(rune(bo))
	case float64:
		return strconv.FormatFloat(bo, 'f', 4, 64)
	case float32:
		return strconv.FormatFloat(float64(bo), 'f', 4, 64)
	default:
		data, err := json.Marshal(val)
		if err != nil {
			return "-"
		}

		return string(data)
	}
}
