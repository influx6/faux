// Package pkg provides a package interfacing API using the go build
// libraries to provide high level functions that simplify and running
// building processes.
package pkg

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var multispaces = regexp.MustCompile(`\s+`)

// GoDeps calls go get for specific package
func GoDeps(targetdir string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("godeps.Error: %+s", err)
		}
	}()

	cmdline := []string{"go", "get"}

	cmdline = append(cmdline, targetdir)

	//setup the executor and use a shard buffer
	cmd := exec.Command("go", cmdline[1:]...)
	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	cmd.Stderr = buf

	err := cmd.Run()

	if buf.Len() > 0 {
		return fmt.Errorf("go get failed: %s: %s", buf.String(), err.Error())
	}

	return nil
}

// GoRun runs the runs a command
func GoRun(cmd string) string {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("gorun.Error: %+s", err)
		}
	}()
	var cmdline []string
	com := strings.Split(cmd, " ")

	if len(com) < 0 {
		return ""
	}

	if len(com) == 1 {
		cmdline = append(cmdline, com...)
	} else {
		cmdline = append(cmdline, com[0])
		cmdline = append(cmdline, com[1:]...)
	}

	//setup the executor and use a shard buffer
	cmdo := exec.Command(cmdline[0], cmdline[1:]...)
	buf := bytes.NewBuffer([]byte{})
	cmdo.Stdout = buf
	cmdo.Stderr = buf

	_ = cmdo.Run()

	return buf.String()
}

// GobuildArgs runs the build process against a directory, using the giving
// arguments. Returns a non-nil error if it fails.
func GobuildArgs(args []string) error {
	if len(args) <= 0 {
		return nil
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("gobuild.Error: %+s", err)
		}
	}()

	cmdline := []string{"go", "build"}

	// target := filepath.Join(dir, name)
	cmdline = append(cmdline, args...)

	//setup the executor and use a shard buffer
	cmd := exec.Command("go", cmdline[1:]...)
	buf := bytes.NewBuffer([]byte{})

	msg, err := cmd.CombinedOutput()

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("go.build failed: %s: %s -> Msg: %s", buf.String(), err.Error(), msg)
	}

	return nil
}

// Gobuild runs the build process against a directory, using a giving name for the
// build file. Returns a non-nil error if it fails.
func Gobuild(dir, name string, args []string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("gobuild.Error: %+s", err)
		}
	}()

	cmdline := []string{"go", "build"}

	if runtime.GOOS == "windows" {
		name = fmt.Sprintf("%s.exe", name)
	}

	target := filepath.Join(dir, name)
	cmdline = append(cmdline, args...)
	cmdline = append(cmdline, "-o", target)

	//setup the executor and use a shard buffer
	cmd := exec.Command("go", cmdline[1:]...)
	buf := bytes.NewBuffer([]byte{})

	msg, err := cmd.CombinedOutput()

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("go.build failed: %s: %s -> Msg: %s", buf.String(), err.Error(), msg)
	}

	return nil
}

// RunCMD runs the a set of commands from a list while skipping any one-length command, panics if it gets an empty lists
func RunCMD(cmds []string, done func()) chan bool {
	if len(cmds) < 0 {
		panic("commands list cant be empty")
	}

	var relunch = make(chan bool)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("cmdRun.Error: %+s", err)
			}
		}()

	cmdloop:
		for {
			select {
			case do, ok := <-relunch:

				if !ok {
					break cmdloop
				}

				if !do {
					continue
				}

				for _, cox := range cmds {

					cmd := strings.Split(cox, " ")

					if len(cmd) <= 1 {
						continue
					}

					cmdo := exec.Command(cmd[0], cmd[1:]...)
					cmdo.Stdout = os.Stdout
					cmdo.Stderr = os.Stderr

					if err := cmdo.Start(); err != nil {
						fmt.Printf("---> Error executing command: %s -> %s\n", cmd, err)
					}
				}

				if done != nil {
					done()
				}
			}
		}

	}()
	return relunch
}

