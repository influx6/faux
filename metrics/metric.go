// Package metrics defines a basic structure foundation for handling logs without
// much hassle, allow more different entries to be created.
// Inspired by https://medium.com/@tjholowaychuk/apex-log-e8d9627f4a9a.
package metrics

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// TraceKey defines the key which is used to store the trace object.
	TraceKey = "FuncTrace"

	// DefaultMessage defines a default message used by SentryJSON where
	// fields contains no messages to be used.
	DefaultMessage = "No Message"

	// StackSize defines the max size for an expected stack.
	StackSize = 1 << 6
)

//==============================================================================

// Metrics defines an interface which exposes a method which receives given
// Entry which will be sorted accordingly to it's registered entry.
type Metrics interface {
	Emit(Entry) error
}

// Sentry exposes an interface which allows Entries to be transformed into
// a structure which delivers the json data to remote APIs, services, etc.
type Sentry interface {
	Emit(SentryJSON) error
}

//==============================================================================

// SentryJSON defines a json style structure for delivery entry data to
// other APIs.
type SentryJSON struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Fields  Fields    `json:"fields"`
}

// SentryPipe defines a pipe which will expose a method to allow piping into a
// metrics to deliver entries as centries.
type SentryPipe struct {
	sentries []Sentry
}

// Sentries returns a metrics which will pipe all recieved Entrys to provided
// sentries.
func Sentries(sx ...Sentry) Metrics {
	return SentryPipe{
		sentries: sx,
	}
}

// Emit delivers the giving entry to all available metricss.
func (pipe SentryPipe) Emit(e Entry) error {
	var sentryJSON SentryJSON
	sentryJSON.Fields = e.Fields()
	sentryJSON.Time = time.Now()

	var message string
	if e.Message != "" {
		message = e.Message
	} else if mo, ok := sentryJSON.Fields["message"].(string); ok {
		message = mo
	} else {
		message = DefaultMessage
	}

	sentryJSON.Message = message

	for _, sentry := range pipe.sentries {
		if err := sentry.Emit(sentryJSON); err != nil {
			return err
		}
	}

	return nil
}

//==============================================================================

// FilterecMaster defines that filters out Entry objects based on a provided function.
type FilteredMaster struct {
	Master
	filterFn func(Entry) bool
}

// Filter returns a new instance of a FilteredMaster.
func Filter(filterFn func(Entry) bool, metrics ...interface{}) FilteredMaster {
	return FilteredMaster{
		filterFn: filterFn,
		Master:   New(metrics...),
	}
}

// Emit delivers the giving entry to all available metricss.
func (fm FilteredMaster) Emit(e Entry) error {
	if !fm.filterFn(e) {
		return nil
	}

	return fm.Master.Emit(e)
}

// Master defines a core metrics structure to pipe Entry values to registed metricss.
type Master struct {
	metrics []Metrics
}

// New returns a new metricsMaster for which will recieve all expected Entry values.
func New(metrics ...interface{}) Master {
	var sentries []Sentry
	var entries []Metrics

	for _, item := range metrics {
		switch rItem := item.(type) {
		case Metrics:
			entries = append(entries, rItem)
		case Sentry:
			sentries = append(sentries, rItem)
		}
	}

	return Master{
		metrics: append(entries, Sentries(sentries...)),
	}
}

// With returns a new Master with a new list of metricss.
func (metrics Master) With(m Metrics) Master {
	return Master{
		metrics: append([]Metrics{m}, metrics.metrics...),
	}
}

// Emit delivers the giving entry to all available metricss.
func (metrics Master) Emit(e Entry) error {
	for _, metrics := range metrics.metrics {
		if err := metrics.Emit(e); err != nil {
			return err
		}
	}

	return nil
}

//==============================================================================

// NilPair defines a nil starting pair.
var NilPair = (*Pair)(nil)

// Pair defines a struct for storing a linked pair of key and values.
type Pair struct {
	prev  *Pair
	key   string
	value interface{}
}

