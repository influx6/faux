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
	"github.com/influx6/faux/reflection"
)

// BlockDislay writes giving Entries as seperated blocks of contents where the each content is
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
func BlockDisplay(w io.Writer) metrics.Metrics {
	return BlockDisplayWith(w, "Message:", nil)
}

// BlockDislay writes giving Entries as seperated blocks of contents where the each content is
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
func BlockDisplayWith(w io.Writer, header string, filterFn func(metrics.Entry) bool) metrics.Metrics {
	return NewCustomEmitter(w, func(en metrics.Entry) []byte {
		if filterFn != nil && !filterFn(en) {
			return nil
		}

		var bu bytes.Buffer
		if header != "" {
			fmt.Fprintf(&bu, "%s %+s\n", header, en.Message)
		} else {
			fmt.Fprintf(&bu, "%+s\n", en.Message)
		}

		if en.Function != "" {
			fmt.Fprintf(&bu, "Function: %+s\n", en.Function)
		}

		print(en.Field, func(key []string, value string) {
			keyVal := strings.Join(key, ".")
			keyLength := len(keyVal) + 2
			valLength := len(value) + 2

			keyLines := printBlockLine(keyLength)
			valLines := printBlockLine(valLength)
			spaceLines := printSpaceLine(1)

			fmt.Fprintf(&bu, "+%s+%s+\n", keyLines, valLines)
			fmt.Fprintf(&bu, "|%s%s%s|%s%s%s|\n", spaceLines, keyVal, spaceLines, spaceLines, value, spaceLines)
			fmt.Fprintf(&bu, "+%s+%s+", keyLines, valLines)
			fmt.Fprintf(&bu, "\n")

		})

		bu.WriteString("\n")
		return bu.Bytes()
	})
}

// StackDislay writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  Message: We must create new standard behaviour
//	Function: BuildPack
//  - displayrange.address.bolder: "No 20 tokura flag"
//  - displayrange.bolder.size:  20
//
func StackDisplay(w io.Writer) metrics.Metrics {
	return StackDisplayWith(w, "Message:", "-", nil)
}

// StackDislayWith writes giving Entries as seperated blocks of contents where the each content is
// converted within a block like below:
//
//  [Header]: We must create new standard behaviour
//	Function: BuildPack
//  [tag] displayrange.address.bolder: "No 20 tokura flag"
//  [tag] displayrange.bolder.size:  20
//
func StackDisplayWith(w io.Writer, header string, tag string, filterFn func(metrics.Entry) bool) metrics.Metrics {
	return NewCustomEmitter(w, func(en metrics.Entry) []byte {
		if filterFn != nil && !filterFn(en) {
			return nil
		}

		var bu bytes.Buffer
		if header != "" {
			fmt.Fprintf(&bu, "%s %+s\n", header, en.Message)
		} else {
			fmt.Fprintf(&bu, "%+s\n", en.Message)
		}

		if tag == "" {
			tag = "-"
		}

		if en.Function != "" {
			fmt.Fprintf(&bu, "Function: %+s\n", en.Function)
		}

		print(en.Field, func(key []string, value string) {
			fmt.Fprintf(&bu, "%s %s: %+s\n", tag, strings.Join(key, "."), value)
		})

		bu.WriteString("\n")
		return bu.Bytes()
	})
}

//=====================================================================================

