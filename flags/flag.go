// Package flags creates a augmentation on the native flags package.
package flags

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/influx6/faux/context"
)

const (
	usageTml = `Usage: {{ toLower .Title}} [flags] [command] 

⡿ COMMANDS:{{ range .Commands }}
	⠙ {{toLower .Name }}	{{if isEmpty .ShortDesc }}{{cutoff .Desc 40 }}{{else}}{{cutoff .ShortDesc 40 }}{{end}}
{{end}}

⡿ HELP:
	Run [command] help

⡿ OTHERS:
	Run '{{toLower .Title}} printflags' to print all flags of all commands.

⡿ WARNING:
	Uses internal flag package so flags must precede command name. 
	e.g '{{toLower .Title}} -cmd.flag=4 run'
`

	cmdUsageTml = `Command: {{toLower .Title}} [flags] {{ toLower .Cmd.Name}} 

⡿ DESC:
	{{.Cmd.Desc}}

⡿ Flags:
	{{$title := toLower .Title}}{{$cmdName := .Cmd.Name}}{{ range $_, $fl := .Cmd.Flags }}
	⠙ {{toLower $cmdName}}.{{toLower $fl.FlagName}}		Default: {{.Default}}	Desc: {{.Desc }}
	{{end}}

⡿ USAGE:
	{{ range $_, $fl := .Cmd.Flags }}
	⠙ {{$title}} -{{toLower $cmdName}}.{{toLower $fl.FlagName}}={{.Default}} {{toLower $cmdName}} 
	{{end}}

⡿ OTHERS:
	Commands which respect context.Context, can set timeout by using the -timeout flag.
	e.g -timeout=4m, -timeout=4h

⡿ WARNING:
	Uses internal flag package so flags must precede command name. 
	e.g '{{toLower .Title}} -cmd.flag=4 run'
`
)

var (
	timeout = flag.Duration("timeout", 0, "-timeout=4m to set deadline for function execution")

	defs = template.FuncMap{
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
		"isEmpty": func(val string) bool {
			return strings.TrimSpace(val) == ""
		},
		"cutoff": func(val string, limit int) string {
			if len(val) > limit {
				return val[:limit]
			}
			return val
		},
	}
)

// Flag defines a interface exposing a single function for parsing
// a giving flag for attaching and data collection.
type Flag interface {
	FlagName() string
	Value() interface{}
	Parse(string) error
	DefaultValue() interface{}
}

// DurationFlag implements a structure for parsing duration flags.
type DurationFlag struct {
	Name       string
	Desc       string
	Default    time.Duration
	value      *time.Duration
	Validation func(time.Duration) error
}

// FlagName returns name of flag.
func (s *DurationFlag) FlagName() string {
	return s.Name
}

// DefaultValue returns default value of flag pointer.
func (s *DurationFlag) DefaultValue() interface{} {
	return s.Default
}

// Value returns internal value of flag pointer.
func (s *DurationFlag) Value() interface{} {
	return *s.value
}

