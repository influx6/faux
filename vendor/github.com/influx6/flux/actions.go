package flux

import (
	"sync"
	"time"
)

//ActionInterface defines member functions
type ActionInterface interface {
	Fullfill(b interface{})
	WhenOnly(fx func(interface{})) ActionInterface
	When(fx func(interface{}, ActionInterface)) ActionInterface
	Then(fx func(interface{}, ActionInterface)) ActionInterface
	UseThen(fx func(interface{}, ActionInterface), a ActionInterface) ActionInterface
	Fullfilled() bool
	ChainAction(ActionInterface) ActionInterface
	ChainLastAction(ActionInterface) ActionInterface
	Chain(int) *ActDepend
	ChainWith(...ActionInterface) *ActDepend
	Wrap() *ActionWrap
	Sync(int) <-chan interface{}
}

//ActionStackInterface defines actionstack member method rules
type ActionStackInterface interface {
	Complete(b interface{}) ActionInterface
	Done() ActionInterface
	Error() ActionInterface
}

//ActionWrap safty wraps action for limited access to its fullfill function
type ActionWrap struct {
	action ActionInterface
}

//UnwrapAny returns two values where the first is not nil if the ActionInterface is a Action
//or the second non-nil if its a ActDepend
func UnwrapAny(a ActionInterface) (*Action, *ActDepend) {
	wa := UnwrapActionWrap(a)
	wad := UnwrapActDependWrap(a)

	return wa, wad
}

//UnwrapAction unwraps an ActionInterface to a *Action
func UnwrapAction(a ActionInterface) *Action {
	w, ok := a.(*Action)

	if !ok {
		return nil
	}

	return w
}

//UnwrapActionWrap unwraps an action that has being wrapped with ActionWrap
func UnwrapActionWrap(a ActionInterface) *Action {
	w, ok := a.(*ActionWrap)

	if !ok {
		return UnwrapAction(a)
	}

	ax, ok := w.action.(*Action)

	if !ok {
		return UnwrapActionWrap(ax)
	}

	return ax
}

//UnwrapActDependWrap unwraps an ActDepend that has being wrapped with ActionWrap
func UnwrapActDependWrap(a ActionInterface) *ActDepend {
	w, ok := a.(*ActionWrap)

	if !ok {
		return UnwrapActDepend(a)
	}

	ax, ok := w.action.(*ActDepend)

	if !ok {
		return UnwrapActDependWrap(ax)
	}

	return ax

}

//UnwrapActDepend unwraps an ActDepend that has being wrapped with ActionWrap
func UnwrapActDepend(a ActionInterface) *ActDepend {
	ax, ok := a.(*ActDepend)

	if ok {
		return ax
	}

	return nil
}

//NewActionWrap returns a action wrapped in a actionwrap
func NewActionWrap(a *Action) *ActionWrap {
	return &ActionWrap{a}
}

//Chain returns ActDepend(ActionDepend) with this action as the root
func (a *ActionWrap) Chain(m int) *ActDepend {
	return a.action.Chain(m)
}

//Fullfilled returns true or false if the action is done
func (a *ActionWrap) Fullfilled() bool {
	return a.action.Fullfilled()
}

//Fullfill meets this action of this structure
func (a *ActionWrap) Fullfill(b interface{}) {
	return
}

//WhenOnly adds a function to the action stack with the action as the second arg
func (a *ActionWrap) WhenOnly(fx func(b interface{})) ActionInterface {
	return a.action.WhenOnly(fx)
}

//When adds a function to the action stack with the action as the second arg
func (a *ActionWrap) When(fx func(b interface{}, a ActionInterface)) ActionInterface {
	return a.action.When(fx)
}

//ChainAction is a convenience method auto completes another action and returns that action,it uses
//UseThen underneath
func (a *ActionWrap) ChainAction(f ActionInterface) ActionInterface {
	return a.action.ChainAction(f)
}

//ChainLastAction is a convenience method auto completes another action and returns that action,it uses
//UseThen underneath
func (a *ActionWrap) ChainLastAction(f ActionInterface) ActionInterface {
	return a.action.ChainLastAction(f)
}

