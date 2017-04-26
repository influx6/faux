package subscriptions

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"sync"
)

var (
	sbl                    = "."
	slash                  = byte('/')
	slashSlice             = []byte{slash}
	colon                  = byte(':')
	colonSlice             = []byte{colon}
	emptySlice             = []byte("")
	sublist                = byte('.')
	sublistSlice           = []byte{sublist}
	edges                  = byte('^')
	edgesSlice             = []byte{edges}
	contains               = byte('*')
	containsSlice          = []byte{contains}
	containsArraySlice     = [][]byte{containsSlice}
	startBracket           = byte('[')
	startBracketSlice      = []byte{startBracket}
	endBracket             = byte(']')
	endBracketSlice        = []byte{endBracket}
	startCurlyBracket      = byte('{')
	startCurlyBracketSlice = []byte{startCurlyBracket}
	endCurlyBracket        = byte('}')
	endCurlyBracketSlice   = []byte{endCurlyBracket}
)

// Trace defines an interface which receives data trace data logs.
type Trace interface {
	Trace(context interface{}, msg []byte)
}

// Message defines the structure returned to a subscriber once a publish matching
// its criteria is found.
type Message struct {
	Msid    []byte
	Topic   []byte
	Match   []byte
	Params  map[string]string
	Payload interface{}
	Source  interface{}
}

// Subscriber defines an interface for routes to be fired upon when matched.
type Subscriber interface {
	Fire(context interface{}, sm *Message) error
}

// Subscription defines a struct for storing subscriptions.
type Subscription struct {
	root  *level
	cache *subCache
}

// New returns a new Subscription which can be used to route to specific
// subscribers. The variadic argument is used for elegance and to allow zero,
// zero value tracers.
func New(t ...Trace) *Subscription {
	var sub Subscription

	var cache subCache
	sub.cache = &cache

	var tr Trace

	if len(t) > 0 {
		tr = t[0]
	}

	sub.root = newLevel(tr)

	return &sub
}

// RoutesFor returns a giving route for a giving element.
func (s *Subscription) RoutesFor(sub Subscriber) ([][]byte, error) {
	sb, found := s.cache.Find(sub)
	if !found {
		return nil, errors.New("Not found")
	}

	var ju [][]byte

	for _, ns := range sb.ns {
		ju = append(ju, ns)
	}

	return ju, nil
}

// Routes returns the lists of routes which the subscription holds.
func (s *Subscription) Routes() [][]byte {
	return s.root.Routes()
}

// Handle calls the giving path slice, if found and applies the payload else
// returns an error.
func (s *Subscription) Handle(context interface{}, path []byte, payload interface{}, source interface{}) {
	msg := Message{
		Topic:   path,
		Payload: payload,
		Params:  map[string]string{},
		Source:  source,
	}

	s.root.Resolve(context, path, &msg)
}

// Register adds the new giving path slice into the subscription for routing.
func (s *Subscription) Register(path []byte, sub Subscriber) error {
	if err := s.root.Add(path, sub); err != nil {
		return err
	}

	s.cache.Add(sub, path)
	return nil
}

// Unregister removes the existing giving path slice into the subscription for routing.
func (s *Subscription) Unregister(path []byte, sub Subscriber) error {
	if err := s.root.Remove(path, sub); err != nil {
		return err
	}

	s.cache.Remove(sub, path)
	return nil
}

//==============================================================================

type subCache struct {
	rw    sync.RWMutex
	cache []subyList
}

type subyList struct {
	sub Subscriber
	ns  [][]byte
}

func (s *subCache) Remove(sub Subscriber, path []byte) {
	s.rw.RLock()
	{
		for _, target := range s.cache {
			if target.sub != sub {
				continue
			}

			for n, pn := range target.ns {
				if bytes.Equal(path, pn) {
					target.ns = append(target.ns[:n], target.ns[n+1:]...)
				}
			}

			break
		}
	}
	s.rw.RUnlock()
}

func (s *subCache) Add(sub Subscriber, path []byte) {
	var found bool
	s.rw.RLock()
	{
		for _, target := range s.cache {
			if target.sub != sub {
				continue
			}

			target.ns = append(target.ns, path)
			found = true
			break
		}
	}
	s.rw.RUnlock()

	if !found {
		s.rw.Lock()
		{
			s.cache = append(s.cache, subyList{
				sub: sub,
				ns:  [][]byte{path},
			})
		}
		s.rw.Unlock()
	}
}

func (s *subCache) Find(sub Subscriber) (subyList, bool) {
	var target subyList

	s.rw.RLock()
	{
		for _, target = range s.cache {
			if target.sub != sub {
				continue
			}

			return target, true
		}
	}
	s.rw.RUnlock()

	return target, false
}

//==============================================================================

