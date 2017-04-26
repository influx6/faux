package mque_test

import (
	"testing"

	"github.com/influx6/faux/mque"
)

// TestQueue validates the behaviour of que api.
func TestQueue(t *testing.T) {
	t.Logf("Should be able to use a argument selective queue")
	{

		t.Logf("\tWhen giving a mque.Que and a string only allowed constraint")
		{

			passed := make(chan int)
			failed := make(chan int)

			q := mque.New()

			q.Q(func(letter string) {
				go func() { passed <- 1 }()
			})

			q.Q(func(item int) {
				go func() { failed <- 1 }()
			})

			q.Run("letter")

			select {
			case <-passed:
				t.Logf("Should have received a string")
			case <-failed:
				t.Errorf("Should have received a string")
			}

		}

		t.Logf("\tWhen giving a mque.Que and a int only allowed constraint")
		{
			passed := make(chan int)
			failed := make(chan int)

			q := mque.New()

			q.Q(func(letter string) {
				go func() { failed <- 1 }()
			})

			q.Q(func(item int) {
				go func() { passed <- 1 }()
			})

			q.Run(30)

			select {
			case <-failed:
				t.Fatalf("Should have received a integer")
			case <-passed:
				t.Logf("Should have received a integer")
			}
		}
	}
}

func BenchmarkQueue(b *testing.B) {
	b.ResetTimer()
	defer b.ReportAllocs()

	q := mque.New()

	q.Q(func(item int) {})
	q.Q(func(item int) {})

	for i := 0; i < b.N; i++ {
		q.Run(i)
	}
}

func BenchmarkQueueWithMultitypes(b *testing.B) {
	b.ResetTimer()
	defer b.ReportAllocs()

	q := mque.New()

	q.Q(func(item string) {})
	q.Q(func(item int) {})

	for i := 0; i < b.N; i++ {
		q.Run(i)
	}
}

// TestQueueEnd validates the behaviour of que subscriber End call.
func TestQueueEnd(t *testing.T) {

	t.Logf("Should be able to use a argument selective queue")
	{

		t.Logf("\tWhen needing to unsubscribe from a queue")
		{

			var count int
			var ended bool

			q := mque.New()

			q.Q(func(item int) {
				count++
			})

			sub := q.Q(func(item int) {
				count++
			}, func() {
				ended = true
			})

			q.Run(20)
			sub.End()

			q.Run(40)

			if count <= 2 || count >= 4 {
				t.Fatalf("Should have received only two events")
			}
			t.Logf("Should have received only two events")

			if !ended {
				t.Fatalf("Should have ended subscriber")
			}
			t.Logf("Should have ended subscriber")

		}
	}
}
