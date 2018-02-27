package metrics

import (
	"fmt"
	"strings"
	"time"
)

// level constants
const (
	RedAlertLvl    Level = iota // Immediately notify everyone by mail level, because this is bad
	YellowAlertLvl              // Immediately notify everyone but we can wait to tomorrow
	ErrorLvl                    // Error occured with some code due to normal opperation or odd behaviour (not critical)
	InfoLvl                     // Information for view about code operation (replaces Debug, Notice, Trace).
)

// Entry represent a giving record of data at a giving period of time.
type Entry struct {
	ID        string      `json:"id"`
	Function  string      `json:"function"`
	File      string      `json:"file"`
	Type      string      `json:"type"`
	Line      int         `json:"line"`
	Level     Level       `json:"level"`
	Field     Field       `json:"fields"`
	Time      time.Time   `json:"time"`
	Message   string      `json:"message"`
	Filter    interface{} `json:"filter"`
	Tags      []string    `json:"tags"`
	Trace     Trace       `json:"trace"`
	Timelapse []Timelapse `json:"timelapse"`
}

// Level defines a int type which represent the a giving level of entry for a giving entry.
type Level int

// GetLevel returns Level value for the giving string.
// It returns -1 if it does not know the level string.
func GetLevel(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "redalert", "redalartlvl":
		return RedAlertLvl
	case "yellowalert", "yellowalertlvl":
		return YellowAlertLvl
	case "error", "errorlvl":
		return ErrorLvl
	case "info", "infolvl":
		return InfoLvl
	}

	return -1
}

// String returns the string version of the Level.
func (l Level) String() string {
	switch l {
	case RedAlertLvl:
		return "REDALERT"
	case YellowAlertLvl:
		return "YELLOWALERT"
	case ErrorLvl:
		return "ERROR"
	case InfoLvl:
		return "INFO"
	}

	return "UNKNOWN"
}

// EntryMod defines a function type which receives a pointer to an entry.
type EntryMod func(*Entry)

// Partial returns a new EntryMod which will always apply provided EntryMod
// to all provided Entry.
func Partial(mods ...EntryMod) EntryMod {
	if len(mods) == 1 {
		return mods[0]
	}

	return func(en *Entry) {
		for _, mod := range mods {
			mod(en)
		}
	}
}

// Apply runs all giving EntryMod functions provided on the provided Entry.
func Apply(en *Entry, mods ...EntryMod) {
	for _, mod := range mods {
		mod(en)
	}
}

// Timelapse defines a message attached with a giving time value.
type Timelapse struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
	Field   Field     `json:"fields"`
}

// WithTimelapse returns a Timelapse with associated field and message.
func WithTimelapse(message string, f Field) EntryMod {
	return func(en *Entry) {
		en.Timelapse = append(en.Timelapse, Timelapse{
			Field:   f,
			Message: message,
			Time:    time.Now(),
		})
	}
}

// YellowAlert returns an Entry with the level set to YellowAlertLvl.
func YellowAlert(err error, message string, m ...interface{}) EntryMod {
	return Partial(withMessageAt(4, YellowAlertLvl, message, m...), func(en *Entry) {
		en.Field["error"] = err
	})
}

// RedAlert returns an Entry with the level set to RedAlertLvl.
func RedAlert(err error, message string, m ...interface{}) EntryMod {
	return Partial(withMessageAt(4, RedAlertLvl, message, m...), func(en *Entry) {
		en.Field["error"] = err
	})
}

// Errorf returns a entry where the message is the provided error.Error() value
// produced from the message and its provided values
// and the error is added as a key-value within the Entry fields.
func Errorf(message string, m ...interface{}) EntryMod {
	err := fmt.Errorf(message, m...)
	return Partial(withMessageAt(4, ErrorLvl, err.Error()), func(en *Entry) {
		en.Field["error"] = err
	})
}

// Error returns a entry where the message is the provided error.Error() value
// and the error is added as a key-value within the Entry fields.
func Error(err error) EntryMod {
	return Partial(withMessageAt(4, ErrorLvl, err.Error()), func(en *Entry) {
		en.Field["error"] = err
	})
}

// Tags returns an Entry with the tags value set to ts.
func Tags(ts ...string) EntryMod {
	return func(en *Entry) {
		en.Tags = ts
	}
}

// Type returns an Entry with the type value set to t.
func Type(t string) EntryMod {
	return func(en *Entry) {
		en.Type = t
	}
}

// Info returns an Entry with the level set to Info.
func Info(message string, m ...interface{}) EntryMod {
	return withMessageAt(4, InfoLvl, message, m...)
}

// Message returns a new Entry with the provided Level and message used.
func Message(message string, m ...interface{}) EntryMod {
	function, file, line := getFunctionName(3)
	return func(en *Entry) {
		en.Message = fmt.Sprintf(message, m...)
		en.Function, en.File, en.Line = function, file, line
	}
}

// WithMessage returns a new Entry with the provided Level and message used.
func WithMessage(level Level, message string, m ...interface{}) EntryMod {
	return withMessageAt(4, level, message, m...)
}

// withMessage returns a new Entry with the provided Level and message used.
func withMessageAt(depth int, level Level, message string, m ...interface{}) EntryMod {
	function, file, line := getFunctionName(depth)
	return func(e *Entry) {
		e.Level = level
		e.Field = make(Field)
		e.Time = time.Now()
		e.Function, e.File, e.Line = function, file, line

		if len(m) == 0 {
			e.Message = message
			return
		}
		e.Message = fmt.Sprintf(message, m...)
	}
}

// WithTrace returns itself after setting the giving trace value
// has the method trace for the giving Entry.
func WithTrace(t Trace) EntryMod {
	function, file, line := getFunctionName(3)
	return func(en *Entry) {
		en.Trace = t
		en.Function, en.File, en.Line = function, file, line
	}
}

// WithFilter returns a Entry and set the Filter to the provided value.
func WithFilter(filter interface{}) EntryMod {
	return func(en *Entry) {
		en.Filter = filter
	}
}

// WithID returns a Entry and set the ID to the provided value.
func WithID(id string) EntryMod {
	function, file, line := getFunctionName(3)
	return func(en *Entry) {
		en.ID = id
		en.Function, en.File, en.Line = function, file, line
	}
}

// With returns a Entry set to the LogLevel of the previous and
// adds the giving key-value pair to the entry.
func With(key string, value interface{}) EntryMod {
	return func(en *Entry) {
		if en.Field == nil {
			en.Field = make(Field)
		}

		en.Field[key] = value
	}
}

// WithFields adds all field key-value pair into associated Entry
// returning the Entry.
func WithFields(f Field) EntryMod {
	return func(en *Entry) {
		if en.Field == nil {
			en.Field = make(Field)
		}

		for k, v := range f {
			en.Field[k] = v
		}
	}
}