type level struct {
	rw     sync.RWMutex
	tracer Trace
	all    *node
	nodes  map[string]*node
}

func newLevel(tracer Trace) *level {
	return &level{
		tracer: tracer,
		nodes:  make(map[string]*node),
		all: &node{
			sid:     containsSlice[0:],
			matcher: func(b []byte) bool { return true },
		},
	}
}

type node struct {
	next    *level
	sid     []byte
	ns      []byte
	subs    []Subscriber
	matcher func([]byte) bool
}

func (n *node) resolve(context interface{}, tracer Trace, tokens [][]byte, msg *Message) error {
	tLen := len(tokens)

	if tLen == 0 {
		return errors.New("Empty tokens for resolving")
	}

	token := tokens[0]

	if tLen == 1 {
		tokens = tokens[:0]
	} else {
		tokens = tokens[1:]
	}

	if !n.matcher(token) && !bytes.Equal(token, containsSlice) {
		return errors.New("token does not match route")
	}

	if !bytes.Equal(token, containsSlice) {
		msg.Params[string(n.ns)] = string(token)
	}

	msg.Match = bytes.Join([][]byte{msg.Match, n.sid}, sublistSlice)

	for _, sub := range n.subs {
		recovers(context, func() {
			if err := sub.Fire(context, msg); err != nil {
				if tracer != nil {
					tracer.Trace(context, []byte(fmt.Sprintf("Error firing for route %+s: %+s", n.sid, err.Error())))
				}
			}
		}, tracer)
	}

	if n.next != nil {
		if bytes.Equal(token, containsSlice) && len(tokens) == 0 {
			n.next.resolve(context, containsArraySlice, msg)
			return nil
		}

		n.next.resolve(context, tokens, msg)
	}

	return nil
}

// Routes returns the route for this level including its subroute.
func (s *level) Routes() [][]byte {
	var subs [][]byte

	s.rw.RLock()
	{
		for _, node := range s.nodes {
			subs = append(subs, node.sid)
			if node.next != nil {
				for _, lroute := range node.next.Routes() {
					subs = append(subs, bytes.Join([][]byte{node.sid, lroute}, sublistSlice))
				}
			}
		}
	}
	s.rw.RUnlock()

	return subs
}

// Size returns the total number of nodes on this level.
func (s *level) Size() int {
	s.rw.RLock()
	n := len(s.nodes)
	s.rw.RUnlock()
	return n
}

// Resolve checks if the giving path is a match within the giving level's routes.
func (s *level) Resolve(context interface{}, pattern []byte, msg *Message) {

	pLen := len(pattern)

	s.rw.RLock()
	tracer := s.tracer
	s.rw.RUnlock()

	if pLen == 0 {
		if s.tracer != nil {
			err := errors.New("Invalid Path to route")
			tracer.Trace(context, []byte(fmt.Sprintf("Error routing %+s: %+s", pattern, err.Error())))
		}
	}

	// if pLen == 1 && pattern[0] == contains {
	// 	s.rw.RLock()
	// 	{
	// 		for _, sub := range s.all.subs {
	// 			recovers(context, func() {
	// 				if err := sub.Fire(context, msg); err != nil {
	// 					if tracer != nil {
	// 						tracer.Trace(context, []byte(fmt.Sprintf("Error firing for route %+s: %+s", pattern, err.Error())))
	// 					}
	// 				}
	// 			}, tracer)
	// 		}
	// 	}
	// 	s.rw.RUnlock()
	// 	return
	// }

	tokens, err := splitResolveToken(pattern)
	if err != nil {
		if s.tracer != nil {
			s.tracer.Trace(context, []byte(fmt.Sprintf("Error routing %+s: %+s", pattern, err.Error())))
		}

		return
	}

	s.resolve(context, tokens, msg)
}

func (s *level) resolve(context interface{}, tokens [][]byte, msg *Message) {
	s.rw.RLock()
	tracer := s.tracer
	s.rw.RUnlock()

	s.rw.RLock()
	{
		for _, sub := range s.all.subs {
			recovers(context, func() {
				if err := sub.Fire(context, msg); err != nil {
					if tracer != nil {
						tracer.Trace(context, []byte(fmt.Sprintf("Error firing for route %+s: %+s", tokens, err.Error())))
					}
				}
			}, tracer)
		}
	}
	s.rw.RUnlock()

	s.rw.RLock()
	{
		for _, node := range s.nodes {
			if err := node.resolve(context, tracer, tokens, msg); err != nil && tracer != nil {
				tracer.Trace(context, []byte(fmt.Sprintf("Error routing %+s: %+s", tokens, err.Error())))
			}
		}
	}
	s.rw.RUnlock()

}

