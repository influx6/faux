package flux

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestStackers tests the strict stacking structure for reactor as an alternative to its more library branch structure,basically it enforces a stacking down the tree,as far as their is a tail Reactor,all binding will be done with that
func TestStackers(t *testing.T) {

	mo := ReactStack(ReactIdentity())

	mo.React(func(r Reactor, err error, data interface{}) {
		if 2 != data {
			FatalFailed(t, "Data is inacurrate expected %d but got %d", 2, data)
		}

		LogPassed(t, "Data is corrected expected %d and got %d", 2, data)

		r.Reply(data.(int) * 40)
	}, true)

	if lo := mo.Length(); lo == 0 {
		FatalFailed(t, "Length is inacurrate expected 1 but got %d", lo)
	}

	mo.React(func(r Reactor, err error, data interface{}) {
		if 80 != data {
			FatalFailed(t, "Data is inacurrate expected %d but got %d", 80, data)
		}

		LogPassed(t, "Data is corrected expected %d and got %d", 80, data)

		r.Reply(data.(int) * 200)
	}, true)

	if lo := mo.Length(); lo < 2 {
		FatalFailed(t, "Length is inacurrate expected 2 but got %d", lo)
	}

	mi := mo.React(func(r Reactor, err error, data interface{}) {
		if 16000 != data {
			FatalFailed(t, "Data is inacurrate expected %d but got %d", 16000, data)
		}

		LogPassed(t, "Data is corrected expected %d and got %d", 16000, data)
	}, true)

	if lo := mo.Length(); lo < 3 {
		FatalFailed(t, "Length is inacurrate expected 3 but got %d", lo)
	}

	LogPassed(t, "Successfully stacked 3 Reactor. Length: %d", mo.Length())

	lmi, _ := mo.Last()

	if mi != lmi {
		FatalFailed(t, "Last Reactor is incorrect")
	}

	LogPassed(t, "Last Reactor matches the reference Reactor.")

	mo.Send(2)

	mo.Close()
}

// TestMousePosition provides a example test of a reactor that process mouse position (a slice of two values)
func TestMousePosition(t *testing.T) {
	var count int64
	var ws sync.WaitGroup

	ws.Add(2)

	pos := Reactive(func(r Reactor, err error, data interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}

		atomic.AddInt64(&count, 1)
		r.Reply(data)
	})

	go func() {
		//add about 3000 mouse events
		for i := 0; i < 2000; i++ {
			pos.Send([]int{i * 3, i + 1})
		}
		ws.Done()
	}()

	go func() {
		//add about 3000 mouse events
		for i := 0; i < 1000; i++ {
			pos.Send([]int{i * 3, i + 1})
		}
		ws.Done()
	}()

	ws.Wait()
	LogPassed(t, "Delivered all 3000 events")

	pos.Close()

	if atomic.LoadInt64(&count) != 3000 {
		FatalFailed(t, "Total processed values is not equal, expected %d but got %d", 3000, count)
	}

	LogPassed(t, "Total mouse data was processed with count %d", count)
}

// TestPartyOf2 test two independent reactors binding
func TestPartyOf2(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(2)

	dude := Reactive(func(r Reactor, err error, data interface{}) {
		r.Reply(data)
		ws.Done()
	})

	dudette := Reactive(func(r Reactor, err error, data interface{}) {
		r.Reply("4000")
		ws.Done()
	})

	dude.Bind(dudette, true)

	dude.Send("3000")

	// dudette.Close()
	ws.Wait()

	dude.Close()
}

func TestMerge(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(2)

	mo := ReactIdentity()
	mp := ReactIdentity()

	me := MergeReactors(mo, mp)

	me.React(func(v Reactor, err error, data interface{}) {
		ws.Done()
	}, true)

	mo.Send(1)
	mp.Send(2)

	ws.Wait()
	me.Close()
	//merge will not react to this
	mp.Send(4)

	<-me.CloseNotify()

	mo.Close()
	mp.Close()
}

func TestDistribute(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(300)

	master := ReactIdentity()

	slave := Reactive(func(r Reactor, err error, data interface{}) {
		if _, ok := data.(int); !ok {
			FatalFailed(t, "Data %+v is not int type", data)
		}
		ws.Done()
	})

	slave2 := Reactive(func(r Reactor, err error, data interface{}) {
		if _, ok := data.(int); !ok {
			FatalFailed(t, "Data %+v is not int type", data)
		}
		ws.Done()
	})

	DistributeSignals(master, slave, slave2)

	for i := 0; i < 150; i++ {
		master.Send(i)
	}
	LogPassed(t, "Successfully Sent 150 numbers")

	ws.Wait()

	LogPassed(t, "Successfully Processed 150 numbers")
	master.Close()
	slave.Close()
	slave2.Close()
}

func TestLift(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(2)

	master := ReactIdentity()

	slave := Reactive(func(r Reactor, err error, data interface{}) {
		ws.Done()
		if data != 40 {
			FatalFailed(t, "Incorrect value recieved,expect %d got %d", 40, data)
		}
		r.Reply(data.(int) * 20)
	})

	slave.React(func(r Reactor, err error, data interface{}) {
		ws.Done()
		if data != 800 {
			FatalFailed(t, "Incorrect value recieved,expect %d got %d", 800, data)
		}
	}, true)

	Lift(true, master, slave)

	master.Send(40)

	ws.Wait()

	slave.Close()
	// <-slave.CloseSignal()

	LogPassed(t, "Successfully Lifted numbers between 2 Reactors")
	master.Close()
}
