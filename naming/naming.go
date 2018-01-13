// Package naming defines a package for naming things, it provides us a baseline
// package for providing naming standards for things we use.
package naming

import "fmt"

// Name defines an interface which exposes a single method to deliver a new name
// based on an internal NameGenerator.
type Name interface {
	New() string
}

// FeedNamer defines an interface which exposes a single method to deliver a new name
// based on an internal NameGenerator and the supplied value.
type FeedNamer interface {
	New(string) string
}

// NameGenerator defines a type whch exposes a method which generates a new name
// based on the provided parameters.
type NameGenerator interface {
	Generate(template string, base string) string
}

//================================================================================

// Namer defines a struct which implements the Namer interface and uses the
// provided template and generator to generate a new name for use.
type Namer struct {
	template string
	gen      NameGenerator
}

// NewNamer returns a new instance of the NamerCon struct.
func NewNamer(template string, generator NameGenerator) *Namer {
	return &Namer{
		template: template,
		gen:      generator,
	}
}

// New returns a new name based on the provided base value and the internal template
// and NameGenerator.
func (n *Namer) New(base string) string {
	return n.gen.Generate(n.template, base)
}

//================================================================================

// Names defines a struct which implements the Namer interface and uses the
// provided template and generator to generate a new name for use.
type Names struct {
	template  string
	base      string
	generator NameGenerator
}

// NewNames returns a new instance of the NamerCon struct.
func NewNames(template string, base string, generator NameGenerator) *Names {
	return &Names{
		base:      base,
		template:  template,
		generator: generator,
	}
}

// New returns a new name based on the provided he internal template
// and generator.
func (n *Names) New() string {
	return n.generator.Generate(n.template, n.base)
}

//================================================================================

// GenerateNamer defines a function which accepts a base generator and uses the provided
// template and returns a function which will use this template and generator for all
// name generation.
func GenerateNamer(gen NameGenerator, template string) func(string) string {
	return func(base string) string {
		return gen.Generate(template, base)
	}
}

//================================================================================

// SuffixNamer returns a struct which returns a new name based on a giving
// suffix.
type SuffixNamer struct {
	Suffix string
}

// Generate returns a new name based on the internal suffix and provided
// template.
func (s SuffixNamer) Generate(template, base string) string {
	return fmt.Sprintf(template, s.Suffix, base)
}

//================================================================================

// PrefixNamer returns a struct which returns a new name based on a giving
// prefix.
type PrefixNamer struct {
	Prefix string
}

// Generate returns a new name based on the internal prefix and provided
// template.
func (p PrefixNamer) Generate(template, base string) string {
	return fmt.Sprintf(template, p.Prefix, base)
}

//================================================================================

// Basic defines a struct which implements the NameGenerator interface and generates
// a giving name based on the template and a giving set of directives.
type Basic struct{}

// Generate returns a new name based on the provided arguments.
func (Basic) Generate(template string, base string) string {
	return fmt.Sprintf(template, base)
}

//================================================================================

// LimitedNamer defines a struct which implements the NameGenerator interface and generates
// a giving name based on the template and a giving set of directives.
type LimitedNamer struct {
	maxLen     int
	maxBaseLen int
}

// NewLimitNamer returns a new instance of a LimitedNamer.
func NewLimitNamer(maxChars int, maxbaseChar int) *LimitedNamer {
	return &LimitedNamer{
		maxLen:     maxChars,
		maxBaseLen: maxbaseChar,
	}
}

// Generate returns a new name from the provided template and base to generate based
// on the rules of the simple rules.
func (l LimitedNamer) Generate(template string, base string) string {
	if len(base) > l.maxBaseLen {
		base = base[:l.maxBaseLen]
	}

	rem := l.maxLen - l.maxBaseLen
	gened := fmt.Sprintf(template, base, String(rem))

	if len(gened) > l.maxLen {
		gened = gened[:l.maxLen]
	}

	return gened
}

//=================================================================================

var (
	maxCharacter      = 60
	maxBaseCharacter  = 10
	maxAvailableSpace = maxCharacter - maxBaseCharacter
)

// SimpleNamer defines a struct which implements the NameGenerator interface and generates
// a giving name based on the template and a giving set of directives.
type SimpleNamer struct{}

// Generate returns a new name from the provided template and base to generate based
// on the rules of the simple rules.
func (s SimpleNamer) Generate(template string, base string) string {
	if len(base) > maxBaseCharacter {
		base = base[:maxBaseCharacter]
	}

	gened := fmt.Sprintf(template, base, String(maxBaseCharacter))

	if len(gened) > maxCharacter {
		gened = gened[:maxCharacter]
	}

	return gened
}
