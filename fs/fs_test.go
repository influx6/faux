package fs

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/influx6/flux"
)

func TestReader(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(1)

	read := FileReader()

	read.React(func(r pub.Publisher, err error, ev interface{}) {
		ws.Done()
	}, true)

	read.Send("./fs.go")
	ws.Wait()
	read.Close()
}

func TestWriter(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(1)

	read := FileWriter(nil)

	read.React(func(r pub.Publisher, err error, ev interface{}) {
		ws.Done()
	}, true)

	read.Send(&FileWrite{Path: "../fixtures/markdown/book.md", Data: []byte("# Book\n Write in love")})

	ws.Wait()
	read.Close()
}

func TestWatch(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(2)

	watcher := Watch(WatchConfig{
		Path: "../fixtures",
	})

	watcher.React(func(r pub.Publisher, err error, ev interface{}) {
		ws.Done()
	}, true)

	go func() {
		if md, err := os.Create("../fixtures/read.md"); err == nil {
			md.Close()
		}
		<-time.After(3 * time.Second)
		os.Remove("../fixtures/read.md")
		ws.Done()
	}()

	ws.Wait()
	watcher.Close()
}

func TestWatchSet(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(2)

	watcher := WatchSet(WatchSetConfig{
		Path: []string{"../fixtures"},
	})

	watcher.React(func(r pub.Publisher, err error, ev interface{}) {
		// log.Printf("err: %s %s", err, ev)
		ws.Done()
	}, true)

	go func() {
		if md, err := os.Create("../fixtures/dust.md"); err == nil {
			md.Close()
		}
		<-time.After(2 * time.Second)
		os.Remove("../fixtures/dust.md")
		ws.Done()
	}()

	ws.Wait()
	watcher.Close()
}

func TestListStreaming(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(9)

	lists, err := StreamListings(ListingConfig{
		Path:        "../builders",
		DirAlso:     true,
		UseRelative: true,
	})

	if err != nil {
		flux.FatalFailed(t, "Failed to build list streamer", err)
	}

	lists.React(func(r pub.Publisher, err error, ev interface{}) {
		// log.Printf("File: %s", ev)
		ws.Done()
	}, true)

	lists.Send(true)
	ws.Wait()
	lists.Close()
}
