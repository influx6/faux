package builders

import (
	"sync"
	"testing"

	"github.com/influx6/flux"
	"github.com/influx6/reactors/fs"
)

func TestGoStream(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(2)

	mark, err := GoFridayStream(MarkStreamConfig{
		InputDir: "../fixtures/markdown",
		SaveDir:  "../fixtures/templates",
		Ext:      ".tmpl",
	})

	if err != nil {
		flux.FatalFailed(t, "Failed to build list streamer", err)
	}

	mark.React((func(root pub.Publisher, err error, data interface{}) {
		if err != nil {
			flux.FatalFailed(t, "Error  occured %+s", err)
		}
		ws.Done()
	}), true)

	mark.Send(true)
	ws.Wait()
	mark.Close()
}

func TestMarkdownStream(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(2)

	mark, err := MarkFridayStream(MarkStreamConfig{
		InputDir: "../fixtures/markdown",
		SaveDir:  "../fixtures/templates",
		Ext:      ".html",
	})

	if err != nil {
		flux.FatalFailed(t, "Failed to build list streamer", err)
	}

	mark.React((func(root pub.Publisher, err error, data interface{}) {
		if err != nil {
			flux.FatalFailed(t, "Error  occured %+s", err)
		}
		ws.Done()
	}), true)

	mark.Send(true)
	ws.Wait()
	mark.Close()
}

func TestMarkdown(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(1)

	mark := BlackFriday()

	mark.React((func(root pub.Publisher, err error, data interface{}) {
		if err != nil {
			flux.FatalFailed(t, "Error  occured %+s", err)
		}
		if rf, ok := data.(*RenderFile); ok {
			flux.LogPassed(t, "Markdown Rendered: %+s", rf.Data)
		} else {
			flux.FatalFailed(t, "Failure in receiving Correct Response: %+s", data)
		}
		ws.Done()
	}), true)

	mark.Send(&RenderFile{
		Path: "./slug.md",
		Data: []byte("# Slug Passive Aggressive -> SPA"),
	})

	ws.Wait()
	mark.Close()
}

func TestMarkFriday(t *testing.T) {
	ws := new(sync.WaitGroup)
	ws.Add(1)

	mark := MarkFriday(MarkConfig{
		SaveDir: "../fixtures/templates",
		Ext:     ".mdr",
	})

	mark.React((func(root pub.Publisher, err error, data interface{}) {
		if err != nil {
			flux.FatalFailed(t, "Error  occured %+s", err)
		}

		if rf, ok := data.(*fs.FileWrite); ok {
			flux.LogPassed(t, "Markdown Rendered: %+s -> %+s", rf.Data, rf)
		} else {
			flux.FatalFailed(t, "Failure in receiving Correct Response: %+s", data)
		}
		ws.Done()
	}), true)

	mark.Send("../fixtures/markdown/base.md")

	ws.Wait()
	mark.Close()
}