// Add adds a new subscriber into the subscription list with the provided pattern.
func (s *level) Add(pattern []byte, subscriber Subscriber) error {
	pLen := len(pattern)

	if pLen == 1 && pattern[0] == contains {
		s.all.subs = append(s.all.subs, subscriber)
		return nil
	}

	tokens := splitToken(pattern)
	return s.add(tokens, subscriber)
}

func (s *level) add(patterns [][]byte, subscriber Subscriber) error {
	pLen := len(patterns)

	for i := 0; i < len(patterns); i++ {
		item := patterns[i]
		itemLen := len(item)

		var match func([]byte) bool

		if bytes.Contains(item, startCurlyBracketSlice) && bytes.Contains(item, endCurlyBracketSlice) {
			word, regex := yankRegExp(item)

			if len(word) == 0 {
				return fmt.Errorf("Regexp token[%q] must include name", item)
			}

			matchex := regexp.MustCompile(string(regex))

			match = func(d []byte) bool {
				return matchex.Match(d)
			}

			s.rw.RLock()
			nodeItem, ok := s.nodes[string(item)]
			s.rw.RUnlock()
			if !ok {
				nodeItem = &node{
					next:    newLevel(s.tracer),
					sid:     item,
					ns:      word,
					matcher: match,
				}

				s.rw.Lock()
				s.nodes[string(item)] = nodeItem
				s.rw.Unlock()
			}

			if i+1 >= pLen {
				s.rw.Lock()
				nodeItem.subs = append(nodeItem.subs, subscriber)
				s.rw.Unlock()
				return nil
			}

			return nodeItem.next.add(patterns[i+1:], subscriber)
		}

		if bytes.Contains(item, edgesSlice) {
			if itemLen == 1 {
				return errors.New("Invalid Token usage, Edges('^') must be used at start or end of section")
			}

			ditem := item

			if bytes.HasPrefix(item, edgesSlice) {
				item = item[1:]

				match = func(d []byte) bool {
					if bytes.HasPrefix(d, item) {
						return true
					}

					return false
				}
			}

			if bytes.HasSuffix(item, edgesSlice) {
				item = item[:itemLen-1]

				match = func(d []byte) bool {
					if bytes.HasSuffix(d, item) {
						return true
					}

					return false
				}
			}

			s.rw.RLock()
			nodeItem, ok := s.nodes[string(ditem)]
			s.rw.RUnlock()
			if !ok {
				nodeItem = &node{
					next:    newLevel(s.tracer),
					sid:     ditem,
					ns:      item,
					matcher: match,
				}

				s.rw.Lock()
				s.nodes[string(ditem)] = nodeItem
				s.rw.Unlock()
			}

			if i+1 >= pLen {
				s.rw.Lock()
				nodeItem.subs = append(nodeItem.subs, subscriber)
				s.rw.Unlock()
				return nil
			}

			return nodeItem.next.add(patterns[i+1:], subscriber)
		}

		if bytes.Contains(item, containsSlice) {
			if itemLen == 1 && bytes.Equal(item, containsSlice) {
				s.all.subs = append(s.all.subs, subscriber)
				return nil
			}

			ditem := item
			item = bytes.Replace(item, containsSlice, emptySlice, 1)
			match = func(d []byte) bool {
				if bytes.Contains(d, item) {
					return true
				}

				return false
			}

			s.rw.RLock()
			nodeItem, ok := s.nodes[string(ditem)]
			s.rw.RUnlock()
			if !ok {
				nodeItem = &node{
					next:    newLevel(s.tracer),
					sid:     ditem,
					ns:      item,
					matcher: match,
				}

				s.rw.Lock()
				s.nodes[string(ditem)] = nodeItem
				s.rw.Unlock()
			}

			if i+1 >= pLen {
				s.rw.Lock()
				nodeItem.subs = append(nodeItem.subs, subscriber)
				s.rw.Unlock()
				return nil
			}

			return nodeItem.next.add(patterns[i+1:], subscriber)
		}

		match = func(d []byte) bool {
			return bytes.Equal(d, item)
		}

		s.rw.RLock()
		nodeItem, ok := s.nodes[string(item)]
		s.rw.RUnlock()
		if !ok {
			nodeItem = &node{
				next:    newLevel(s.tracer),
				sid:     item,
				ns:      item,
				matcher: match,
			}

			s.rw.Lock()
			s.nodes[string(item)] = nodeItem
			s.rw.Unlock()
		}

		if i+1 >= pLen {
			s.rw.Lock()
			nodeItem.subs = append(nodeItem.subs, subscriber)
			s.rw.Unlock()
			return nil
		}

		return nodeItem.next.add(patterns[i+1:], subscriber)
	}

	return nil
}