//Then adds a function to the action stack or fires immediately if done
func (a *ActionWrap) Then(fx func(b interface{}, a ActionInterface)) ActionInterface {
	return a.action.Then(fx)
}

//UseThen adds a function with a ActionInterface to the action stack or fires immediately if done
//once done that action interface is returned
func (a *ActionWrap) UseThen(fx func(b interface{}, a ActionInterface), f ActionInterface) ActionInterface {
	return a.action.UseThen(fx, a)
}

//Wrap returns actionwrap for the action
func (a *ActionWrap) Wrap() *ActionWrap {
	return a
}

//ChainWith returns ActDepend(ActionDepend) with this action as the root
func (a *ActionWrap) ChainWith(r ...ActionInterface) *ActDepend {
	return a.action.ChainWith(r...)
}

//Sync returns unbuffered channel which will get resolved with the
//value of the action when fullfilled
func (a *ActionWrap) Sync(ms int) <-chan interface{} {
	return a.action.Sync(ms)
}

//Action provides a future-style connect approach
type Action struct {
	fired bool
	stack *SingleStack
	cache interface{}
	lock  *sync.RWMutex
}

//Sync returns unbuffered channel which will get resolved with the
//value of the action when fullfilled or when the supplied value of
//time has passed it will eject
func (a *Action) Sync(ms int) <-chan interface{} {
	m := make(chan interface{})

	if a.cache != nil {
		go func() {
			m <- a.cache
			close(m)
		}()
	} else {
		if ms <= 0 {
			ms = 1
		}

		md := time.Duration(ms) * time.Millisecond

		a.When(func(b interface{}, _ ActionInterface) {
			go func() {
				m <- b
				close(m)
			}()
		})

		go func() {
			<-time.After(md)
			<-m
		}()
	}

	return m
}

//ChainLastAction is a convenience method auto completes another action and returns that action,it uses
//UseThen underneath
func (a *Action) ChainLastAction(f ActionInterface) ActionInterface {
	return a.ChainAction(f)
}

//ChainAction is a convenience method auto completes another action and returns that action,it uses
//UseThen underneath
func (a *Action) ChainAction(f ActionInterface) ActionInterface {
	if a == f {
		return a
	}

	return a.UseThen(func(b interface{}, to ActionInterface) {
		to.Fullfill(b)
	}, f)
}

//Chain returns ActDepend(ActionDepend) with this action as the root
func (a *Action) Chain(max int) *ActDepend {
	return NewActDepend(a, max)
}

//ChainWith returns ActDepend(ActionDepend) with this action as the root
func (a *Action) ChainWith(r ...ActionInterface) *ActDepend {
	return NewActDependWith(a, r...)
}

//ActDepend provides a nice means of creating a new action depending on
//unfullfilled action
type ActDepend struct {
	root    ActionInterface
	waiters []ActionInterface
	states  map[int]bool
	ind     int
	max     int
	ended   bool
}

//resolveType is used internall by ActDepend to fix a bug when trying to base
//other ActDepend to another especially when they share same roots,this
//cause a no-fullfillment as we are not using the .Then mechanism when calling
//.Mix or .MixLast
func resolveType(r ActionInterface) ActionInterface {
	ax, ad := UnwrapAny(r)

	if ax != nil {
		return ax
	}

	if ad != nil {
		xa := NewAction()
		ad.When(func(x interface{}, _ ActionInterface) {
			xa.Fullfill(x)
		})
		return xa
	}

	return r
}

//NewActDepend returns a action resolver based on a root action,when this root
//action is resolved,it waits on the user to call the actdepend then method to complete
//the next action,why so has to allow user-based chains where the user must partake in the
//completion of the final action
func NewActDepend(r ActionInterface, max int) *ActDepend {
	var set = make([]ActionInterface, max)
	var count = 0

	act := &ActDepend{
		resolveType(r),
		set,
		make(map[int]bool),
		0,
		max,
		false,
	}

	for count < max {
		m := NewAction()
		c := count
		m.When(func(b interface{}, _ ActionInterface) {
			act.states[c] = true
		})
		set[count] = m
		act.states[count] = false
		count++
	}

	return act
}

