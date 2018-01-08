package seeker

// BufferedPeeker implements a custom buffer structure which
// allows peeks, reversal of index location of provided byte slice.
// It helps to minimize memory allocation.
type BufferedPeeker struct {
	l int
	d []byte
	c int
}

// NewBufferedPeeker returns new instance of BufferedPeeker.
func NewBufferedPeeker(d []byte) *BufferedPeeker {
	return &BufferedPeeker{
		d: d,
		l: len(d),
	}
}

// Reset sets giving buffers memory slice and sets appropriate
// settings.
func (b *BufferedPeeker) Reset(bm []byte) {
	b.d = bm
	b.l = len(bm)
	b.c = 0
}

// Length returns total length of slice.
func (b *BufferedPeeker) Length() int {
	return len(b.d)
}

// Area returns current available length from current index.
func (b *BufferedPeeker) Area() int {
	return len(b.d[b.c:])
}

// Reverse reverses previous index back a giving length
// It reverse any length back to 0 if length provided exceeds
// underline slice length.
func (b *BufferedPeeker) Reverse(n int) {
	back := b.c - n
	if back <= 0 && b.c == 0 {
		return
	}

	if back <= 0 {
		b.c = 0
		return
	}

	b.c = back
}

// Next returns bytes around giving range.
// If area is beyond slice length, then the rest of slice is returned.
func (b *BufferedPeeker) Next(n int) []byte {
	if b.c+n >= b.l {
		p := b.c
		b.c = b.l
		return b.d[p:]
	}

	area := b.d[b.c : b.c+n]
	b.c += n
	return area
}

// Peek returns the next bytes of giving n length. It
// will not move index beyond current bounds, but returns
// current area if within slice length.
// If area is beyond slice length, then the rest of slice is returned.
func (b *BufferedPeeker) Peek(n int) []byte {
	if b.c+n >= b.l {
		return b.d[b.c:]
	}

	return b.d[b.c : b.c+n]
}
