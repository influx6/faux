// Package metrics defines a basic structure foundation for handling logs without
// much hassle, allow more different entries to be created.
// Inspired by https://medium.com/@tjholowaychuk/apex-log-e8d9627f4a9a.
package metrics

// Processors implements a single method to process a Entry.
type Processors interface {
	Handle(Entry) error
}

// Collector defines an interface which exposes a single method to collect
// internal data which is then returned as an Entry.
type Collector interface {
	Collect(string) Entry
}

// Metrics defines an interface with a single method for receiving
// new Entry objects.
type Metrics interface {
	Send(Entry) error
	Emit(...EntryMod) error
	CollectMetrics(string) error
}

// New returns a Metrics object with the provided Augmenters and  Metrics
// implemement objects for receiving metric Entries.
func New(vals ...interface{}) Metrics {
	var mods []EntryMod
	var procs []Processors
	var collectors []Collector

	for _, val := range vals {
		switch item := val.(type) {
		case Collector:
			collectors = append(collectors, item)
		case EntryMod:
			mods = append(mods, item)
		case Processors:
			procs = append(procs, item)
		}
	}

	var modder EntryMod
	if len(mods) != 0 {
		modder = Partial(mods...)
	}

	return metrics{
		collectors: collectors,
		processors: procs,
		mod:        modder,
	}
}

type metrics struct {
	mod        EntryMod
	processors []Processors
	collectors []Collector
}

// CollectMetrics runs internal indepent collectors to
// grap metrics.
func (m metrics) CollectMetrics(id string) error {
	for _, collector := range m.collectors {
		if err := m.Send(collector.Collect(id)); err != nil {
			return err
		}
	}
	return nil
}

// Send delivers Entry to processors
func (m metrics) Send(en Entry) error {
	for _, met := range m.processors {
		if err := met.Handle(en); err != nil {
			return err
		}
	}
	return nil
}

// Emit implements the Metrics interface and delivers Entry
// to undeline metrics.
func (m metrics) Emit(mods ...EntryMod) error {
	if len(m.processors) == 0 || len(mods) == 0 {
		return nil
	}

	var en Entry
	Apply(&en, mods...)
	if m.mod != nil {
		m.mod(&en)
	}

	return m.Send(en)
}

// FilterLevel will return a metrics where all Entry will be filtered by their Entry.Level
// if the level giving is greater or equal to the provided, then it will be received by
// the metrics subscribers.
func FilterLevel(l Level, procs ...Processors) Processors {
	return Case(func(en Entry) bool { return en.Level >= l }, procs...)
}

// DoFn defines a function type which takes a giving Entry.
type DoFn func(Entry) error

type fnMetrics struct {
	do DoFn
}

// DoWith returns a Metrics object where all entries are applied to the provided function.
func DoWith(do DoFn) Processors {
	return fnMetrics{
		do: do,
	}
}

// Handle implements the Processors interface and delivers Entry
// to undeline metrics.
func (m fnMetrics) Handle(en Entry) error {
	return m.do(en)
}

// ConditionalProcessors defines a Processor which first validate it's
// ability to process a giving Entry.
type ConditionalProcessors interface {
	Processors
	Can(Entry) bool
}

// FilterFn defines a function type which takes a giving Entry returning a bool to indicate filtering state.
type FilterFn func(Entry) bool

type caseProcessor struct {
	condition FilterFn
	procs     []Processors
}

// Case returns a Processor object with the provided Augmenters and  Metrics
// implemement objects for receiving metric Entries, where entries are filtered
// out based on a provided function.
func Case(fn FilterFn, procs ...Processors) ConditionalProcessors {
	return caseProcessor{
		condition: fn,
		procs:     procs,
	}
}

// Can returns true/false if we can handle giving Entry.
func (m caseProcessor) Can(en Entry) bool {
	return m.condition(en)
}

// Handle implements the Processors interface and delivers Entry
// to undeline metrics.
func (m caseProcessor) Handle(en Entry) error {
	if m.condition(en) {
		for _, proc := range m.procs {
			if err := proc.Handle(en); err != nil {
				return err
			}
		}
	}
	return nil
}

// EntryEmitter defines a type which returns a entry when runned.
type EntryEmitter func(string) Entry

// Collect returns a Collector which executes provided function when
// called by Metric to run.
func Collect(fn EntryEmitter) Collector {
	return fnCollector{fn: fn}
}

// fnCollector implements the Collector interface.
type fnCollector struct {
	fn EntryEmitter
}

// Collect runs the internal function and returning the produced entry.
func (fn fnCollector) Collect(name string) Entry {
	return fn.fn(name)
}

// switchMaster defines that mod out Entry objects based on a provided function.
type switchMaster struct {
	cases []ConditionalProcessors
}

// Switch returns a new instance of a SwitchMaster.
func Switch(conditions ...ConditionalProcessors) Processors {
	return switchMaster{
		cases: conditions,
	}
}

// Handle delivers the giving entry to all available metricss.
func (fm switchMaster) Handle(e Entry) error {
	for _, proc := range fm.cases {
		if proc.Can(e) {
			if err := proc.Handle(e); err != nil {
				return err
			}
		}
	}
	return nil
}