//NewActDependWith provides the actdepend struct but allows specifying the next call in the chan
func NewActDependWith(root ActionInterface, r ...ActionInterface) *ActDepend {
	var max = len(r)

	act := &ActDepend{
		resolveType(root),
		r,
		make(map[int]bool),
		0,
		max,
		false,
	}

	for k, m := range r {
		act.states[k] = false
		c := k
		m.When(func(b interface{}, _ ActionInterface) {
			act.states[c] = false
		})
	}

	return act
}

//NewActDependBy provides the actdepend struct but allows specifying the next call in the chan
func NewActDependBy(r ActionInterface, v ActionInterface, max int) *ActDepend {
	var set = make([]ActionInterface, max)
	var count = 1

	act := &ActDepend{
		resolveType(r),
		set,
		make(map[int]bool),
		0,
		max,
		false,
	}

	for count < max {
		m := NewAction()
		c := count
		m.When(func(b interface{}, _ ActionInterface) {
			act.states[c] = false
		})
		set[count] = m
		act.states[count] = false
		count++
	}

	return act
}

func (a *ActDepend) correctIndex(index int) (int, bool) {
	ind := 0

	if index < 0 {
		ind = a.Size() - ind
		if ind < 0 || ind > a.Size() {
			return 0, false
		}
	} else {
		ind = index
		if ind >= a.Size() {
			ind--
		}
	}

	return ind, true

}

//First returns the first ActionInterface in the dependency stack
func (a *ActDepend) First() ActionInterface {
	return a.getIndex(0)
}

//Last returns the last ActionInterface in the dependency stack
func (a *ActDepend) Last() ActionInterface {
	return a.getIndex(a.Size() - 1)
}

//Use returns the ActionInterface wrapped by an ActionWrap
//at the index or nil and supports negative indexing
func (a *ActDepend) Use(ind int) ActionInterface {
	return a.getIndex(ind).Wrap()
}

//GetIndex returns the ActionInterface at the index or nil
//supports negative indexing
func (a *ActDepend) getIndex(ind int) ActionInterface {
	ind, ok := a.correctIndex(ind)

	if !ok {
		return nil
	}

	return a.waiters[ind]
}

//IsIndexFullfilled returns true/false if the action at the index is fullfilled
func (a *ActDepend) IsIndexFullfilled(ind int) bool {
	return a.getIndex(ind).Fullfilled()
}

//Mix base the completion of action at a index with a custom action
//point using OverrideBefore and allows
//adding an extra step into the dependency action roadmap
//i.e when the next chain at this index which will complete the
//next chain if it is not the last as the normal operation of OverrideBefore
//it will base the completion of that next action on the action being mixed
//instead of the action at that index,like adding a middleman to a middleman :)
func (a *ActDepend) Mix(ind int, base ActionInterface) {
	a.OverrideBefore(ind, func(b interface{}, na ActionInterface) {
		base.ChainAction(na)
		base.Fullfill(b)
	})
}

//MixLast base adds a new action into the current action stack and
//calls inserts a ghost action inbetween the action at the index and the next
//action,when the ghost action is fullfilled the next action is fullfilled
//It underneaths calls the Action.ChainLastAction which when an ActDepend will
//resolve the next after the last action has been dissolved
func (a *ActDepend) MixLast(ind int, base ActionInterface) {
	a.OverrideBefore(ind, func(b interface{}, na ActionInterface) {
		base.ChainLastAction(na)
		base.Fullfill(b)
	})
}

//OverrideAfter allows calling Then with an action after the current index
//that is you want to listen to the action at this index to fullfill the
//next index
func (a *ActDepend) OverrideAfter(index int, fx func(b interface{}, a ActionInterface)) ActionInterface {
	ind, ok := a.correctIndex(index)

	if !ok {
		return a
	}

	ax := a.waiters[ind]
	var dx ActionInterface
	var nx = ind + 1

	if _, ok := a.states[nx]; ok {
		a.states[nx] = true
	}

	if nx >= a.Size() {
		return ax.Then(fx)
	}

	dx = a.waiters[nx]

	return ax.UseThen(fx, dx)

}