// Parse sets the underline flag ready for value receiving.
func (s *DurationFlag) Parse(cmd string) error {
	s.value = new(time.Duration)
	flag.DurationVar(s.value, fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), s.Default, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// Float64Flag implements a structure for parsing float64 flags.
type Float64Flag struct {
	Name       string
	Desc       string
	Default    float64
	value      *float64
	Validation func(float64) error
}

// FlagName returns name of flag.
func (s *Float64Flag) FlagName() string {
	return s.Name
}

// DefaultValue returns default value of flag pointer.
func (s *Float64Flag) DefaultValue() interface{} {
	return s.Default
}

// Value returns internal value of flag pointer.
func (s *Float64Flag) Value() interface{} {
	return *s.value
}

// Parse sets the underline flag ready for value receiving.
func (s *Float64Flag) Parse(cmd string) error {
	s.value = new(float64)
	flag.Float64Var(s.value, fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), s.Default, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// UInt64Flag implements a structure for parsing uint64 flags.
type UInt64Flag struct {
	Name       string
	Desc       string
	Default    uint64
	value      *uint64
	Validation func(uint64) error
}

// FlagName returns name of flag.
func (s *UInt64Flag) FlagName() string {
	return s.Name
}

// DefaultValue returns default value of flag pointer.
func (s *UInt64Flag) DefaultValue() interface{} {
	return s.Default
}

// Value returns internal value of flag pointer.
func (s *UInt64Flag) Value() interface{} {
	return *s.value
}

// Parse sets the underline flag ready for value receiving.
func (s *UInt64Flag) Parse(cmd string) error {
	s.value = new(uint64)
	flag.Uint64Var(s.value, fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), s.Default, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// Int64Flag implements a structure for parsing int64 flags.
type Int64Flag struct {
	Name       string
	Desc       string
	Default    int64
	value      *int64
	Validation func(int64) error
}

// FlagName returns name of flag.
func (s *Int64Flag) FlagName() string {
	return s.Name
}

// Value returns internal value of flag pointer.
func (s *Int64Flag) Value() interface{} {
	return *s.value
}

// DefaultValue returns default value of flag pointer.
func (s *Int64Flag) DefaultValue() interface{} {
	return s.Default
}

// Parse sets the underline flag ready for value receiving.
func (s *Int64Flag) Parse(cmd string) error {
	s.value = new(int64)
	flag.Int64Var(s.value, fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), s.Default, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// UIntFlag implements a structure for parsing uint flags.
type UIntFlag struct {
	Name       string
	Desc       string
	Default    uint
	value      *uint
	Validation func(uint) error
}

// Parse sets the underline flag ready for value receiving.
func (s *UIntFlag) Parse(cmd string) error {
	s.value = new(uint)
	flag.UintVar(s.value, fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), s.Default, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// FlagName returns name of flag.
func (s *UIntFlag) FlagName() string {
	return s.Name
}

// DefaultValue returns default value of flag pointer.
func (s *UIntFlag) DefaultValue() interface{} {
	return s.Default
}

// Value returns internal value of flag pointer.
func (s *UIntFlag) Value() interface{} {
	return *s.value
}

// IntFlag implements a structure for parsing int flags.
type IntFlag struct {
	Name       string
	Desc       string
	Default    int
	value      *int
	Validation func(int) error
}

// FlagName returns name of flag.
func (s *IntFlag) FlagName() string {
	return s.Name
}

// Value returns internal value of flag pointer.
func (s *IntFlag) Value() interface{} {
	return *s.value
}

// DefaultValue returns default value of flag pointer.
func (s *IntFlag) DefaultValue() interface{} {
	return s.Default
}

// Parse sets the underline flag ready for value receiving.
func (s *IntFlag) Parse(cmd string) error {
	s.value = flag.Int(fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), s.Default, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// BoolFlag implements a structure for parsing bool flags.
type BoolFlag struct {
	Name       string
	Desc       string
	Default    bool
	value      *bool
	Validation func(bool) error
}

// FlagName returns name of flag.
func (s *BoolFlag) FlagName() string {
	return s.Name
}

// Value returns internal value of flag pointer.
func (s *BoolFlag) Value() interface{} {
	return *s.value
}

// DefaultValue returns default value of flag pointer.
func (s *BoolFlag) DefaultValue() interface{} {
	return s.Default
}

// Parse sets the underline flag ready for value receiving.
func (s *BoolFlag) Parse(cmd string) error {
	s.value = new(bool)
	flag.BoolVar(s.value, fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), s.Default, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// TBoolFlag implements a structure for parsing bool flags that are true by default.
type TBoolFlag struct {
	Name       string
	Desc       string
	Default    bool
	value      *bool
	Validation func(bool) error
}

// DefaultValue returns default value of flag pointer.
func (s *TBoolFlag) DefaultValue() interface{} {
	return true
}

// FlagName returns name of flag.
func (s *TBoolFlag) FlagName() string {
	return s.Name
}

// Value returns internal value of flag pointer.
func (s *TBoolFlag) Value() interface{} {
	return *s.value
}

// Parse sets the underline flag ready for value receiving.
func (s *TBoolFlag) Parse(cmd string) error {
	s.Default = true
	s.value = new(bool)
	flag.BoolVar(s.value, fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), true, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// StringFlag implements a structure for parsing string flags.
type StringFlag struct {
	Name       string
	Desc       string
	Default    string
	value      *string
	Validation func(string) error
}

// FlagName returns name of flag.
func (s *StringFlag) FlagName() string {
	return s.Name
}

// DefaultValue returns default value of flag pointer.
func (s *StringFlag) DefaultValue() interface{} {
	return s.Default
}

// Value returns internal value of flag pointer.
func (s *StringFlag) Value() interface{} {
	return *s.value
}

// Parse sets the underline flag ready for value receiving.
func (s *StringFlag) Parse(cmd string) error {
	s.value = new(string)
	flag.StringVar(s.value, fmt.Sprintf("%s.%s", strings.ToLower(cmd), s.Name), s.Default, s.Desc)
	if s.Validation != nil {
		return s.Validation(*s.value)
	}
	return nil
}

// Action defines a giving function to be executed for a Command.
type Action func(context.Context) error

// Command defines structures which define specific actions to be executed
// with associated flags.
type Command struct {
	Name        string
	Desc        string
	ShortDesc   string
	Flags       []Flag
	Action      Action
	WaitOnCtrlC bool
}

// Run adds all commands and appropriate flags for each commands.
// There is no need to call flag.Parse, has this calls it underneath and
// parses appropriate commands.
func Run(title string, cmds ...Command) {
	commandHelp := make(map[string]string)

	if tml, err := template.New("flags.Usage").Funcs(defs).Parse(usageTml); err == nil {
		var bu bytes.Buffer
		if err := tml.Execute(&bu, struct {
			Title    string
			Commands []Command
		}{
			Title:    title,
			Commands: cmds,
		}); err == nil {
			flag.Usage = func() {
				fmt.Println(bu.String())
			}
		}
	}

	// Register all flags first.
	for _, cmd := range cmds {
		if tml, err := template.New("command.Usage").Funcs(defs).Parse(cmdUsageTml); err == nil {
			var bu bytes.Buffer
			if err := tml.Execute(&bu, struct {
				Title string
				Cmd   Command
			}{
				Title: title,
				Cmd:   cmd,
			}); err == nil {
				commandHelp[cmd.Name] = bu.String()
			} else {
				commandHelp[cmd.Name] = err.Error()
			}
		}

		for _, flag := range cmd.Flags {
			if err := flag.Parse(cmd.Name); err != nil {
				log.Fatalf("Flags error: %+q : %+s", flag.FlagName(), err)
				return
			}
		}
	}

	flag.Parse()

	command := strings.ToLower(flag.Arg(0))
	subCommand := strings.ToLower(flag.Arg(1))
	if command == "printflags" {
		flag.PrintDefaults()
		return
	}

	var wg sync.WaitGroup

	for _, cmd := range cmds {
		if strings.ToLower(cmd.Name) == command {
			if subCommand != "help" {
				valCtx := context.NewValueBag()
				for _, flag := range cmd.Flags {
					valCtx.Set(flag.FlagName(), flag.Value())
				}

				var ctx context.CancelableContext
				if *timeout == 0 {
					ctx = context.NewCnclContext(valCtx)
				} else {
					ctx = context.NewExpiringCnclContext(nil, *timeout, valCtx)
				}

				if !cmd.WaitOnCtrlC {
					if err := cmd.Action(ctx); err != nil {
						fmt.Fprint(os.Stderr, err.Error())
					}
					return
				}

				ch := make(chan os.Signal, 3)
				signal.Notify(ch, syscall.SIGQUIT)
				signal.Notify(ch, syscall.SIGTERM)
				signal.Notify(ch, os.Interrupt)

				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := cmd.Action(ctx); err != nil {
						fmt.Fprint(os.Stderr, err.Error())
						close(ch)
					}
				}()

				<-ch
				ctx.Cancel()
				return
			}

			fmt.Println(commandHelp[cmd.Name])
			return
		}
	}

	if flag.Usage != nil {
		flag.Usage()
	}

	wg.Wait()
}
