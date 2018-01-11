package custom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/influx6/faux/metrics"
)

var (
	red     = color.New(color.FgRed)
	green   = color.New(color.FgGreen)
	white   = color.New(color.FgWhite)
	yellow  = color.New(color.FgHiYellow)
	magenta = color.New(color.FgHiMagenta)
)

// FlatDisplay writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  Message: We must create new standard behaviour 	Function: BuildPack  |  display: red,  words: 20,
//
//  Message: We must create new standard behaviour 	Function: BuildPack  |  display: red,  words: 20,
//
func FlatDisplay(w io.Writer) metrics.Processors {
	return FlatDisplayWith(w, "Message:", nil)
}

// FlatDisplayWith writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  [Header]: We must create new standard behaviour 	Function: BuildPack  |  display: red,  words: 20,
//
//  [Header]: We must create new standard behaviour 	Function: BuildPack  |  display: red,  words: 20,
//
func FlatDisplayWith(w io.Writer, header string, filterFn func(metrics.Entry) bool) metrics.Processors {
	return NewEmitter(w, func(en metrics.Entry) []byte {
		if filterFn != nil && !filterFn(en) {
			return nil
		}

		var bu bytes.Buffer
		bu.WriteString("\n")

		if header != "" {
			fmt.Fprintf(&bu, "%s %+s", green.Sprint(header), printAtLevel(en.Level, en.Message))
		} else {
			fmt.Fprintf(&bu, "%+s", printAtLevel(en.Level, en.Message))
		}

		if en.ID != "" {
			fmt.Fprintf(&bu, "ID: %+s\n", printAtLevel(en.Level, en.ID))
		}

		fmt.Fprint(&bu, printSpaceLine(2))

		if en.Function != "" {
			fmt.Fprintf(&bu, "%s: %+s\n", green.Sprint("Function"), en.Function)
			fmt.Fprint(&bu, printSpaceLine(2))
			fmt.Fprintf(&bu, "%s: %+s:%d", green.Sprint("File"), en.File, en.Line)
			fmt.Fprint(&bu, printSpaceLine(2))
		}

		fmt.Fprint(&bu, printSpaceLine(2))

		for key, value := range en.Field {
			fmt.Fprintf(&bu, "%+s: %+s", green.Sprint(key), printItem(value))
			fmt.Fprint(&bu, printSpaceLine(2))
		}

		bu.WriteString("\n")
		return bu.Bytes()
	})
}

//=====================================================================================

// BlockDisplay writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  Message: We must create new standard behaviour
//	Function: BuildPack
//  +-----------------------------+------------------------------+
//  | displayrange.address.bolder | "No 20 tokura flag"          |
//  +-----------------------------+------------------------------+
//  +--------------------------+----------+
//  | displayrange.bolder.size |  20      |
//  +--------------------------+----------+
//
func BlockDisplay(w io.Writer) metrics.Processors {
	return BlockDisplayWith(w, "Message:", nil)
}

// BlockDisplayWith writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  Message: We must create new standard behaviour
//	Function: BuildPack
//  +-----------------------------+------------------------------+
//  | displayrange.address.bolder | "No 20 tokura flag"          |
//  +-----------------------------+------------------------------+
//  +--------------------------+----------+
//  | displayrange.bolder.size |  20      |
//  +--------------------------+----------+
//
func BlockDisplayWith(w io.Writer, header string, filterFn func(metrics.Entry) bool) metrics.Processors {
	return NewEmitter(w, func(en metrics.Entry) []byte {
		if filterFn != nil && !filterFn(en) {
			return nil
		}

		var bu bytes.Buffer
		if header != "" {
			fmt.Fprintf(&bu, "%s %+s\n", green.Sprint(header), printAtLevel(en.Level, en.Message))
		} else {
			fmt.Fprintf(&bu, "%+s\n", printAtLevel(en.Level, en.Message))
		}

		if en.ID != "" {
			fmt.Fprintf(&bu, "ID: %+s\n", printAtLevel(en.Level, en.ID))
		}

		if en.Function != "" {
			fmt.Fprintf(&bu, "%s: %+s\n", green.Sprint("Function"), en.Function)
			fmt.Fprintf(&bu, "%s: %+s:%d\n", green.Sprint("File"), en.File, en.Line)
		}

		for key, val := range en.Field {
			value := printItem(val)
			keyLength := len(key) + 2
			valLength := len(value) + 2

			keyLines := printBlockLine(keyLength)
			valLines := printBlockLine(valLength)
			spaceLines := printSpaceLine(1)

			fmt.Fprintf(&bu, "+%s+%s+\n", keyLines, valLines)
			fmt.Fprintf(&bu, "|%s%s%s|%s%s%s|\n", spaceLines, green.Sprint(key), spaceLines, spaceLines, value, spaceLines)
			fmt.Fprintf(&bu, "+%s+%s+", keyLines, valLines)
			fmt.Fprintf(&bu, "\n")
		}

		bu.WriteString("\n")
		return bu.Bytes()
	})
}

