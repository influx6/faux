package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/influx6/assets"
	"github.com/influx6/faux/builders"
	"github.com/influx6/faux/fs"
	"github.com/influx6/faux/pkg"
	"github.com/influx6/faux/pub"
	"github.com/influx6/faux/pubro"
)

// WatchCommandDirective provides the configuration for the WatchCommand publisher.
type WatchCommandDirective struct {
	Dir      string
	Commands []string
}

// WatchCommand returns a new publisher using the WatchCommandConfig config.
func WatchCommand(wc WatchCommandDirective) pub.Publisher {
	if wc.Dir == "" {
		panic("confi must contain a valid Dir")
	}

	//create the file watcher
	watcher := fs.Watch(fs.WatchConfig{
		Path: wc.Dir,
	})

	watcher.Bind(builders.CommandLauncher(wc.Commands), true)
	return watcher
}

// JSClientDirective provides a configuration for using the JSClient publisher.
type JSClientDirective struct {
	Package string
	Verbose bool
	Tags    []string
}

// JSClient returns a publisher for watching and executing a gopherjs build
// process for a go package.
func JSClient(js JSClientDirective) pub.Publisher {
	_, jsName := filepath.Split(js.Package)

	pkgs := append([]string{}, js.Package)
	packages, err := assets.GetAllPackageLists(pkgs)

	if err != nil {
		panic(err)
	}

	dir, err := assets.GetPackageDir(js.Package)
	if err != nil {
		panic(err)
	}

	jsbuild := builders.JSLauncher(builders.JSBuildConfig{
		Package:  js.Package,
		Folder:   dir,
		FileName: jsName,
		Tags:     js.Tags,
		Verbose:  js.Verbose,
	})

	watcher := fs.WatchSet(fs.WatchSetConfig{
		Path: packages,
		Validator: func(base string, info os.FileInfo) bool {
			if strings.Contains(base, ".git") {
				return false
			}

			if info != nil && info.IsDir() {
				return true
			}

			if filepath.Ext(base) != ".go" {
				return false
			}

			return true
		},
	})

	watcher.Bind(jsbuild, true)
	jsbuild.Send(true)

	return jsbuild
}

// GoBinaryDirective provides configuration for building a go binary from a project
// directory while watching using a filewatcher for changes.
type GoBinaryDirective struct {
	Package string
	Name    string
	OutDir  string
	BinArgs []string
}

// GoBinary returns a publisher for watching and building a go binary for a
// specific package name.
func GoBinary(gb GoBinaryDirective) pub.Publisher {
	pkgDir, err := pkg.GetPackageDir(gb.Package)
	if err != nil {
		panic(err)
	}

	// pkgDirBin defines where we wish to store the binary file.
	pkgDirBin := filepath.Join(pkgDir, "bin")

	// If directive provides its own output path then use that.
	if gb.OutDir != "" {
		pkgDirBin = gb.OutDir
	}

	pkgLists, err := pkg.GetPackageLists(gb.Package)
	if err != nil {
		panic(err)
	}

	pkg.GoDeps(pkgDir)

	var binName string

	if gb.Name != "" {
		binName = gb.Name
	} else {
		binName = filepath.Base(gb.Package)
	}

	buildbin := builders.BinaryBuildLauncher(builders.BinaryBuildConfig{
		Path:    pkgDir,
		Name:    binName,
		RunArgs: gb.BinArgs,
	})

	watcher := fs.WatchSet(fs.WatchSetConfig{
		Path: pkgLists,
		Validator: func(base string, info os.FileInfo) bool {
			if strings.Contains(base, ".git") {
				return false
			}

			if strings.Contains(base, pkgDirBin) {
				return false
			}

			if info != nil && info.IsDir() {
				return true
			}

			if filepath.Ext(base) != ".go" {
				return false
			}

			return true
		},
	})

	watcher.Bind(buildbin, true)

	return watcher
}

// MarkdownDirective provides a configuration for the GoFriday markdwon publisher.
// Its InDir and OutDir specify the directory to find markdwon files and the dir
// to store the output go template with the extension provided in Ext else
// defaults to .tmpl.
type MarkdownDirective struct {
	InDir    string
	OutDir   string
	Ext      string
	Sanitize bool
}

// Markdown2Templates returns a publisher which watches a specific directory of markdwon
// files and produces the equivalent rendered files as go template files,
// ensuring to keep the necessary dir structures appropriately.
func Markdown2Templates(md MarkdownDirective) pub.Publisher {
	if md.InDir == "" {
		panic(fmt.Sprintf("MarkdownDirective contains a empty In directory"))
	}

	if md.OutDir == "" {
		panic(fmt.Sprintf("MarkdownDirective contains a empty Out directory"))
	}

	gofriday, err := builders.GoFridayStream(builders.MarkStreamConfig{
		InputDir: md.InDir,
		SaveDir:  md.OutDir,
		Ext:      md.Ext,
		Sanitize: md.Sanitize,
	})

	if err != nil {
		panic(err)
	}

	//create the file watcher
	watcher := fs.Watch(fs.WatchConfig{
		Path: md.InDir,
	})

	// create the command runner set to run the args
	watcher.Bind(gofriday, true)

	return watcher
}

// StaticDirective provides configuration for using the gostatic publisher.
type StaticDirective struct {
	InDir       string
	OutDir      string
	PackageName string
	FileName    string
	Gzip        bool
	Production  bool
	Decompress  bool
	Ignore      *regexp.Regexp
}

// Static returns a publisher which watches a directory and builds a
func Static(sm StaticDirective) pub.Publisher {
	gostatic, err := builders.BundleAssets(&assets.BindFSConfig{
		InDir:           sm.InDir,
		OutDir:          sm.OutDir,
		Package:         sm.PackageName,
		File:            sm.FileName,
		Gzipped:         sm.Gzip,
		NoDecompression: sm.Decompress,
		Production:      sm.Production,
	})

	if err != nil {
		panic(err)
	}

	//bundle up the assets for the main time
	gostatic.Send(true)

	//create the file watcher
	watcher := fs.Watch(fs.WatchConfig{
		Path: sm.InDir,
		Validator: func(path string, info os.FileInfo) bool {
			if sm.Ignore != nil && sm.Ignore.MatchString(path) {
				return false
			}
			return true
		},
	})

	// create the command runner set to run the args
	watcher.Bind(gostatic, true)
	return watcher
}

func init() {
	pubro.Register(pubro.Meta{
		Name: "tasks/static",
		Desc: `Static watches a giving directory and produces go package which
		contains a virtual filesystem of the giving directory,It rebuilds the file
		as the directory changes`,
		Inject: Static,
	})

	pubro.Register(pubro.Meta{
		Name: "tasks/markdown2Templates",
		Desc: `Markdown2Template watches a giving directory and produces go
		template outputs of the markdown files ensuring to respect directory structure`,
		Inject: Markdown2Templates,
	})

	pubro.Register(pubro.Meta{
		Name: "tasks/goBinary",
		Desc: `Gobinary watches a package directory and builds its output
		accordingly into a bin dir in the package or at a custom path.`,
		Inject: GoBinary,
	})

	pubro.Register(pubro.Meta{
		Name: "tasks/watchCommand",
		Desc: `WatchCommand returns a publisher which watches a directory and runs a
		giving set of commands everytime a change occurs within that directory.`,
		Inject: WatchCommand,
	})

	pubro.Register(pubro.Meta{
		Name: "tasks/jsClient",
		Desc: `jsClient returns a publisher that compiles a gopherjs based package,
		and ensures to watch the directory for changes.`,
		Inject: JSClient,
	})
}
