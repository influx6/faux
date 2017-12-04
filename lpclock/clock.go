package lpclock

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	chars        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_abcdefghijklmnopqrstuvwxyz~"
	charsLen     = 28
	charsByte    = byte(28)
	defaultIDLen = 10
	tickOffset   = 30 * time.Second
)

// Clock exposes a giving monotonic lamport clocking structure which returns
// custom continouse ticks for giving origin and id.
type Clock struct {
	id     string
	origin string
	tickT  TickType
	offset time.Duration
	mu     sync.Mutex
	last   *UUID
}

// Lamport returns a new instance of Clock using the LAMPORTTICK has the time tick type.
func Lamport(origin string) *Clock {
	return New(LAMPORTTICK, origin)
}

// Unix returns a new instance of Clock using the UNIXTICK has the time tick type.
func Unix(origin string) *Clock {
	return New(UNIXTICK, origin)
}

// New returns a new instance of Clock.
func New(tickT TickType, origin string) *Clock {
	return NewClock(tickT, tickOffset, origin, string(generateLength(defaultIDLen)))
}

// NewClock returns new instance of Clock struct.
func NewClock(tickT TickType, timeOffset time.Duration, origin string, id string) *Clock {
	if strings.TrimSpace(origin) == "" {
		panic("origin can not be empty")
	}
	if strings.TrimSpace(id) == "" {
		panic("id can not be empty")
	}

	var clock Clock
	clock.id = id
	clock.tickT = tickT
	clock.origin = origin
	clock.offset = timeOffset
	return &clock
}

// Now returns new monotonic UUID which is consistently increasing.
func (c *Clock) Now() UUID {
	c.mu.Lock()
	defer c.mu.Unlock()

	var uuid UUID
	uuid.ID = c.id
	uuid.Type = c.tickT
	uuid.Origin = c.origin

	switch c.tickT {
	case LAMPORTTICK:
		if c.last != nil {
			uuid.Tick = c.last.Tick + 1
		} else {
			uuid.Tick = 1
		}
	case UNIXTICK:
		newTick := time.Now().UTC().Add(c.offset).Unix()
		if c.last != nil && newTick <= c.last.Tick {
			lastTime := time.Unix(c.last.Tick, 0)
			lastTime = lastTime.Add(c.offset).UTC()
			newTick = lastTime.Unix()
		}
		uuid.Tick = newTick
	}

	newUUID := uuid
	c.last = &newUUID

	return uuid
}

// generateLength returns random contents from the provided slice of base64
// encodable chars.
func generateLength(m int) []byte {
	content := make([]byte, m)
	rand.Read(content)
	for index, b := range content {
		content[index] = chars[b%charsByte]
	}
	return content
}

func init() {
	rand.Seed(time.Now().Unix())
}
