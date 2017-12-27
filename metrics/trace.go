package metrics

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// vars
var (
	stackSize = 1 << 6
	question  = "???"
)

// Trace defines a structure which contains the stack, start and endtime
// on a given from a trace call to trace a given call with stack details
// and execution time.
type Trace struct {
	File       string    `json:"file"`
	Package    string    `json:"Package"`
	Function   string    `json:"function"`
	LineNumber int       `json:"line_number"`
	Stack      []byte    `json:"stack"`
	Comments   []string  `json:"comments"`
	Time       time.Time `json:"end_time"`
}

// NewTrace returns a Trace object which is used to track the execution and
// stack details of a given trace call.
func NewTrace(comments ...string) Trace {
	trace := make([]byte, stackSize)
	trace = trace[:runtime.Stack(trace, false)]

	_, file, line, ok := runtime.Caller(1)
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

	functionName, _, _ := getFunctionName(3)

	return Trace{
		Package:    pkg,
		LineNumber: line,
		Stack:      trace,
		Comments:   comments,
		Time:       time.Now(),
		File:       pkgFile,
		Function:   functionName,
	}

}

// NewTraceWithCallDepth returns a Trace object which is used to track the execution and
// stack details of a given trace call.
func NewTraceWithCallDepth(depth int, comments ...string) Trace {
	trace := make([]byte, stackSize)
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

	functionName, _, _ := getFunctionName(3)

	return Trace{
		Package:    pkg,
		LineNumber: line,
		Stack:      trace,
		File:       pkgFile,
		Comments:   comments,
		Time:       time.Now(),
		Function:   functionName,
	}
}

// String returns the giving trace timestamp for the execution time.
func (t Trace) String() string {
	return fmt.Sprintf("[Package=%q, File=%q, Time=%+q]", t.Package, t.File, t.Time)
}

// getFunctionName returns the caller of the function that called it :)
func getFunctionName(depth int) (string, string, int) {
	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(depth, fpcs)
	if n == 0 {
		return "Unknown()", "???", 0
	}

	funcPtr := fpcs[0]
	funcPtrArea := funcPtr - 1

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(funcPtrArea)
	if fun == nil {
		return "Unknown()", "???", 0
	}

	fileName, line := fun.FileLine(funcPtrArea)

	// return its name
	return fun.Name(), fileName, line
}