//OverrideBefore allows calling Then with an action before the current index
//that is you want to listen to the action at this previous index to fullfill the
//this action at this index
func (a *ActDepend) OverrideBefore(index int, fx func(b interface{}, a ActionInterface)) ActionInterface {
	ind, ok := a.correctIndex(index)

	if !ok {
		return a
	}

	ax := a.waiters[ind]

	var dx ActionInterface

	if ind < 1 {
		dx = a.root
	} else {
		dx = a.waiters[ind-1]
	}

	if _, ok := a.states[ind]; ok {
		a.states[ind] = true
	}

	return dx.UseThen(fx, ax)
}

//Shift pushes the index up if the action is not yet fully ended
//allows the next call of Then to be shifted as this states the user
//wants to manage this on their own
func (a *ActDepend) shift() {
	if a.ended {
		return
	}

	if a.ind >= a.Size() {
		return
	}

	a.ind++
}

//EqualRoot returns true/false if the root is equal
func (a *ActDepend) EqualRoot(r ActionInterface) bool {
	return a.root == r
}

//Size returns the total actions in list
func (a *ActDepend) Size() int {
	return len(a.waiters)
}

//current gets the last ActionInterface in the chain
func (a *ActDepend) current() ActionInterface {
	return a.waiters[a.ind]
}

//ChainWith returns ActDepend(ActionDepend) with this action as the root
func (a *ActDepend) ChainWith(r ...ActionInterface) *ActDepend {
	return a.current().ChainWith(r...)
}

//ChainLastAction is a convenience method auto completes another action and returns that action,it uses
//UseThen underneath
func (a *ActDepend) ChainLastAction(f ActionInterface) ActionInterface {
	return a.Last().ChainLastAction(f)
}

//ChainAction is a convenience method auto completes another action and returns that action,it uses
//UseThen underneath
func (a *ActDepend) ChainAction(f ActionInterface) ActionInterface {
	return a.current().ChainAction(f)
}

//Sync returns unbuffered channel which will get resolved with the
//value of the action when fullfilled
func (a *ActDepend) Sync(ms int) <-chan interface{} {
	return a.current().Sync(ms)
}

//End stops the generation of new chain
func (a *ActDepend) End() {
	a.ended = true
}

//Chain returns ActDepend(ActionDepend) with this action as the root
func (a *ActDepend) Chain(max int) *ActDepend {
	return a
}

//Wrap returns actionwrap for the action
func (a *ActDepend) Wrap() *ActionWrap {
	return &ActionWrap{a}
}

//Fullfilled returns true or false if the action is done
func (a *ActDepend) Fullfilled() bool {
	return a.current().Fullfilled()
}

//Fullfill actually fullfills the root action if its not fullfilled already
func (a *ActDepend) Fullfill(b interface{}) {
	a.root.Fullfill(b)
}

//WhenOnly adds a function to the action stack with the action as the second arg
func (a *ActDepend) WhenOnly(fx func(b interface{})) ActionInterface {
	return a.current().WhenOnly(fx)
}

//When adds a function to the action stack with the action as the second arg
func (a *ActDepend) When(fx func(b interface{}, a ActionInterface)) ActionInterface {
	return a.current().When(fx)
}

//Then adds a function to the action stack or fires immediately if done
func (a *ActDepend) Then(fx func(b interface{}, a ActionInterface)) ActionInterface {
	ind := a.ind
	sz := a.Size()

	if a.ended || (ind+1 > sz) {
		a.ended = true
		return a.waiters[sz-1].Then(fx)
	}

	cur := a.current()

	if a.states[ind] {
		a.ind++
		return a.Then(fx)
	}

	if ind <= 0 {
		_ = a.root.UseThen(fx, cur)
	} else {
		act := a.waiters[ind-1]
		act.UseThen(fx, cur)
	}

	a.states[ind] = true

	if ind < sz && (ind+1 < sz) {
		a.ind++
	} else {
		a.ended = true
	}

	return a
}

