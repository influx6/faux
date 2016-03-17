package sumex_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/sumex"
)

type writer struct{}

// Do writes returns a giving value as a byte slice else returns a non-nil error.
// Uses json.Marshal internally.
func (w writer) Do(value interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return json, nil
}

type reader struct{}

// Do takes a expected value as string and returns the internal expected
// structure else returns a non-nil error.
func (r reader) Do(value interface{}, err error) (interface{}, error) {
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

	ws := sumex.New(3, nil, writer{})
	rs := ws.Stream(sumex.New(3, nil, reader{}))
	rc, _ := sumex.Receive(rs)
	erc, _ := sumex.ReceiveError(rs)

	defer ws.Shutdown()
	defer rs.Shutdown()

	ws.Inject(&monster{Name: monsterName})

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
	ws.InjectError(ex)

	err, ok = <-erc
	if err == nil {
		t.Fatalf("\t%s\tShould have received an error(%s) from stream: %s", tests.Failed, ex, err)
	}
	t.Logf("\t%s\tShould have received an error(%s) from stream", tests.Success, ex)

}

// BenchmarkStreams measures the performance of streamers using one worker.
func BenchmarkOneWorkerStreams(t *testing.B) {
	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		ws := sumex.New(1, nil, writer{})
		rc, _ := sumex.Receive(ws)

		defer ws.Shutdown()

		for j := 0; j < 10; j++ {
			go ws.Inject(&monster{Name: fmt.Sprintf("%d", j)})
		}

		<-rc
	}
}

// BenchmarkNWorkerStreams measures the performance of streamers using the provided
// Nth value from the benchmark testing instance,
func BenchmarkNWorkerStreams(t *testing.B) {
	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		ws := sumex.New(40, nil, writer{})
		rc, _ := sumex.Receive(ws)

		defer ws.Shutdown()

		for j := 0; j < 10; j++ {
			go ws.Inject(&monster{Name: fmt.Sprintf("%d", j)})
		}

		<-rc
	}
}
