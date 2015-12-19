package tasks

import (
	"github.com/influx6/faux/builders"
	"github.com/influx6/faux/fs"
	"github.com/influx6/faux/pub"
	"github.com/influx6/faux/pubro"
)

// WatchCommandConfig provides the configuration for the WatchCommand publisher.
type WatchCommandConfig struct {
	Dir      string
	Commands []string
}

// WatchCommand returns a new publisher using the WatchCommandConfig config.
func WatchCommand(wc WatchCommandConfig) pub.Publisher {
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

func init() {

	pubro.Register(pubro.Meta{
		Name: "tasks/watchCommand",
		Desc: `WatchCommand is a publisher which watches a directory and runs a
giving set of commands everytime a change occurs within that directory.`,
		Inject: WatchCommand,
	})

}
