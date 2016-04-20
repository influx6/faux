package flux

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
)

func BenchMarkBuffer(t *testing.B) {
	for i := 0; i < t.N; i++ {
		digits := NewBuffer()

		//add absurd amount of digits (3000) into the buffer
		for i := 0; i <= 7000; i++ {
			digits.Enqueue(i)
		}

		//remove 1000 total digits from buffer
		for i := 0; i <= 7000; i++ {
			digits.Dequeue()
		}

		//clear out buffer and ensure its empty
		digits.Clear()
	}
}

func TestBuffer(t *testing.T) {
	digits := NewBuffer()

	//add absurd amount of digits (3000) into the buffer
	for i := 0; i <= 3000; i++ {
		digits.Enqueue(i)
	}

	//length should be equal to 3000 in total
	if ln := digits.Length(); ln < 3000 {
		FatalFailed(t, "Total length of buffer incorrect. Execpted %d got %d", 3000, ln)
	}

	LogPassed(t, "Total length of buffer correct at %d", 3000)

	//remove 1000 total digits from buffer
	for i := 0; i <= 1000; i++ {
		digits.Dequeue()
	}

	if ln := digits.Length(); ln != 2000 {
		FatalFailed(t, "Total length of buffer incorrect. Execpted %d got %d", 2000, ln)
	}

	LogPassed(t, "Total length of buffer correct at %d", 2000)

	//clear out buffer and ensure its empty
	digits.Clear()

	// if ln := digits.Length(); ln != 0 {
	// 	FatalFailed(t, "Total length of buffer incorrect. Execpted %d got %d", 0, ln)
	// }

	LogPassed(t, "Buffer correct at %d", 0)
}

func TestQueue(t *testing.T) {
	dq := make(chan interface{})
	qo := NewQueue(dq)

	//queue around 5000 items into it
	for i := 0; i < 5000; i++ {
		qo.Enqueue(i)
	}

	//length should be equal to 5000 in total
	if ln := qo.Length(); ln > 5000 {
		FatalFailed(t, "Total length of queue incorrect. Execpted %d got %d", 5000, ln)
	}

	LogPassed(t, "Total length of queue correct at %d", 5000)

	//remove 1000 total digits from buffer
	for i := 0; i < 1000; i++ {
		if nxt := <-dq; nxt != i {
			FatalFailed(t, "Incorrect value received in dequeue, expected %d got %d", i, nxt)
		}
	}

	runtime.Gosched()

	//length should equal to around 4000
	if ln := qo.Length(); ln != 4000 {
		FatalFailed(t, "Total length of queue incorrect. Execpted %d got %d", 4000, ln)
	}

	LogPassed(t, "Total length of queue correct at %d", 4000)

	//removing remaining items in queue first
	for i := 0; i < 4000; i++ {
		<-dq
	}

	runtime.Gosched()

	if ln := qo.Length(); ln != 0 {
		FatalFailed(t, "Total length of queue incorrect. Execpted %d got %d", 0, ln)
	}

	//close out buffer and ensure its empty
	qo.Close()

	//ensure closed state of receive channel
	_, ok := <-dq

	if ok {
		FatalFailed(t, "Expected receive channel to be in closed state")
	}

	LogPassed(t, "Queue is empty and channel is closed")
}

func TestUnfinishedQueue(t *testing.T) {
	dq := make(chan interface{})
	qo := NewQueue(dq)

	//queue around 5000 items into it
	for i := 0; i < 5000; i++ {
		qo.Enqueue(i)
	}

	//length should be equal to 5000 in total
	if ln := qo.Length(); ln != 5000 {
		FatalFailed(t, "Total length of queue incorrect. Execpted %d got %d", 5000, ln)
	}

	LogPassed(t, "Total length of queue correct at %d", 3000)

	//remove 1000 total digits from buffer
	for i := 0; i < 1000; i++ {
		if nxt := <-dq; nxt != i {
			FatalFailed(t, "Incorrect value received in dequeue, expected %d got %d", i, nxt)
		}
	}

	runtime.Gosched()

	//length should equal to around 4000
	if ln := qo.Length(); ln > 4000 {
		FatalFailed(t, "Total length of queue incorrect. Execpted %d got %d", 4000, ln)
	}

	LogPassed(t, "Total length of queue correct at %d", 4000)

	//close queue even though we still have contents
	qo.Close()

	//ensure closed state of receive channel
	_, ok := <-dq

	if ok {
		FatalFailed(t, "Expected receive channel to be in closed state")
	}

	LogPassed(t, "Queue is empty and channel is closed")
}
func BenchmarkQueue(t *testing.B) {
	for i := 0; i < t.N; i++ {
		dqo := make(chan interface{})
		digits := NewQueue(dqo)

		//add absurd amount of digits (3000) into the buffer
		for i := 0; i <= 10000; i++ {
			digits.Enqueue(i)
		}

		//remove 1000 total digits from buffer
		for i := 0; i <= 10000; i++ {
			<-dqo
		}

		//clear out buffer and ensure its empty
		digits.Close()
	}
}

func TestPressureStream(t *testing.T) {
	ps := NewPressureStream()

	//add absurd amount of digits (3000) into the buffer
	for i := 0; i < 10000; i++ {
		ps.SendSignal(i)
	}

	//add absurd amount of digits (3000) into the buffer
	for i := 0; i < 10000; i++ {
		ps.SendError(fmt.Errorf("%d", i))
	}

	//removall the data
	for i := 0; i < 10000; i++ {
		if signal := <-ps.Signals; signal != i {
			FatalFailed(t, "Received in inaccurate signal, expected %d got %d", i, signal)
		}
	}

	LogPassed(t, "All signals were removed")

	//removall the data
	for i := 0; i < 10000; i++ {
		signal, ok := (<-ps.Errors).(error)
		if !ok {
			FatalFailed(t, "Signal type is incorrect expected error %+s", signal)
		}

		ind, err := strconv.Atoi(signal.Error())

		if err != nil {
			FatalFailed(t, "Error value not a stringified number %+s", signal.Error())
		}

		if ind != i {
			FatalFailed(t, "Received in inaccurate error, expected %d got %d", i, ind)
		}
	}

	LogPassed(t, "All signals were removed")

	//close and ensure we are closed
	ps.Close()

	_, ok := <-ps.Signals

	if ok {
		FatalFailed(t, "Signal channel is not closed")
	}

	_, ok = <-ps.Errors

	if ok {
		FatalFailed(t, "Signal channel is not closed")
	}

}

func BenchmarkPressureStream(t *testing.B) {
	for i := 0; i < t.N; i++ {
		ps := NewPressureStream()

		//add absurd amount of digits (3000) into the buffer
		for i := 0; i < 10000; i++ {
			ps.SendSignal(i)
		}

		//add absurd amount of digits (3000) into the buffer
		for i := 0; i < 10000; i++ {
			ps.SendError(fmt.Errorf("%d", i))
		}

		for i := 0; i < 10000; i++ {
			<-ps.Signals
		}

		for i := 0; i < 10000; i++ {
			<-ps.Errors
		}

		ps.Close()
	}
}
