package flux

import (
	"sync"
	"testing"
	"time"
)

func TestResetTimer(t *testing.T) {
	ws := new(sync.WaitGroup)
	rs := NewResetTimer(func() {
	}, func() {
		ws.Done()
	}, time.Duration(600)*time.Millisecond, true, true)

	ws.Add(2)

	go func() {
		time.Sleep(time.Duration(650) * time.Millisecond)
		rs.Add()
	}()

	ws.Wait()
	rs.Close()
}

func TestOneTimeWaiter(t *testing.T) {
	ms := time.Duration(3) * time.Second

	w := NewTimeWait(0, ms)
	ws := new(sync.WaitGroup)

	w.Add()
	ws.Add(1)

	w.Then().When(func(v interface{}, _ ActionInterface) {
		_, ok := v.(int)

		if !ok {
			t.Fatal("Waiter completed with non-int value:", v)
		}

		ws.Done()
	})

	if w == nil {
		t.Fatal("Unable to create waiter")
	}

	if w.Count() < 0 {
		t.Fatal("Waiter completed before use")
	}

	if w.Count() < -1 {
		t.Fatal("Waiter just did a bad logic and went below -1")
	}

	// t.Log("Calling Done()")
	w.Done()

	ws.Wait()
}

func TestTimeWaiter(t *testing.T) {
	ms := time.Duration(1) * time.Second

	w := NewTimeWait(0, ms)
	ws := new(sync.WaitGroup)

	ws.Add(1)

	w.Then().When(func(v interface{}, _ ActionInterface) {
		_, ok := v.(int)

		if !ok {
			t.Fatal("Waiter completed with non-int value:", v)
		}

		ws.Done()
	})

	if w == nil {
		t.Fatal("Unable to create waiter")
	}

	if w.Count() < 0 {
		t.Fatal("Waiter completed before use")
	}

	w.Add()
	w.Add()

	if w.Count() < -1 {
		t.Fatal("Waiter just did a bad logic and went below -1")
	}

	// w.Done()

	ws.Wait()
}

func TestWaiter(t *testing.T) {
	w := NewWait()

	w.Then().When(func(v interface{}, _ ActionInterface) {
		_, ok := v.(int)

		if !ok {
			t.Fatal("Waiter completed with non-int value:", v)
		}

	})

	if w == nil {
		t.Fatal("Unable to create waiter")
	}

	if w.Count() < 0 {
		t.Fatal("Waiter completed before use")
	}

	w.Add()

	w.Done()

	cu := w.Count()
	if w.Count() > 1 {
		t.Fatalf("Waiter count (%v) is still greater (%v) despite call to Done()", w.Count(), cu)
	}

	w.Done()
	w.Done()

	if w.Count() < -1 {
		t.Fatal("Waiter just did a bad logic and went below -1")
	}
}