//UseThen adds a function with a ActionInterface to the action stack or fires immediately if done
//once done that action interface is returned
func (a *ActDepend) UseThen(fx func(b interface{}, a ActionInterface), f ActionInterface) ActionInterface {
	return a.current().UseThen(fx, a)
}

//Wrap returns actionwrap for the action
func (a *Action) Wrap() *ActionWrap {
	return NewActionWrap(a)
}

//Fullfilled returns true or false if the action is done
func (a *Action) Fullfilled() bool {
	return a.fired
}

//Fullfill meets this action of this structure
func (a *Action) Fullfill(b interface{}) {
	a.lock.RLock()
	state := a.fired
	a.lock.RUnlock()

	if state {
		return
	}

	a.cache = b
	a.lock.Lock()
	a.fired = true
	a.lock.Unlock()
	a.stack.Each(b)
	a.stack.Clear()
}

//WhenOnly adds a function to the action stack with the action as the second arg
func (a *Action) WhenOnly(fx func(b interface{})) ActionInterface {
	return a.When(func(v interface{}, _ ActionInterface) {
		fx(v)
	})
}

//When adds a function to the action stack with the action as the second arg
func (a *Action) When(fx func(b interface{}, e ActionInterface)) ActionInterface {
	a.lock.RLock()
	state := a.fired
	a.lock.RUnlock()

	if state {
		fx(a.cache, a)
	} else {
		a.stack.Add(func(res interface{}) {
			fx(res, a)
		})
	}

	return a
}

//Then adds a function to the action stack or fires immediately if done
func (a *Action) Then(fx func(b interface{}, a ActionInterface)) ActionInterface {
	ac := NewAction()
	if a.fired {
		fx(a.cache, ac)
	} else {
		a.stack.Add(func(res interface{}) {
			fx(res, ac)
		})
	}

	return ac
}

//UseThen adds a function with a ActionInterface to the action stack or fires immediately if done
//once done that action interface is returned
func (a *Action) UseThen(fx func(b interface{}, a ActionInterface), f ActionInterface) ActionInterface {
	if a.fired {
		fx(a.cache, f)
	} else {
		a.stack.Add(func(res interface{}) {
			fx(res, f)
		})
	}

	return f
}

//NewAction returns a new Action struct
func NewAction() *Action {
	return &Action{false, NewSingleStack(), nil, new(sync.RWMutex)}
}

//ActionStack provides two internal stack for success and error
type ActionStack struct {
	success ActionInterface
	failed  ActionInterface
}

//Done returns the action for the done state
func (a *ActionStack) Done() ActionInterface {
	return a.success
}

//Error returns the action for the error state
func (a *ActionStack) Error() ActionInterface {
	return a.failed
}

//Complete allows completion of an action stack
func (a *ActionStack) Complete(b interface{}) ActionInterface {
	if a.failed.Fullfilled() {
		return a.failed.Wrap()
	}

	if a.success.Fullfilled() {
		return a.success.Wrap()
	}

	e, ok := b.(error)

	if ok {
		a.failed.Fullfill(e)
		return a.failed.Wrap()
	}

	a.success.Fullfill(b)
	return a.success.Wrap()
}

//NewActionStack returns a new actionStack
func NewActionStack() *ActionStack {
	return &ActionStack{
		NewAction(),
		NewAction(),
	}
}

//ActionMod defines a function type that modifies a actionstack actions
//and returns them or the new actions
type ActionMod func(a ActionStackInterface) (ActionInterface, ActionInterface)

//NewActionStackFrom returns a new actionstack with the predefined actions
//from a previous actionstack with modification
func NewActionStackFrom(a ActionStackInterface, mod ActionMod) *ActionStack {
	d, e := mod(a)
	return &ActionStack{
		d,
		e,
	}
}

//NewActionStackBy returns a new actionstack with the predefined actions
//from a previous actionstack with modification
func NewActionStackBy(d ActionInterface, e ActionInterface) *ActionStack {
	return &ActionStack{
		d,
		e,
	}
}
