package pub

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// succeedMark is the Unicode codepoint for a check mark.
const succeedMark = "\u2713"

// failedMark is the Unicode codepoint for an X mark.
const failedMark = "\u2717"

// TestStackers tests the strict stacking structure for Publisher as an alternative to its more library branch structure,basically it enforces a stacking down the tree,as far as their is a tail Publisher,all binding will be done with that
func TestStackers(t *testing.T) {

	mo := PublisherStack()

	mo.React(func(r Publisher, err error, data interface{}) {
		if 2 != data {
			fatalFailed(t, "Data is inacurrate expected %d but got %d", 2, data)
		}

		logPassed(t, "Data is corrected expected %d and got %d", 2, data)

		r.Reply(data.(int) * 40)
	}, true)

	mo.React(func(r Publisher, err error, data interface{}) {
		if 80 != data {
			fatalFailed(t, "Data is inacurrate expected %d but got %d", 80, data)
		}

		logPassed(t, "Data is corrected expected %d and got %d", 80, data)

		r.Reply(data.(int) * 200)
	}, true)

	mi := mo.React(func(r Publisher, err error, data interface{}) {
		if 16000 != data {
			fatalFailed(t, "Data is inacurrate expected %d but got %d", 16000, data)
		}

		logPassed(t, "Data is corrected expected %d and got %d", 16000, data)
	}, true)

	mo.Send(2)

	mi.Close()
	mo.Close()
}

// TestMousePosition provides a example test of a Publisher that process mouse position (a slice of two values)
func TestMousePosition(t *testing.T) {
	var count int64
	var ws sync.WaitGroup

	ws.Add(2)

	pos := Pubb(func(r Publisher, err error, data interface{}) {
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
	logPassed(t, "Delivered all 3000 events")

	pos.Close()

	if atomic.LoadInt64(&count) != 3000 {
		fatalFailed(t, "Total processed values is not equal, expected %d but got %d", 3000, count)
	}

	logPassed(t, "Total mouse data was processed with count %d", count)
}

// TestPartyOf2 test two independent Publishers binding
func TestPartyOf2(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(2)

	dude := Pubb(func(r Publisher, err error, data interface{}) {
		r.Reply(data)
		ws.Done()
	})

	dudette := Pubb(func(r Publisher, err error, data interface{}) {
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

	mo := Identity()
	mp := Identity()

	me := Merge(mo, mp)

	me.React(func(v Publisher, err error, data interface{}) {
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

	master := Identity()

	slave := Pubb(func(r Publisher, err error, data interface{}) {
		if _, ok := data.(int); !ok {
			fatalFailed(t, "Data %+v is not int type", data)
		}
		ws.Done()
	})

	slave2 := Pubb(func(r Publisher, err error, data interface{}) {
		if _, ok := data.(int); !ok {
			fatalFailed(t, "Data %+v is not int type", data)
		}
		ws.Done()
	})

	Distribute(master, slave, slave2)

	for i := 0; i < 150; i++ {
		master.Send(i)
	}
	logPassed(t, "Successfully Sent 150 numbers")

	ws.Wait()

	logPassed(t, "Successfully Processed 150 numbers")
	master.Close()
	slave.Close()
	slave2.Close()
}

func TestLift(t *testing.T) {
	var ws sync.WaitGroup
	ws.Add(2)

	master := Identity()

	slave := Pubb(func(r Publisher, err error, data interface{}) {
		ws.Done()
		if data != 40 {
			fatalFailed(t, "Incorrect value recieved,expect %d got %d", 40, data)
		}
		r.Reply(data.(int) * 20)
	})

	slave.React(func(r Publisher, err error, data interface{}) {
		ws.Done()
		if data != 800 {
			fatalFailed(t, "Incorrect value recieved,expect %d got %d", 800, data)
		}
	}, true)

	Lift(true, master, slave)

	master.Send(40)

	ws.Wait()

	slave.Close()
	// <-slave.CloseSignal()

	logPassed(t, "Successfully Lifted numbers between 2 Publishers")
	master.Close()
}

func logPassed(t *testing.T, msg string, data ...interface{}) {
	t.Logf("%s %s", fmt.Sprintf(msg, data...), succeedMark)
}

func fatalFailed(t *testing.T, msg string, data ...interface{}) {
	t.Errorf("%s %s", fmt.Sprintf(msg, data...), failedMark)
}