// Remove delets the subscriber from the subscription list with the provided pattern.
func (s *level) Remove(pattern []byte, subscriber Subscriber) error {
	pLen := len(pattern)

	if pLen == 1 && pattern[0] == contains {
		nodeItem := s.all

		s.rw.RLock()
		subs := nodeItem.subs
		subLen := len(nodeItem.subs)
		s.rw.RUnlock()

		for j := 0; j < subLen; j++ {
			if sub := subs[j]; sub == subscriber {
				if subLen == 1 {
					s.rw.Lock()
					nodeItem.subs = subs[:0]
					s.rw.Unlock()
					return nil
				}

				s.rw.Lock()
				nodeItem.subs[j] = subs[subLen-1]
				s.rw.Unlock()
				return nil
			}
		}

		return errors.New("Subscriber not found in registry")
	}

	tokens := splitToken(pattern)
	return s.remove(tokens, subscriber)
}

func (s *level) remove(patterns [][]byte, subscriber Subscriber) error {
	pLen := len(patterns)

	for i := 0; i < len(patterns); i++ {
		item := patterns[i]

		s.rw.RLock()
		nodeItem, ok := s.nodes[string(item)]
		s.rw.RUnlock()

		if !ok {
			return errors.New("Invalid route")
		}

		if i+1 >= pLen {
			s.rw.RLock()
			subs := nodeItem.subs
			s.rw.RUnlock()

			subLen := len(nodeItem.subs)

			for j := 0; j < subLen; j++ {
				if sub := subs[j]; sub == subscriber {
					if subLen == 1 {
						s.rw.Lock()
						nodeItem.subs = subs[:0]
						s.rw.Unlock()
						return nil
					}

					s.rw.Lock()
					nodeItem.subs[j] = subs[subLen-1]
					s.rw.Unlock()
					return nil
				}
			}

			return errors.New("Subscriber not found in registry")
		}

		return nodeItem.next.remove(patterns[i+1:], subscriber)
	}

	return nil
}

// PathToByte provides a quick function to transform a path string (`/ded/fr/fg`)
// into a slice of point bytes (`/ded/fr/fg => ded.fr.fg ==> []byte{11,33,...}`).
func PathToByte(path string) []byte {
	tokens := []byte(path)
	tokens = bytes.TrimSpace(tokens)

	if len(tokens) != 1 {
		tokens = bytes.TrimPrefix(tokens, slashSlice)
		tokens = bytes.TrimSuffix(tokens, slashSlice)
	}

	tokens = bytes.Replace(tokens, slashSlice, sublistSlice, -1)

	if len(tokens) == 1 && tokens[0] == sublist {
		tokens[0] = byte('/')
	}

	return tokens
}

func nsToken(token []byte) []byte {
	return bytes.Join([][]byte{colonSlice, token}, emptySlice)
}

func splitToken(pattern []byte) [][]byte {
	pLen := len(pattern)

	var tokens [][]byte
	var token []byte

	for i := 0; i < pLen; i++ {
		item := pattern[i]

		if item == sublist {
			tokens = append(tokens, token)
			token = nil
			continue
		}

		token = append(token, item)
	}

	if len(token) != 0 {
		tokens = append(tokens, token)
	}

	return tokens
}

func splitResolveToken(pattern []byte) ([][]byte, error) {
	var tokens [][]byte
	var token []byte

	pLen := len(pattern)

	for i := 0; i < pLen; i++ {
		item := pattern[i]

		switch item {
		case colon, edges, startCurlyBracket, endCurlyBracket, startBracket, endBracket:
			return nil, errors.New("Resolve Path cant contain pattern characters(`^`,`>`,`*`,`[`,`]`,`[]`")
		case sublist:
			tokens = append(tokens, token)
			token = nil
			continue
		default:
			token = append(token, item)
		}
	}

	if len(token) != 0 {
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func recovers(context interface{}, fx func(), tracer Trace) {
	defer func() {
		if err := recover(); err != nil {
			if tracer != nil {
				size := 1 << 16
				trace := make([]byte, size)
				trace = trace[:runtime.Stack(trace, true)]
				if tracer != nil {
					tracer.Trace(context, trace)
				}
			}
			return
		}
	}()

	fx()
}

// yankRegExp splits the giving pattern `id[\d+]`.
func yankRegExp(pattern []byte) ([]byte, []byte) {
	var word, rgu []byte

	var foundReg bool

	pattern = bytes.Replace(pattern, endCurlyBracketSlice, emptySlice, 1)
	pattern = bytes.Replace(pattern, startCurlyBracketSlice, emptySlice, 1)

	for _, item := range pattern {
		if item == colon {
			continue
		}

		if item == startBracket {
			foundReg = true
			continue
		}

		if item == endBracket {
			continue
		}

		if foundReg {
			rgu = append(rgu, item)
			continue
		}

		word = append(word, item)
	}

	return word, rgu
}
