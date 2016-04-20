package flux

import (
	"log"
	"sync"
	"testing"
)

func TestFlatReactor(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(1)

	fs := FlatReactive(func(r Reactor, err error, d interface{}) {
		ws.Done()
		if d != 1 {
			FatalFailed(t, "Incorrect received value: %d", d)
		}
		LogPassed(t, "Correct Value received: %d", d)
	})

	fs.Send(1)

	ws.Wait()
	fs.Close()
}

func TestFlatReact(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(1)

	ms := FlatIdentity()

	fs := ms.React(func(r Reactor, err error, d interface{}) {
		ws.Done()
		if d != 10 {
			FatalFailed(t, "Incorrect received value: %d", d)
		}
		LogPassed(t, "Correct Value received: %d", d)
	}, true)

	ms.Send(10)

	ws.Wait()
	fs.Close()
}

func TestFlatBind(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(1)

	ms := FlatIdentity()
	fs := FlatReactive(func(r Reactor, err error, d interface{}) {
		ws.Done()
		if d != 100 {
			FatalFailed(t, "Incorrect received value: %d", d)
		}
		LogPassed(t, "Correct Value received: %d", d)
	})

	ms.Bind(fs, true)

	ms.Send(100)

	ws.Wait()
	fs.Close()
}

func TestFlatDistribute(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(2)

	ms := FlatIdentity()

	fs := FlatReactive(func(m Reactor, err error, d interface{}) {
		log.Print("atfs")
		ws.Done()
		if d != 100 {
			FatalFailed(t, "Incorrect received value: %d", d)
		}
		LogPassed(t, "Correct Value received@fs: %d", d)
	})

	gs := FlatReactive(func(r Reactor, err error, d interface{}) {
		ws.Done()
		if d != 100 {
			FatalFailed(t, "Incorrect received value: %d", d)
		}
		LogPassed(t, "Correct Value received@gs: %d", d)
	})

	_ = DistributeSignals(ms, fs, gs)

	ms.Send(100)

	ws.Wait()
	fs.Close()
}