//=====================================================================================

// StackDisplay writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  Message: We must create new standard behaviour
//	Function: BuildPack
//  - displayrange.address.bolder: "No 20 tokura flag"
//  - displayrange.bolder.size:  20
//
func StackDisplay(w io.Writer) metrics.Processors {
	return StackDisplayWith(w, "Message:", "-", nil)
}

// StackDisplayWith writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  [Header]: We must create new standard behaviour
//	Function: BuildPack
//  [tag] displayrange.address.bolder: "No 20 tokura flag"
//  [tag] displayrange.bolder.size:  20
//
func StackDisplayWith(w io.Writer, header string, tag string, filterFn func(metrics.Entry) bool) metrics.Processors {
	return NewEmitter(w, func(en metrics.Entry) []byte {
		if filterFn != nil && !filterFn(en) {
			return nil
		}

		var bu bytes.Buffer
		if header != "" {
			fmt.Fprintf(&bu, "%s %+s\n", green.Sprint(header), printAtLevel(en.Level, en.Message))
		} else {
			fmt.Fprintf(&bu, "%+s\n", printAtLevel(en.Level, en.Message))
		}

		if en.ID != "" {
			fmt.Fprintf(&bu, "ID: %+s\n", printAtLevel(en.Level, en.ID))
		}

		if tag == "" {
			tag = "-"
		}

		if en.Function != "" {
			fmt.Fprintf(&bu, "%s: %+s\n", green.Sprint("Function"), en.Function)
			fmt.Fprintf(&bu, "%s: %+s:%d\n", green.Sprint("File"), en.File, en.Line)
		}

		for key, value := range en.Field {
			fmt.Fprintf(&bu, "%s %s: %+s\n", tag, green.Sprintf(key), printItem(value))
		}

		bu.WriteString("\n")
		return bu.Bytes()
	})
}

//=====================================================================================

// Emitter emits all entries into the entries into a sink io.writer after
// transformation from giving transformer function..
type Emitter struct {
	Sink      io.Writer
	Transform func(metrics.Entry) []byte
}

// NewEmitter returns a new instance of Emitter.
func NewEmitter(w io.Writer, transform func(metrics.Entry) []byte) *Emitter {
	return &Emitter{
		Sink:      w,
		Transform: transform,
	}
}

// Handle implements the metrics.metrics interface.
func (ce *Emitter) Handle(e metrics.Entry) error {
	_, err := ce.Sink.Write(ce.Transform(e))
	return err
}

//=====================================================================================

func printAtLevel(lvl metrics.Level, message string) string {
	switch lvl {
	case metrics.ErrorLvl:
		return red.Sprint(message)
	case metrics.InfoLvl:
		return white.Sprint(message)
	case metrics.RedAlertLvl:
		return magenta.Sprint(message)
	case metrics.YellowAlertLvl:
		return yellow.Sprint(message)
	}

	return message
}

func printSpaceLine(length int) string {
	var lines []string

	for i := 0; i < length; i++ {
		lines = append(lines, " ")
	}

	return strings.Join(lines, "")
}

func printBlockLine(length int) string {
	var lines []string

	for i := 0; i < length; i++ {
		lines = append(lines, "-")
	}

	return strings.Join(lines, "")
}

type stringer interface {
	String() string
}

func printItem(item interface{}) string {
	switch bo := item.(type) {
	case stringer:
		return bo.String()
	case string:
		return `"` + bo + `"`
	case error:
		return bo.Error()
	case int:
		return strconv.Itoa(bo)
	case int8:
		return strconv.Itoa(int(bo))
	case int16:
		return strconv.Itoa(int(bo))
	case int64:
		return strconv.Itoa(int(bo))
	case time.Time:
		return bo.UTC().String()
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
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Sprintf("%#v", item)
	}

	return string(data)
}
