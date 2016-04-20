package flux

import "fmt"

// Provides simple printer printers

// succeedMark is the Unicode codepoint for a check mark.
const succeedMark = "\u2713"

// failedMark is the Unicode codepoint for an X mark.
const failedMark = "\u2717"

//Passed returns a msg with a check mark
func Passed(msg string, v ...interface{}) string {
	return render(succeedMark, msg, v...)
}

//Failed returns a msg with a x mark
func Failed(msg string, v ...interface{}) string {
	return render(failedMark, msg, v...)
}

// FatalPrinter provides a hasher,stop after print interface with Fatal function
type FatalPrinter interface {
	Fatal(v ...interface{})
}

//FatalFailed uses the log to print out the failed message
func FatalFailed(fr FatalPrinter, msg string, v ...interface{}) {
	fr.Fatal(Failed(msg, v...))
}

//FatalPassed uses the log to print out the passed message
func FatalPassed(fr FatalPrinter, msg string, v ...interface{}) {
	fr.Fatal(Passed(msg, v...))
}

// LogPrinter provides a simple printer interface with normal log function
type LogPrinter interface {
	Log(v ...interface{})
}

//LogFailed uses the log to print out the failed message
func LogFailed(pr LogPrinter, msg string, v ...interface{}) {
	pr.Log(Failed(msg, v...))
}

//LogPassed uses the log to print out the passed message
func LogPassed(pr LogPrinter, msg string, v ...interface{}) {
	pr.Log(Passed(msg, v...))
}

// SimplePrinter provides a simple printer interface with normal print function
type SimplePrinter interface {
	Print(v ...interface{})
}

//PrintFailed uses the log to print out the failed message
func PrintFailed(pr SimplePrinter, msg string, v ...interface{}) {
	pr.Print(Passed(msg, v...))
}

//PrintPassed uses the log to print out the passed message
func PrintPassed(pr SimplePrinter, msg string, v ...interface{}) {
	pr.Print(Passed(msg, v...))
}

//Render returns a new string with the mark put in place
func render(mark, msg string, v ...interface{}) string {
	rms := fmt.Sprintf("\t%s: %s", msg, mark)
	return fmt.Sprintf(rms, v...)
}