// NewPair returns a a key-value pair chain for setting fields.
func NewPair(key string, value interface{}) *Pair {
	return &Pair{
		key:   key,
		value: value,
	}
}

// Append returns a new Pair with the giving key and with the provded Pair set as
// it's previous link.
func Append(p *Pair, key string, value interface{}) *Pair {
	return p.Append(key, value)
}

// Fields defines a type for key-value pairs which defines the field values to be stored.
type Fields map[string]interface{}

// Fields returns all internal pair data as a map.
func (p *Pair) Fields() Fields {
	var f Fields

	if p.prev == nil {
		f = make(Fields)
		f[p.key] = p.value
		return f
	}

	f = p.prev.Fields()

	if p.key != "" {
		f[p.key] = p.value
	}

	return f
}

// Append returns a new pair with the giving key and value and its previous
// set to this pair.
func (p *Pair) Append(key string, val interface{}) *Pair {
	return &Pair{
		prev:  p,
		key:   key,
		value: val,
	}
}

// Root returns the root Pair in the chain which links all pairs together.
func (p *Pair) Root() *Pair {
	if p.prev == nil {
		return p
	}

	return p.prev.Root()
}

// GetBool collects the string value of a key if it exists.
func (p *Pair) GetBool(key string) (bool, bool) {
	val, found := p.Get(key)
	if !found {
		return false, false
	}

	value, ok := val.(bool)
	return value, ok
}

// GetFloat64 collects the string value of a key if it exists.
func (p *Pair) GetFloat64(key string) (float64, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(float64)
	return value, ok
}

// GetFloat32 collects the string value of a key if it exists.
func (p *Pair) GetFloat32(key string) (float32, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(float32)
	return value, ok
}

// GetInt8 collects the string value of a key if it exists.
func (p *Pair) GetInt8(key string) (int8, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int8)
	return value, ok
}

// GetInt16 collects the string value of a key if it exists.
func (p *Pair) GetInt16(key string) (int16, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int16)
	return value, ok
}

// GetInt64 collects the string value of a key if it exists.
func (p *Pair) GetInt64(key string) (int64, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int64)
	return value, ok
}

// GetInt32 collects the string value of a key if it exists.
func (p *Pair) GetInt32(key string) (int32, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int32)
	return value, ok
}

// GetInt collects the string value of a key if it exists.
func (p *Pair) GetInt(key string) (int, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int)
	return value, ok
}

// GetString collects the string value of a key if it exists.
func (p *Pair) GetString(key string) (string, bool) {
	val, found := p.Get(key)
	if !found {
		return "", false
	}

	value, ok := val.(string)
	return value, ok
}

// Get collects the value of a key if it exists.
func (p *Pair) Get(key string) (value interface{}, found bool) {
	if p == nil {
		return
	}

	if p.key == key {
		return p.value, true
	}

	if p.prev == nil {
		return
	}

	return p.prev.Get(key)
}

//==============================================================================

// Entry defines a data type which encapuslates data related to a giving
// Log event.
type Entry struct {
	*Pair
	Message string
}

// WithFields returns a new try with the provided key-value pair with the set ID.
func WithFields(f Fields) Entry {
	entry := Entry{
		Pair: (*Pair)(nil),
	}

	for k, v := range f {
		entry.Pair = entry.Pair.Append(k, v)
	}

	return entry
}

// With returns a new try with the provided key-value pair with the set ID.
func With(key string, value interface{}) Entry {
	return Entry{
		Pair: NewPair(key, value),
	}
}

