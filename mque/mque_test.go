package mque_test

import (
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/mque"
)

func init() {
	tests.Init("")
}

// TestQueue validates the behaviour of que api.
func TestQueue(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

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
				t.Logf("\t%s\tShould have received a string", tests.Success)
			case <-failed:
				t.Errorf("\t%s\tShould have received a string", tests.Failed)
			}

		}

		t.Logf("\tWhen giving a mque.Que and a int only allowed constraint")
		{
			passed := make(chan int)
			failed := make(chan int)

			q := mque.New()

			q.Q(func(letter string) {
				// fmt.Printf("letter: %s\n", letter)
				go func() { failed <- 1 }()
			})

			q.Q(func(item int) {
				// fmt.Printf("digit: %d\n", item)
				go func() { passed <- 1 }()
			})

			q.Run(30)

			select {
			case <-failed:
				t.Logf("\t%s\tShould have received a integer", tests.Success)
			case <-passed:
				t.Errorf("\t%s\tShould have received a integer", tests.Failed)
			}

		}
	}
}

// TestQueueEnd validates the behaviour of que subscriber End call.
func TestQueueEnd(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Logf("Should be able to use a argument selective queue")
	{

		t.Logf("\tWhen needing to unsubscribe from a queue")
		{

			var count int

			q := mque.New()

			q.Q(func(item *int) {
				count++
			})

			sub := q.Q(func(item int) {
				count++
			})

			q.Run(20)
			sub.End()

			var numb = 40

			q.Run(&numb)

			if count < 2 {
				t.Fatalf("\t%s\tShould have received more than two events", tests.Failed)
			}
			t.Logf("\t%s\tShould have received more than two events", tests.Success)

		}
	}
}
