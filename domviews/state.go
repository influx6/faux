package domviews

import (
	"errors"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

// excessStops is a regexp for matching more than one fullstops in the state address which then gets replaced into a single fullstop
var excessStops = regexp.MustCompile(`.\+`)

// ErrStateNotFound is returned when the state address is inaccurate and a state was not found in the path
var ErrStateNotFound = errors.New("State Not Found")

// ErrInvalidStateAddr is returned when a wrong length or format state address is found
var ErrInvalidStateAddr = errors.New("State Not Found")

// StateResponse defines the function type used by state in response to a state call
type StateResponse func()

//StateValidator defines a function type used in the state validator process
type StateValidator func(string, string) bool

// States represent the interface defining a state type
type States interface {
	Active() bool
	// Tag() string
	Engine() *StateEngine
	Activate()
	Deactivate()
	UseActivator(StateResponse) States
	UseDeactivator(StateResponse) States
	OverrideValidator(StateValidator) States
	acceptable(string, string) bool
}

// State represents a single state of with a specific tag and address
// where the address is a single piece 'home' item in the '.home.files' style state address
type State struct {
	// atomic bit used to indicate active state or inactive state
	active int64

	// tag represent the identifier key used in a super-state
	// tag string

	// activator and deactivator provide actions to occur when the state is set to be active
	activator   StateResponse
	deactivator StateResponse

	// validator represents an option argument which also takes part in the validation process of a state in validating its rightness/wrongness from a giving state address
	optionalValidator StateValidator

	//internal engine that allows sub-states from a root state
	engine *StateEngine

	// the parent state this is connected to
	// parent States

	vo, ro, do sync.Mutex
}

// NewState builds a new state with a tag and single address point .eg home or files ..etc
func NewState() *State {
	ns := State{}

	ns.engine = BuildStateEngine(&ns)

	return &ns
}

// Active returns true/false if this state is active
func (s *State) Active() bool {
	return atomic.LoadInt64(&s.active) == 1
}

// Engine returns the internal nested StateEngine
func (s *State) Engine() *StateEngine {
	return s.engine
}

// UseDeactivator assigns the state a new deactivate respone handler
func (s *State) UseDeactivator(so StateResponse) States {
	s.do.Lock()
	s.deactivator = so
	s.do.Unlock()
	return s
}

// UseActivator assigns the state a new active respone handler
func (s *State) UseActivator(so StateResponse) States {
	s.ro.Lock()
	s.activator = so
	s.ro.Unlock()
	return s
}

// OverrideValidator assigns an validator to perform custom matching of the state
func (s *State) OverrideValidator(so StateValidator) States {
	s.vo.Lock()
	s.optionalValidator = so
	s.vo.Unlock()
	return s
}

// Activate activates the state
func (s *State) Activate() {
	if s.active > 1 {
		return
	}

	atomic.StoreInt64(&s.active, 1)

	subs := s.engine.diffOnlySubs()

	//activate all the subroot states first so they can
	//do any population they want
	for _, ko := range subs {
		ko.Activate()
	}

	s.ro.Lock()

	if s.activator != nil {
		s.activator()
	}

	s.ro.Unlock()
}

// Deactivate deactivates the state
func (s *State) Deactivate() {
	if s.active < 1 {
		return
	}
	atomic.StoreInt64(&s.active, 0)

	s.do.Lock()
	if s.deactivator != nil {
		s.deactivator()
	}
	s.do.Unlock()
}

// acceptable checks if the state matches the current point
func (s *State) acceptable(addr string, point string) bool {
	if s.optionalValidator == nil {
		if addr == point {
			return true
		}
		return false
	}

	s.vo.Lock()
	state := s.optionalValidator(addr, point)
	s.vo.Unlock()
	return state
}

// StateEngine represents the engine that handles the state machine based operations for state-address based states
type StateEngine struct {
	rw     sync.RWMutex
	states map[States]string
	owner  States
	curr   States
}

// NewStateEngine returns a new engine with a default empty state
func NewStateEngine() *StateEngine {
	return BuildStateEngine(nil)
}

// BuildStateEngine returns a new StateEngine instance set with a particular state as its owner
func BuildStateEngine(s States) *StateEngine {
	es := StateEngine{
		states: make(map[States]string),
		owner:  s,
	}
	return &es
}

// AddState adds a new state into the engine with the tag used to identify the state, if the address is a empty string then the address recieves the tag as its value, remember the address is a single address point .eg home or files and not the length of the extend address eg .root.home.files
func (se *StateEngine) AddState(addr string) States {
	sa := NewState()
	se.add(addr, sa)

	return sa
}

// UseState adds a state into the StateEngine with a specific tag, the state address point is still used in matching against it
func (se *StateEngine) UseState(addr string, s States) States {
	if addr == "" {
		addr = "."
	}

	se.add(addr, s)

	return s
}

// ShallowState returns the current state of the engine and not the final state i.e with a state address of '.home.files' from its root, it will return State(:home) object
func (se *StateEngine) ShallowState() States {
	if se.curr == nil {
		return nil
	}

	return se.curr
}

// State returns the current last state of the engine with respect to any nested state that is with the state address of '.home.files', it will return State(:files) object
func (se *StateEngine) State() States {
	co := se.curr

	if co == nil {
		// return se.owner
		return nil
	}

	return co.Engine().State()
}

// Partial renders the partial of the last state of the state address
func (se *StateEngine) Partial(addr string) error {
	points, err := se.prepare(addr)

	if err != nil {
		return err
	}

	return se.trajectory(points, true)
}

// All renders the partial of the last state of the state address
func (se *StateEngine) All(addr string) error {
	points, err := se.prepare(addr)

	if err != nil {
		return err
	}

	return se.trajectory(points, false)
}

// DeactivateAll deactivates all states connected to this engine
func (se *StateEngine) DeactivateAll() {
	se.eachState(func(so States, tag string, _ func()) {
		so.Deactivate()
	})
}

func (se *StateEngine) eachState(fx func(States, string, func())) {
	if fx == nil {
		return
	}
	se.rw.RLock()
	defer se.rw.RUnlock()

	var stop bool

	for so, addr := range se.states {
		if stop {
			break
		}
		fx(so, addr, func() {
			stop = true
		})
	}
}

func (se *StateEngine) getAddr(s States) string {
	se.rw.RLock()
	defer se.rw.RUnlock()

	return se.states[s]
}

func (se *StateEngine) get(addr string) States {
	se.rw.RLock()
	defer se.rw.RUnlock()

	for sm, ao := range se.states {
		if ao != addr {
			continue
		}

		return sm
	}

	return nil
}

func (se *StateEngine) add(addr string, s States) {
	se.rw.RLock()
	_, ok := se.states[s]
	se.rw.RUnlock()

	if ok {
		return
	}

	se.rw.Lock()
	se.states[s] = addr
	se.rw.Unlock()
}

// trajectory is the real engine which checks the path and passes down the StateStat to the sub-states and determines wether its a full view or partial view
func (se *StateEngine) trajectory(points []string, partial bool) error {

	subs, nosubs := se.diffSubs()

	//are we out of points to walk through? if so then fire acive and tell others to be inactive
	if len(points) < 1 {
		//deactivate all the non-subroot states
		for _, ko := range nosubs {
			ko.Deactivate()
		}

		//if the engine has a root state activate it since in doing so,it will activate its own children else manually activate the children
		if se.owner != nil {
			se.owner.Activate()
		} else {
			//activate all the subroot states first so they can
			//do be ready for the root. We call this here incase the StateEngine has no root state
			for _, ko := range subs {
				ko.Activate()
			}
		}

		return nil
	}

	//cache the first point so we dont loose it
	point := points[0]

	var state = se.get(point)

	if state == nil {
		// for _, ko := range nosubs {
		// 	if sko.acceptable(se.getAddr(ko), point, so) {
		// 		state = ko
		// 		break
		// 	}
		// }
		//
		// if state == nil {
		return ErrStateNotFound
		// }
	}

	//set this state as the current active state
	se.curr = state

	//shift the list one more bit for the points
	points = points[1:]

	//we pass down the points since that will handle the loadup downwards
	err := state.Engine().trajectory(points, partial)

	if err != nil {
		return err
	}

	if !partial {
		// //activate all the subroot states first so they can
		// //do any population they want
		// for _, ko := range subs {
		// 	ko.Activate(so)
		// }

		if se.owner != nil {
			se.owner.Activate()
		}
	}

	return nil
}

// preparePoints prepares the state address into a list of walk points
func (se *StateEngine) prepare(addr string) ([]string, error) {

	var points []string
	var polen int

	if addr != "." {
		addr = excessStops.ReplaceAllString(addr, ".")
		points = strings.Split(addr, ".")
		polen = len(points)
	} else {
		polen = 1
		points = []string{""}
	}

	//check if the length is below 1 then return appropriately
	if polen < 1 {
		return nil, ErrInvalidStateAddr
	}

	//if the first is an empty string, meaning the '.' root was supplied, then we shift so we just start from the first state point else we ignore and use the list as-is.
	if points[0] == "" {
		points = points[1:]
	}

	return points, nil
}

// diffOnlyNotSubs returns all states with a '.' root state address
func (se *StateEngine) diffOnlySubs() []States {
	var subs []States

	se.eachState(func(so States, addr string, _ func()) {
		if addr == "." {
			subs = append(subs, so)
		}
	})

	return subs
}

// diffOnlyNotSubs returns all states not with a '.' root state address
func (se *StateEngine) diffOnlyNotSubs() []States {
	var subs []States

	se.eachState(func(so States, addr string, _ func()) {
		if addr != "." {
			subs = append(subs, so)
		}
	})

	return subs
}

func (se *StateEngine) diffSubs() ([]States, []States) {
	var nosubs, subs []States

	se.eachState(func(so States, addr string, _ func()) {
		if addr == "." {
			subs = append(subs, so)
		} else {
			nosubs = append(nosubs, so)
		}
	})

	return subs, nosubs
}
