package pub

import "sync/atomic"

// ChannelStream provides a simple struct for exposing outputs from Publisher to outside
type ChannelStream struct {
	Data   chan interface{}
	Error  chan error
	closed int64
}

// NewChannelStream returns a new channel stream instance with blocked channels, so ensure to fullfill the contract of removing the data you need only
func NewChannelStream() *ChannelStream {
	cs := ChannelStream{
		Data:  make(chan interface{}),
		Error: make(chan error),
	}
	return &cs
}

func (c *ChannelStream) error(err error) {
	if atomic.LoadInt64(&c.closed) <= 0 {
		go func() { c.Error <- err }()
		// c.Error <- err
	}
}

func (c *ChannelStream) data(d interface{}) {
	if atomic.LoadInt64(&c.closed) <= 0 {
		go func() { c.Data <- d }()
		// c.Data <- d
	}
}

// Close ends the capability to use the ChannelStream channels
func (c *ChannelStream) Close() {
	atomic.StoreInt64(&c.closed, 1)
}

// Listen binds into a Publisher and will pipe any response into its Data or Error channels, always use this to bind to Publishers, to ensure safety in code use i.e dont try to pipe into the channels your own way
func (c *ChannelStream) Listen(m Publisher) {
	if atomic.LoadInt64(&c.closed) > 0 {
		return
	}

	m.React(func(m Publisher, err error, data interface{}) {
		if err != nil {
			c.error(err)
			m.ReplyError(err)
			return
		}
		c.data(data)
		m.Reply(err)
	}, true)
}
