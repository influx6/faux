package stacks

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// contains stack space-delimiters controls.
const ()

// Tracer provides a interface that allows adding stack trace details, which
// provides a flexible registry formatting how stack trace gets collected.
type Tracer interface {
	Add(context string, Line int, Goroutine, PkgRoot, Pkg, File, Method, MethodAddress, Address, OriginalLine string)
}

// Run collects the stack trace and compiles into the giving trace store.
// If all is false, only the current go-routine is checked.
// The size determines the size of the stacktrace []byte list
// t Trace provides the trace storage mechanism,desired.
func Run(context string, all bool, size int, t Tracer) {
	traceLines := PullTrace(all, size)

	var goroutine, addr, root, pkg, file, method, mName string
	var lineNumber int

	for index, line := range traceLines {
		if len(line) == 0 {
			continue
		}

		// if we contain the 'goroutine' word, then this is the go-routine info,
		// cache and skip.
		if strings.Contains(strings.ToLower(line), "goroutine") {
			goroutine = line
			continue
		}

		if index%2 != 0 {
			root, pkg, method, mName = ExtractPackageAndMethod(line)
		} else {
			addr, file, lineNumber = ExtractFileAndLineNumber(line)
			t.Add(context, lineNumber, goroutine, root, pkg, file, mName, method, addr, line)
		}
	}
}

// Stack represents a piece of stack strace defining a specific line within a
// received stack.
type Stack struct {
	Context    string
	LineNumber int
	MethodName string
	Method     string
	File       string
	Root       string
	Package    string
	Address    string
	Goroutine  string
	Line       string
}

// Trace provides a lists of stack traces which gets registered from a trace
// lists.
type Trace []Stack

// TraceNow returns a new trace of the current goroutine unless all is set to true.
// Returns a list of trace Stacks.
func TraceNow(context string, all bool, size int) Trace {
	tr := make(Trace, 0)
	Run(context, all, size, &tr)
	return tr
}

// FindMethod returns the stack trace which has the specific method information.
// If the index is -1 then the stack trace was not found.
func (s Trace) FindMethod(method string) (Stack, int) {
	var stack Stack
	var index = -1

	for ind, us := range s {
		if us.MethodName != method {
			continue
		}
		stack = us
		index = ind
		break
	}

	return stack, index
}

// Add adds the stack trace information into the list of traces.
func (s *Trace) Add(context string, line int, gr, pkgroot, pkg, file, methodName, method, addr, lm string) {
	*s = append(*s, Stack{
		Context:    context,
		LineNumber: line,
		Method:     method,
		MethodName: methodName,
		File:       file,
		Root:       pkgroot,
		Package:    pkg,
		Address:    addr,
		Goroutine:  gr,
		Line:       lm,
	})
}

// PullTrace pulls the current stack strace, if the all bool is true, gets all
// goroutine stack traces else only the current go-routine it runs gets traced.
// Uses runtime.Stack underneath. Uses the size to limit the list returned.
// If size is zero, it defaults to a 1 << 16 stack list length.
func PullTrace(all bool, size int) []string {
	if size == 0 {
		size = 1 << 16
	}

	trace := make([]byte, size)
	trace = trace[:runtime.Stack(trace, all)]
	return strings.Split(string(trace), "\n")
}

var methodInfo = regexp.MustCompile(`(\(.+?\))`)

// ExtractPackageAndMethod extracts the package root, pkg path and method from
// a trace line.
func ExtractPackageAndMethod(line string) (string, string, string, string) {
	var pkgRoot, pkg, method, methodName string
	pkgRoot, method = splitPath(line)

	parts := strings.Split(method, ".")
	pkg = parts[0]

	if len(pkgRoot) > 0 {
		pkg = fmt.Sprintf("%s/%s", pkgRoot, parts[0])
	}

	methodName = parts[len(parts)-1]
	method = strings.Join(parts[1:], ".")

	if methodInfo.MatchString(methodName) {
		methodName = methodInfo.ReplaceAllString(methodName, "")
	} else {
		methodName = strings.Replace(methodName, "()", "", -1)
	}

	return pkgRoot, pkg, method, methodName
}

// ExtractFileAndLineNumber extracts the method name by which package, the line
// numbe and other meta information for the stack trace.
func ExtractFileAndLineNumber(line string) (string, string, int) {
	var lineNumber int
	var file string

	_, methodMeta := splitPath(line)
	fileMeta, addr := nameAndAddress(methodMeta)
	parts := strings.Split(fileMeta, ":")

	file = parts[0]

	if len(parts) > 1 {
		lnx, _ := strconv.ParseUint(parts[1], 10, 32)
		lineNumber = int(lnx)
	}

	return addr, file, lineNumber
}

// splitAndLastSlash splits a string by a formward slash and returns the left
// and right parts of it.
func splitPath(line string) (string, string) {
	parts := strings.Split(line, "/")
	partsLen := len(parts)
	right := parts[partsLen-1]
	left := strings.Join(parts[:partsLen-1], "/")
	return left, right
}

// removeAllSpace removes uneccessary spaces found and returns the giving string
func nameAndAddress(line string) (string, string) {
	parts := strings.Split(line, " ")
	addr := parts[len(parts)-1]
	return strings.Join(parts[:len(parts)-1], " "), addr
}
