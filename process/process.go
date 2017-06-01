package process

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
)

var log = sink.New(sinks.Stdout{})

// CriticalLevel defines a int type which is used to signal the critical nature of a
// command/script to be executed.
type CriticalLevel int

// Contains possible critical level values for commands execution
const (
	Normal CriticalLevel = iota + 1
	Warning
	RedAlert
	SilentKill
)

// Command defines the command to be executed and it's arguments
type Command struct {
	Name  string        `json:"name" toml:"name"`
	Level CriticalLevel `json:"level" toml:"level"`
	Args  []string      `json:"args" toml:"args"`
	Async bool          `json:"async" toml:"async"`
}

// Run executes the giving command and returns the bytes.Buffer for both
// the Stdout and Stderr.
func (c Command) Run(ctx context.Context, wout, werr io.Writer, pin io.Reader) error {
	proc := exec.Command(c.Name, c.Args...)
	proc.Stdout = wout
	proc.Stdin = pin
	proc.Stderr = werr

	if err := proc.Start(); err != nil {
		log.Emit(sinks.Error("Process : Error : Command : Begin Execution : %q : %q", c.Name, c.Args))
		return err
	}

	go func() {
		<-ctx.Done()
		if proc.Process != nil {
			proc.Process.Kill()
		}
	}()

	if !c.Async {
		if err := proc.Wait(); err != nil {
			log.Emit(sinks.Error("Process : Error : Command : Begin Execution : %q : %q", c.Name, c.Args))

			if c.Level > Warning {
				return err
			}

			return nil
		}
	}

	return nil
}

//============================================================================================

// SyncProcess defines a struct which is used to execute a giving set of
// script values.
type SyncProcess struct {
	Commands []Command `json:"commands"`
}

// Exec executes the giving series of commands attached to the
// process.
func (p SyncProcess) Exec(ctx context.Context, pipeOut, pipeErr io.Writer, pipeIn io.Reader) error {
	for _, command := range p.Commands {
		if err := command.Run(ctx, pipeOut, pipeErr, pipeIn); err != nil {
			return err
		}
	}

	return nil
}

//============================================================================================

// AsyncProcess defines a struct which is used to execute a giving set of
// script values.
type AsyncProcess struct {
	Commands []Command `json:"commands"`
}

// Exec executes the giving series of commands attached to the
// process.
func (p AsyncProcess) Exec(ctx context.Context, pipeOut, pipeErr io.Writer, pipeIn io.Reader) error {
	for _, command := range p.Commands {
		command.Async = true
		command.Run(ctx, pipeOut, pipeErr, pipeIn)
	}

	return nil
}

//============================================================================================

// SyncScripts defines a struct which is used to execute a giving set of
// shell script.
type SyncScripts struct {
	Scripts []ScriptProcess `json:"commands"`
}

// Exec executes the giving series of commands attached to the
// process.
func (p SyncScripts) Exec(ctx context.Context, pipeOut, pipeErr io.Writer, pipeIn io.Reader) error {
	for _, command := range p.Scripts {
		if err := command.Exec(ctx, pipeOut, pipeErr, pipeIn); err != nil {
			return err
		}
	}

	return nil
}

//============================================================================================

// ScriptProcess defines a shell script execution structure which attempts to copy
// given script into a local file path and attempts to execute content.
// Shell states the shell to be used for execution: /bin/sh, /bin/bash
type ScriptProcess struct {
	Shell  string        `json:"shell" toml:"shell"`
	Source string        `json:"source" toml:"source"`
	Level  CriticalLevel `json:"level" toml:"level"`
}

// Exec executes a copy of the giving script source in a temporary file which it then executes
// the contents.
func (c ScriptProcess) Exec(ctx context.Context, pipeOut, pipeErr io.Writer, pipeIn io.Reader) error {
	tmpFile, err := ioutil.TempFile("/tmp", "proc-shell")
	if err != nil {
		log.Emit(sinks.Error("Process : Error : Command : Begin Execution : %q : %+q", c.Shell, err))
		return err
	}

	if _, err := tmpFile.Write([]byte(c.Source)); err != nil {
		log.Emit(sinks.Error("Process : Error : Command : Begin Execution : %q : %+q", c.Shell, err))
		tmpFile.Close()
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		log.Emit(sinks.Error("Process : Error : Command : Begin Execution : %q : %+q", c.Shell, err))
		tmpFile.Close()
		return err
	}

	tmpFile.Close()

	defer os.Remove(tmpFile.Name())

	proc := exec.Command(c.Shell, tmpFile.Name())
	proc.Stdout = pipeOut
	proc.Stderr = pipeErr
	proc.Stdin = pipeIn

	if err := proc.Start(); err != nil {
		log.Emit(sinks.Error("Process : Error : Command : Begin Execution : %q : %+q", c.Shell, err))
		return err
	}

	go func() {
		<-ctx.Done()
		if proc.Process != nil {
			proc.Process.Kill()
		}
	}()

	if err := proc.Wait(); err != nil {
		log.Emit(sinks.Error("Process : Error : Command : Begin Execution : %q : %q", c.Shell, c.Source))

		if c.Level > Warning {
			return err
		}

		return nil
	}

	return nil
}
