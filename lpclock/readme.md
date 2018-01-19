LPClock
-----------
LPClock implements a custom monotonic lamport clock.


## Usage

- Create a clock based on unix timestamps

```go
clock := lpclock.Unix("localhost")

// get new ever increasing monotonic time.
clock.Now()
```

- Create a clock based on ever increasing counter

```go
clock := lpclock.Lamport("localhost")

// get new ever increasing monotonic time.
clock.Now()
```