// SwitchEmitter returns a emitter that converts the behaviour of the output based on giving key and value from
// each Entry.
func SwitchEmitter(keyName string, w io.Writer, transformers map[string]func(metrics.Entry) []byte) metrics.Metrics {
	emitters := make(map[string]metrics.Metrics)

	for id, tm := range transformers {
		emitters[id] = NewCustomEmitter(w, tm)
	}

	return metrics.Switch(keyName, emitters)
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

func print(item interface{}, do func(key []string, val string)) {
	printInDepth(item, do, 0)
}

var maxDepth = 100

func printInDepth(item interface{}, do func(key []string, val string), depth int) {
	if depth >= maxDepth {
		return
	}

	if item == nil {
		return
	}

	itemType := reflect.TypeOf(item)

	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	switch itemType.Kind() {
	case reflect.Array, reflect.Slice:
		printArrays(item, do, depth+1)
	case reflect.Struct:
		if val, err := reflection.ToMap("json", item, true); err == nil {
			printMap(val, do, depth+1)
		}

	case reflect.Map:
		printMap(item, do, depth+1)
	default:
		do([]string{}, printValue(item))
	}
}

func printMap(items interface{}, do func(key []string, val string), depth int) {
	switch bo := items.(type) {
	case map[string]byte:
		for index, item := range bo {
			do([]string{index}, printValue(int(item)))
		}
	case map[string]float32:
		for index, item := range bo {
			do([]string{index}, printValue(item))
		}
	case map[string]float64:
		for index, item := range bo {
			do([]string{index}, printValue(item))
		}
	case map[string]int64:
		for index, item := range bo {
			do([]string{index}, printValue(item))
		}
	case map[string]int32:
		for index, item := range bo {
			do([]string{index}, printValue(item))
		}
	case map[string]int16:
		for index, item := range bo {
			do([]string{index}, printValue(item))
		}
	case map[string]time.Time:
		for index, item := range bo {
			do([]string{index}, printValue(item))
		}
	case map[string]int:
		for index, item := range bo {
			do([]string{index}, printValue(item))
		}
	case map[string][]interface{}:
		for index, item := range bo {
			printInDepth(item, func(key []string, value string) {
				if index == "" {
					do(key, value)
					return
				}

				do(append([]string{index}, key...), value)
			}, depth+1)
		}
	case map[string]interface{}:
		for index, item := range bo {
			printInDepth(item, func(key []string, value string) {
				if index == "" {
					do(key, value)
					return
				}

				do(append([]string{index}, key...), value)
			}, depth+1)
		}
	case map[string]string:
		for index, item := range bo {
			do([]string{index}, printValue(item))
		}
	case map[string][]byte:
		for index, item := range bo {
			do([]string{index}, printValue(string(item)))
		}
	case metrics.Field:
		printMap((map[string]interface{})(bo), do, depth+1)
	}
}

func printArrays(items interface{}, do func(index []string, val string), depth int) {
	switch bo := items.(type) {
	case []metrics.Field:
		for index, item := range bo {
			printMap((map[string]interface{})(item), func(key []string, val string) {
				do(append([]string{printValue(index)}, key...), val)
			}, depth+1)
		}
	case []map[string][]byte:
		for index, item := range bo {
			printMap(item, func(key []string, val string) {
				do(append([]string{printValue(index)}, key...), val)
			}, depth+1)
		}
	case []map[string][]interface{}:
		for index, item := range bo {
			printMap(item, func(key []string, val string) {
				do(append([]string{printValue(index)}, key...), val)
			}, depth+1)
		}
	case []map[string]interface{}:
		for index, item := range bo {
			printMap(item, func(key []string, val string) {
				do(append([]string{printValue(index)}, key...), val)
			}, depth+1)
		}
	case []byte:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(int(item)))
		}
	case []bool:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []interface{}:
		for index, item := range bo {
			printInDepth(item, func(key []string, value string) {
				do(append([]string{printValue(index)}, key...), value)
			}, depth+1)
		}
	case []time.Time:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []string:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []int:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []int64:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []int32:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []int16:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []int8:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []float32:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	case []float64:
		for index, item := range bo {
			do([]string{printValue(index)}, printValue(item))
		}
	}
}

func printValue(item interface{}) string {
	switch bo := item.(type) {
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
	default:
		data, err := json.Marshal(bo)
		if err != nil {
			return "-"
		}
		return string(data)
	}

	return ""
}
