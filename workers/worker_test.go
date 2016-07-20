package workers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/workers"
)

//==============================================================================

var events eventlog

// logg provides a concrete implementation of a logger.
type eventlog struct{}

// Log logs all standard log reports.
func (l eventlog) Log(context interface{}, name string, message string, data ...interface{}) {
	fmt.Printf("Log: %s : %s : %s : %s\n", context, "DEV", name, fmt.Sprintf(message, data...))
}

// Error logs all error reports.
func (l eventlog) Error(context interface{}, name string, err error, message string, data ...interface{}) {
	fmt.Printf("Error: %s : %s : %s : %s : Error %s\n", context, "DEV", name, fmt.Sprintf(message, data...), err)
}

//==============================================================================

type writer struct{}

// Do writes returns a giving value as a byte slice else returns a non-nil error.
// Uses json.Marshal internally.
func (w writer) Do(ctx context.Context, err error, value interface{}) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(&value); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

//==============================================================================
type reader struct{}

// Do takes a expected value as string and returns the internal expected
// structure else returns a non-nil error.
func (r reader) Do(ctx context.Context, err error, value interface{}) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	src, ok := value.([]byte)
	if !ok {
		return nil, errors.New("Invalid Data. Expected []byte")
	}

	var d interface{}
	if err := json.NewDecoder(bytes.NewBuffer(src)).Decode(&d); err != nil {
		return nil, err
	}

	return d, nil
}

//==============================================================================

// basic structure we expect
type monster struct {
	Name string
}

func init() {
	tests.Init("")
}

// TestStreams validates the stream api operation and behaviour
func TestStreams(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	monsterName := "Willow"

	ws := workers.New(workers.Config{Log: events}, writer{})
	rs := ws.Next(workers.New(workers.Config{Log: events}, reader{}))
	rc, _ := workers.Receive(rs)
	erc, _ := workers.ReceiveError(rs)

	ws.Data(nil, &monster{Name: monsterName})

	var res interface{}
	var err error

	res = <-rc

	mo, ok := res.(map[string]interface{})
	if !ok {
		t.Fatalf("\t%s\tShould have received a monster map: %+s", tests.Failed, mo)
	}

	if mo["Name"] != monsterName {
		t.Fatalf("\t%s\tShould have received monster with Name[%s]", tests.Failed, monsterName)
	}
	t.Logf("\t%s\tShould have received monster with Name[%s]", tests.Success, monsterName)

	ex := errors.New("Bad Data")
	ws.Error(nil, ex)

	err = <-erc
	if err == nil {
		t.Fatalf("\t%s\tShould have received an error(%s) from stream: %s", tests.Failed, ex, err)
	}
	t.Logf("\t%s\tShould have received an error(%s) from stream", tests.Success, ex)

	ws.Shutdown()
	rs.Shutdown()

	select {
	case <-ws.CloseNotify():
		t.Logf("\t%s\tShould have closed first stream properly", tests.Success)
	case <-time.After(30 * time.Second):
		t.Errorf("\t%s\tShould have closed first stream properly", tests.Failed)
	}

	select {
	case <-rs.CloseNotify():
		t.Logf("\t%s\tShould have closed second stream properly", tests.Success)
	case <-time.After(30 * time.Second):
		t.Errorf("\t%s\tShould have closed second stream properly", tests.Failed)
	}
}

//==============================================================================

type dsync struct{}

func (dsync) Do(ctx context.Context, err error, value interface{}) (interface{}, error) {
	return value, err
}

type dasync struct{}

func (dasync) Do(ctx context.Context, err error, value interface{}) (interface{}, error) {
	time.Sleep(1 * time.Millisecond)
	return value, err
}

// BenchmarkNodes benches the performance of using the Node api.
func BenchmarkNodes(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	ds := workers.New(workers.Config{Min: 3, Max: 40, Log: events}, dsync{})
	das := workers.New(workers.Config{Min: 3, Max: 40, Log: events}, dasync{})

	ds.Next(das)

	for i := 0; i < b.N; i++ {
		ds.Data(nil, i)
	}
}

//==============================================================================

// BenchmarkStreams measures the performance of streamers using one worker.
func BenchmarkOneWorkerStreams(t *testing.B) {
	t.ResetTimer()
	t.ReportAllocs()

	ws := workers.New(workers.Config{Min: 3, Max: 40, Log: events}, writer{})
	defer ws.Shutdown()

	for i := 0; i < t.N; i++ {
		var wg sync.WaitGroup
		rc, rs := workers.Receive(ws)

		wg.Add(200)
		for j := 0; j < 200; j++ {
			go ws.Data(nil, &monster{Name: fmt.Sprintf("%d", j)})
		}

		go func() {
			for _ = range rc {
				wg.Done()
			}
		}()

		wg.Wait()
		rs.Shutdown()
	}
}

// BenchmarkNWorkerStreams measures the performance of streamers using the provided
// Nth value from the benchmark testing instance,
func BenchmarkNWorkerStreams(t *testing.B) {
	t.ResetTimer()
	t.ReportAllocs()

	ws := workers.New(workers.Config{Min: 20, Max: 1000, Log: events}, writer{})
	defer ws.Shutdown()

	for i := 0; i < t.N; i++ {
		var wg sync.WaitGroup
		rc, rs := workers.Receive(ws)

		wg.Add(200)
		for j := 0; j < 200; j++ {
			go ws.Data(nil, &monster{Name: fmt.Sprintf("%d", j)})
		}

		go func() {
			for _ = range rc {
				wg.Done()
			}
		}()

		wg.Wait()
		rs.Shutdown()
	}
}