// RunGo runs the generated binary file with the arguments expected
func RunGo(gofile string, args []string, done, stopped func()) chan bool {
	var relunch = make(chan bool)

	// if runtime.GOOS == "windows" {
	gofile = filepath.Clean(gofile)
	// }

	go func() {

		// var cmdline = fmt.Sprintf("go run %s", gofile)
		cmdargs := append([]string{"run", gofile}, args...)
		// cmdline = strings.Joinappend([]string{}, "go run", gofile)

		var proc *os.Process

		for dosig := range relunch {
			if proc != nil {
				var err error

				if runtime.GOOS == "windows" {
					err = proc.Kill()
				} else {
					err = proc.Signal(os.Interrupt)
				}

				if err != nil {
					fmt.Printf("---> Error in Sending Kill Signal %s\n", err)
					proc.Kill()
				}
				proc.Wait()
				proc = nil
			}

			if !dosig {
				continue
			}

			cmd := exec.Command("go", cmdargs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				fmt.Printf("---> Error starting process: %s\n", err)
			}

			proc = cmd.Process
			if done != nil {
				done()
			}
		}

		if stopped != nil {
			stopped()
		}
	}()
	return relunch
}

// RunBin runs the generated binary file with the arguments expected
func RunBin(binfile string, args []string, done, stopped func()) chan bool {
	var relunch = make(chan bool)
	go func() {
		// binfile := fmt.Sprintf("%s/%s", bindir, bin)
		// cmdline := append([]string{bin}, args...)
		var proc *os.Process

		for dosig := range relunch {
			if proc != nil {
				var err error

				if runtime.GOOS == "windows" {
					err = proc.Kill()
				} else {
					err = proc.Signal(os.Interrupt)
				}

				if err != nil {
					fmt.Printf("---> Error in Sending Kill Signal: %s\n", err)
					proc.Kill()
				}
				proc.Wait()
				proc = nil
			}

			if !dosig {
				continue
			}

			cmd := exec.Command(binfile, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				fmt.Printf("---> Error starting process: %s -> %s\n", binfile, err)
			}

			proc = cmd.Process
			if done != nil {
				done()
			}
		}

		if stopped != nil {
			stopped()
		}
	}()
	return relunch
}

//SanitizeDuplicates cleans out all duplicates
func SanitizeDuplicates(b []string) []string {
	sz := len(b) - 1
	for i := 0; i < sz; i++ {
		for j := i + 1; j <= sz; j++ {
			if (b)[i] == ((b)[j]) {
				(b)[j] = (b)[sz]
				(b) = (b)[0:sz]
				sz--
				j--
			}
		}
	}
	return b
}

// GetPackageDir returns the directory of a package path from the go src dir.
func GetPackageDir(pkgname string) (string, error) {
	pkg, err := build.Import(pkgname, "", 0)

	if err != nil {
		return "", err
	}

	return pkg.Dir, nil
}

// GetPackageLists retrieves a packages  directory and those of its dependencies
func GetPackageLists(pkgname string) ([]string, error) {
	var paths []string
	var err error

	if paths, err = getPackageLists(pkgname, paths); err != nil {
		return nil, err
	}

	return SanitizeDuplicates(paths), nil
}

// GetAllPackageLists retrieves a set of packages directory and those of its dependencies
func GetAllPackageLists(pkgnames []string) ([]string, error) {
	var packages []string
	var err error

	for _, pkg := range pkgnames {
		if packages, err = getPackageLists(pkg, packages); err != nil {
			return nil, err
		}
	}

	// log.Printf("Packages: %s", packages)
	return SanitizeDuplicates(packages), nil
}

// getPackageLists returns the lists of internal package imports used within
// a giving package.
func getPackageLists(pkgname string, paths []string) ([]string, error) {
	pkg, err := build.Import(pkgname, "", 0)

	if err != nil {
		return nil, err
	}

	if pkg.Goroot {
		return paths, nil
	}

	paths = append(paths, pkg.Dir)

	for _, imp := range pkg.Imports {
		if p, err := getPackageLists(imp, paths); err == nil {
			paths = p
		} else {
			return nil, err
		}
	}

	return paths, nil
}
