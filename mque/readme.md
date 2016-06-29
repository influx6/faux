# Mque
  Mque provides a argument aware pubsub queue.

## Concept
  The callback registered to a Mque queue defines the constraint upon which that
  callback should be called unless the callback expects no argument hence setting
  itself as callable for all data emissions.

## Usage

  ```go

			passed := make(chan int)
			failed := make(chan int)

			q := mque.New()

			q.Q(func() {
        // Will be called on every emission.
			})

			sub := q.Q(func(letter string) {
        // Will only be called when a string is emitted.
			})

			q.Q(func(item int) {
        // Will only be called when a int is emitted.
			})

			q.Run("letter")
			q.Run(60)

      sub.End()

  ```
