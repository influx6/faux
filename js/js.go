package js

import (
	"bytes"
	"errors"

	build "github.com/gopherjs/gopherjs/build"
	"github.com/gopherjs/gopherjs/compiler"
	"github.com/neelance/sourcemap"
)

// ErrNotMain is returned when we find no .go file with 'main' package
var ErrNotMain = errors.New("Package contains no 'main' go package file")

// Session represents a basic build.Session with its option
type Session struct {
	//Dir to use for the virtual files
	dir     string
	Session *build.Session
	Option  *build.Options
}

// New returns a new session for build js files
func New(tags []string, verbose, watch bool) *Session {

	options := &build.Options{
		Verbose:       verbose,
		Watch:         watch,
		CreateMapFile: true,
		Minify:        true,
		BuildTags:     tags,
	}

	session := build.NewSession(options)

	return &Session{
		Session: session,
		Option:  options,
	}
}

// BuildPkg uses the session, to build a package file with the given output name and returns two virtual files containing the js and js.map respectively, or an error
func (j *Session) BuildPkg(pkg, name string) (*bytes.Buffer, *bytes.Buffer, error) {
	var js, jsmap *bytes.Buffer = bytes.NewBuffer(nil), bytes.NewBuffer(nil)

	if err := BuildJS(j, pkg, name, js, jsmap); err != nil {
		return nil, nil, err
	}

	return js, jsmap, nil
}

// BuildDir uses the session, to build a particular dir contain files and using the specified package name and output name returns two virtual files containing the js and js.map respectively, or an error
func (j *Session) BuildDir(dir, importpath, name string) (*bytes.Buffer, *bytes.Buffer, error) {
	var js, jsmap *bytes.Buffer = bytes.NewBuffer(nil), bytes.NewBuffer(nil)

	if err := BuildJSDir(j, dir, importpath, name, js, jsmap); err != nil {
		return nil, nil, err
	}

	return js, jsmap, nil
}

// BuildJSDir builds the js file and returns the content.
// goPkgPath must be a package path eg. github.com/influx6/haiku-examples/app
func BuildJSDir(jsession *Session, dir, importpath, name string, js, jsmap *bytes.Buffer) error {

	session, options := jsession.Session, jsession.Option

	buildpkg, err := build.NewBuildContext(session.InstallSuffix(), options.BuildTags).ImportDir(dir, 0)

	if err != nil {
		return err
	}

	pkg := &build.PackageData{Package: buildpkg}
	pkg.ImportPath = importpath

	//build the package using the sessios
	if err = session.BuildPackage(pkg); err != nil {
		return err
	}

	//build up the source map also
	smfilter := &compiler.SourceMapFilter{Writer: js}
	smsrc := &sourcemap.Map{File: name + ".js"}
	smfilter.MappingCallback = build.NewMappingCallback(smsrc, options.GOROOT, options.GOPATH)
	deps, err := compiler.ImportDependencies(pkg.Archive, session.ImportContext.Import)

	if err != nil {
		return err
	}

	err = compiler.WriteProgramCode(deps, smfilter)

	smsrc.WriteTo(jsmap)
	js.WriteString("//# sourceMappingURL=" + name + ".map.js\n")

	return nil
}

// BuildJS builds the js file and returns the content.
// goPkgPath must be a package path eg. github.com/influx6/haiku-examples/app
func BuildJS(jsession *Session, goPkgPath, name string, js, jsmap *bytes.Buffer) error {

	session, options := jsession.Session, jsession.Option

	//get the build path
	buildpkg, err := build.Import(goPkgPath, 0, session.InstallSuffix(), options.BuildTags)

	if err != nil {
		return err
	}

	if buildpkg.Name != "main" {
		return ErrNotMain
	}

	//build the package data for building
	// pkg := &build.PackageData{Package: buildpkg}

	//build the package using the sessios
	if err = session.BuildPackage(buildpkg); err != nil {
		return err
	}

	//build up the source map also
	smfilter := &compiler.SourceMapFilter{Writer: js}
	smsrc := &sourcemap.Map{File: name + ".js"}
	smfilter.MappingCallback = build.NewMappingCallback(smsrc, options.GOROOT, options.GOPATH)
	deps, err := compiler.ImportDependencies(buildpkg.Archive, session.ImportContext.Import)

	if err != nil {
		return err
	}

	err = compiler.WriteProgramCode(deps, smfilter)

	smsrc.WriteTo(jsmap)
	js.WriteString("//# sourceMappingURL=" + name + ".map.js\n")

	return nil
}