// Trace defines a structure which contains the stack, start and endtime
// on a given from a trace call to trace a given call with stack details
// and execution time.
type Trace struct {
	File       string    `json:"file"`
	Package    string    `json:"Package"`
	Function   string    `json:"function"`
	LineNumber int       `json:"line_number"`
	BeginStack []byte    `json:"begin_stack"`
	EndStack   []byte    `json:"end_stack"`
	Comments   []string  `json:"comments"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	entry      *Entry
}

// String returns the giving trace timestamp for the execution time.
func (t *Trace) String() string {
	return fmt.Sprintf("[Total=%+q, Start=%+q, End=%+q]", t.EndTime.Sub(t.StartTime), t.StartTime, t.EndTime)
}

// End stops the trace, captures the current stack trace and returns the
// entry related to the trace.
func (t *Trace) End() Entry {
	trace := make([]byte, StackSize)
	trace = trace[:runtime.Stack(trace, false)]

	entry := t.entry
	t.entry = nil

	t.EndStack = trace
	t.EndTime = time.Now()

	return entry.With(TraceKey, *t)
}

var question = "???"

// TraceWithCallDepth returns a Trace object which is used to track the execution and
// stack details of a given trace call.
func (e Entry) TraceWithCallDepth(depth int, comments ...string) *Trace {
	trace := make([]byte, StackSize)
	trace = trace[:runtime.Stack(trace, false)]

	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = question
	}

	var pkg, pkgFile string
	pkgFileBase := file

	if file != question {
		pkgPieces := strings.SplitAfter(pkgFileBase, "/src/")
		if len(pkgPieces) > 1 {
			pkgFileBase = pkgPieces[1]
		}

		pkg = filepath.Dir(pkgFileBase)
		pkgFile = filepath.Base(pkgFileBase)
	}

	return &Trace{
		entry:      &e,
		Package:    pkg,
		LineNumber: line,
		BeginStack: trace,
		File:       pkgFile,
		Comments:   comments,
		StartTime:  time.Now(),
		Function:   GetFunctionName(),
	}
}

// Trace returns a Trace object which is used to track the execution and
// stack details of a given trace call.
func (e Entry) Trace(comments ...string) *Trace {
	trace := make([]byte, StackSize)
	trace = trace[:runtime.Stack(trace, false)]

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
	}

	var pkg, pkgFile string
	pkgFileBase := file

	if file != "???" {
		pkgPieces := strings.SplitAfter(pkgFileBase, "/src/")
		if len(pkgPieces) > 1 {
			pkgFileBase = pkgPieces[1]
		}

		pkg = filepath.Dir(pkgFileBase)
		pkgFile = filepath.Base(pkgFileBase)
	}

	return &Trace{
		entry:      &e,
		Package:    pkg,
		LineNumber: line,
		BeginStack: trace,
		Comments:   comments,
		StartTime:  time.Now(),
		File:       pkgFile,
		Function:   GetFunctionName(),
	}
}

// With returns a new Entry set to the LogLevel of the previous and
// adds the giving key-value pair to the entry.
func (e Entry) With(key string, value interface{}) Entry {
	return Entry{
		Pair:    e.Pair.Append(key, value),
		Message: e.Message,
	}
}

// WithMessage sets the message for the giving Entry if it has no message
// else returns a new Entry with the set message.
func (e Entry) WithMessage(message string, m ...interface{}) Entry {
	if len(m) != 0 {
		e.Message = fmt.Sprintf(message, m...)
		return e
	}

	return Entry{
		Pair:    e.Pair,
		Message: fmt.Sprintf(message, m...),
	}
}

// WithFields returns a new Entry set to the LogLevel of the previous and
// adds the all giving key-value pair from the Fields to the entry.
func (e Entry) WithFields(f Fields) Entry {
	entry := Entry{
		Pair:    e.Pair,
		Message: e.Message,
	}

	for k, v := range f {
		entry.Pair = entry.Pair.Append(k, v)
	}

	return entry
}

//==============================================================================

// Hide takes the given message and generates a '***' character sets.
func Hide(message string) string {
	mLen := len(message)

	var mval []string

	for i := 0; i < mLen; i++ {
		mval = append(mval, "*")
	}

	return strings.Join(mval, "")
}

//==============================================================================

// GetFunctionName returns the caller of the function that called it :)
func GetFunctionName() string {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "Unknown()" // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "Unknown()" // proper error her would be better
	}

	// return its name
	return fun.Name()
}
