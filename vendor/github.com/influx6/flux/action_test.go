package flux

import (
	"errors"
	"testing"
)

func TestActDepend(t *testing.T) {
	ax := NewAction()
	ag := NewActDepend(ax, 3)

	_ = ag.OverrideBefore(1, func(b interface{}, next ActionInterface) {
		next.Fullfill(b)
	})

	fx := ag.Then(func(b interface{}, next ActionInterface) {
		next.Fullfill("through!")
	}).Then(func(b interface{}, next ActionInterface) {
		next.Fullfill(b)
	}).Then(func(b interface{}, next ActionInterface) {
		next.Fullfill("josh")
	})

	ax.Fullfill("Sounds!")

	<-fx.Sync(20)

}

func TestAction(t *testing.T) {
	ax := NewAction()

	_ = ax.Then(func(v interface{}, a ActionInterface) {
		_, ok := v.(string)

		if !ok {
			a.Fullfill(false)
			return
		}

		a.Fullfill(true)
	}).When(func(v interface{}, _ ActionInterface) {

		state, _ := v.(bool)

		if !state {
			t.Fatal("Fullfilled value is not a string")
		}

	})

	ax.Fullfill("Sounds!")
}

func TestSuccessActionStack(t *testing.T) {
	ax := NewActionStack()

	_ = ax.Done().When(func(v interface{}, _ ActionInterface) {
		if _, ok := v.(string); !ok {
			t.Fatal("Value received is not a string: ", v)
		}
	})

	_ = ax.Error().When(func(v interface{}, _ ActionInterface) {
		if _, ok := v.(string); !ok {
			t.Fatal("Value received is not a string: ", v)
		}
	})

	ax.Complete("Sounds!")
}

func TestFailedActionStack(t *testing.T) {
	ax := NewActionStack()

	_ = ax.Done().When(func(v interface{}, _ ActionInterface) {
		if _, ok := v.(string); !ok {
			t.Fatal("Value received is not a string: ", v)
		}
	})

	_ = ax.Error().When(func(v interface{}, _ ActionInterface) {
		if _, ok := v.(error); !ok {
			t.Fatal("Value received is not a error: ", v)
		}
	}).Then(func(v interface{}, b ActionInterface) {
		b.Fullfill(v)
	})

	ax.Complete(errors.New("1000"))
}
