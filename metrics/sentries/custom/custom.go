package custom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/influx6/faux/metrics"
)

// BlockDislay writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  Message: We must create new standard behaviour
//  +-----------------------------+------------------------------+
//  | displayrange.address.bolder | "No 20 tokura flag"          |
//  +-----------------------------+------------------------------+
//  +--------------------------+----------+
//  | displayrange.bolder.size |  20      |
//  +--------------------------+----------+
//
func BlockDisplay(w io.Writer, header string) metrics.Metrics {
	return BlockDisplayWith(w, "Message", true, nil)
}

// BlockDislay writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  Message: We must create new standard behaviour
//  +-----------------------------+------------------------------+
//  | displayrange.address.bolder | "No 20 tokura flag"          |
//  +-----------------------------+------------------------------+
//  +--------------------------+----------+
//  | displayrange.bolder.size |  20      |
//  +--------------------------+----------+
//
func BlockDisplayWith(w io.Writer, header string, stacked bool, filterFn func(metrics.Entry) bool) metrics.Metrics {
	return NewCustomEmitter(w, func(en metrics.Entry) []byte {
		if filterFn != nil && !filterFn(en) {
			return nil
		}

		var ok bool
		var message string

		if en.Message == "" {
			message, ok = en.GetString("message")
			if !ok {
				message = metrics.DefaultMessage
			}
		} else {
			message = en.Message
		}

		var bu bytes.Buffer
		if header != "" {
			fmt.Fprintf(&bu, "%s: %+s\n", header, message)
		} else {
			fmt.Fprintf(&bu, "%+s\n", message)
		}

		print(en.Fields(), func(key string, value string) {
			keyLength := len(key) + 2
			valLength := len(value) + 2

			keyLines := printBlockLine(keyLength)
			valLines := printBlockLine(valLength)
			spaceLines := printSpaceLine(1)

			fmt.Fprintf(&bu, "+%s+%s+\n", keyLines, valLines)
			fmt.Fprintf(&bu, "|%s%s%s|%s%s%s|\n", spaceLines, key, spaceLines, spaceLines, value, spaceLines)
			fmt.Fprintf(&bu, "+%s+%s+", keyLines, valLines)

			if stacked {
				fmt.Fprintf(&bu, "\n")
				return
			}

			fmt.Fprintf(&bu, " ")
		})

		bu.WriteString("\n")
		return bu.Bytes()
	})
}

// StackDislay writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  Message: We must create new standard behaviour
//  - displayrange.address.bolder: "No 20 tokura flag"
//  - displayrange.bolder.size:  20
//
func StackDisplay(w io.Writer) metrics.Metrics {
	return StackDisplayWith(w, "Message", "-", nil)
}

// StackDislayWith writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  [Header]: We must create new standard behaviour
//  [tag] displayrange.address.bolder: "No 20 tokura flag"
//  [tag] displayrange.bolder.size:  20
//
func StackDisplayWith(w io.Writer, header string, tag string, filterFn func(metrics.Entry) bool) metrics.Metrics {
	return NewCustomEmitter(w, func(en metrics.Entry) []byte {
		if filterFn != nil && !filterFn(en) {
			return nil
		}
		var ok bool
		var message string

		if en.Message == "" {
			message, ok = en.GetString("message")
			if !ok {
				message = metrics.DefaultMessage
			}
		} else {
			message = en.Message
		}

		var bu bytes.Buffer
		if header != "" {
			fmt.Fprintf(&bu, "%s: %+s\n", header, message)
		} else {
			fmt.Fprintf(&bu, "%+s\n", message)
		}

		if tag == "" {
			tag = "-"
		}

		print(en.Fields(), func(key string, value string) {
			fmt.Fprintf(&bu, "%s %s: %+q\n", tag, key, value)
		})

		bu.WriteString("\n")
		return bu.Bytes()
	})
}

//=====================================================================================

// CustomEmitter emits all entries into the entries into a sink io.writer after
// transformation from giving transformer function..
type CustomEmitter struct {
	Sink      io.Writer
	Transform func(metrics.Entry) []byte
}

// NewCustomEmitter returns a new instance of CustomEmitter.
func NewCustomEmitter(w io.Writer, transform func(metrics.Entry) []byte) *CustomEmitter {
	return &CustomEmitter{
		Sink:      w,
		Transform: transform,
	}
}

// Emit implements the metrics.metrics interface.
func (ce *CustomEmitter) Emit(e metrics.Entry) error {
	_, err := ce.Sink.Write(ce.Transform(e))
	return err
}

//=====================================================================================

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

func print(item interface{}, do func(key string, val string)) {
	itemType := reflect.TypeOf(item)

	switch itemType.Kind() {
	case reflect.Array:
		printArrays(item, do)
	case reflect.Map:
		printMap(item, do)
	default:
		do("-", printValue(item))
	}
}

func printMap(items interface{}, do func(key string, val string)) {
	switch bo := items.(type) {
	case map[string]byte:
		for index, item := range bo {
			do(index, printValue(int(item)))
		}
	case map[string]float32:
		for index, item := range bo {
			do(index, printValue(item))
		}
	case map[string]float64:
		for index, item := range bo {
			do(index, printValue(item))
		}
	case map[string]int64:
		for index, item := range bo {
			do(index, printValue(item))
		}
	case map[string]int32:
		for index, item := range bo {
			do(index, printValue(item))
		}
	case map[string]int16:
		for index, item := range bo {
			do(index, printValue(item))
		}
	case map[string]time.Time:
		for index, item := range bo {
			do(index, printValue(item))
		}
	case map[string]int:
		for index, item := range bo {
			do(index, printValue(item))
		}
	case map[string]string:
		for index, item := range bo {
			do(index, printValue(item))
		}
	case map[string][]byte:
		for index, item := range bo {
			do(index, printValue(string(item)))
		}
	case metrics.Fields:
		for index, item := range bo {
			switch vItem := item.(type) {
			case map[string][]byte:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]time.Time:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]float32:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]float64:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int8:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int16:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int64:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int32:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]string:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			default:
				do(index, printValue(item))
			}
		}
	case map[string]interface{}:
		for index, item := range bo {
			switch vItem := item.(type) {
			case map[string][]byte:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]time.Time:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]float32:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]float64:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int8:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int16:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int64:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int32:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]int:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			case map[string]string:
				printMap(vItem, func(key string, val string) {
					do(fmt.Sprintf("%s.%s", index, key), val)
				})
			default:
				do(index, printValue(item))
			}
		}
	}
}

func printArrays(items interface{}, do func(index string, val string)) {
	switch bo := items.(type) {
	case []map[string][]byte:
		for index, item := range bo {
			printMap(item, func(key string, val string) {
				do(fmt.Sprintf("%d[%s]", index, key), val)
			})
		}
	case []map[string]interface{}:
		for index, item := range bo {
			printMap(item, func(key string, val string) {
				do(fmt.Sprintf("%d[%s]", index, key), val)
			})
		}
	case []byte:
		for index, item := range bo {
			do(printValue(index), printValue(int(item)))
		}
	case []bool:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []interface{}:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []time.Time:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []string:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []int:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []int64:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []int32:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []int16:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []int8:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []float32:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	case []float64:
		for index, item := range bo {
			do(printValue(index), printValue(item))
		}
	}
}

func printValue(item interface{}) string {
	switch bo := item.(type) {
	case string:
		return bo
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
	default:
		data, err := json.Marshal(bo)
		if err != nil {
			return "-"
		}
		return string(data)
	}

	return ""
}
