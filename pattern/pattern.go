// Package pattern provides a simple regexp pattern matching library majorly
// for constructing URL matchers.
//  Patterns in this package follow the follow approach in declaring custom match
// segments.
//
// 		pattern: /name/{id:[/\d+/]}/log/{date:[/\w+\W+/]}
// 		pattern: /name/:id
//
//
package pattern

import (
	"regexp"
	"strings"
)

// Params defines a map of stringed keys and values.
type Params map[string]string

// Matchable defines an interface for matchers.
type Matchable interface {
	IsParam() bool
	Segment() string
	Validate(string) bool
}

// Matchers defines a list of machers for validating patterns with.
type Matchers []Matchable

// URIMatcher defines an interface for a URI matcher.
type URIMatcher interface {
	Validate(string) (Params, string, bool)
	Pattern() string
	Priority() int
}

// matchProvider provides a class array-path matcher
type matchProvider struct {
	pattern  string
	matchers Matchers
	endless  bool
	priority int
}

// New returns a new instance of a URIMatcher.
func New(pattern string) URIMatcher {

	ps := stripAndClean(pattern)

	pm := SegmentList(ps)

	m := matchProvider{
		priority: CheckPriority(pattern),
		pattern:  pattern,
		matchers: pm,
		endless:  IsEndless(pattern),
	}

	return &m
}

// Priority returns the priority status of this giving pattern.
func (m *matchProvider) Priority() int {
	return m.priority
}

// Pattern returns the pattern string for this matcher.
func (m *matchProvider) Pattern() string {
	return m.pattern
}

// Validate returns true/false if the giving string matches the pattern, returning
// a map of parameters match against segments of the pattern.
func (m *matchProvider) Validate(f string) (Params, string, bool) {
	stripped := stripAndClean(f)
	cleaned := cleanPath(stripped)
	src := splitPattern(cleaned)

	total := len(m.matchers)
	srclen := len(src)

	if !m.endless && (total < srclen || total > srclen) {
		return nil, "", false
	}

	var state bool

	param := make(Params)

	for k, v := range m.matchers {
		if k >= srclen {
			state = false
			break
		}

		if v.Validate(src[k]) {

			if v.IsParam() {
				param[v.Segment()] = src[k]
			}

			state = true
			continue
		} else {
			state = false
			break
		}
	}

	var rem string
	if total < srclen {
		csrc := stripAndCleanButHash(f)
		hashIndex := strings.IndexRune(csrc, '#')

		fsrc := stripAndClean(strings.Join(src[:total], "/"))
		fcount := len([]byte(fsrc))

		if hashIndex < fcount {
			rem = strings.Replace(stripped, fsrc, "", 1)
		} else {
			rem = strings.Replace(csrc, fsrc, "", 1)
		}
		// fmt.Printf("Rem: %s : %s -> %s\n", csrc, fsrc, rem)
	}

	return param, rem, state
}

//==============================================================================

// SegmentList returns list of SegmentMatcher which implements the Matchable
// interface, with each made of each segment of the pattern.
func SegmentList(pattern string) Matchers {
	var set Matchers

	for _, val := range splitPattern(pattern) {
		set = append(set, Segment(val))
	}

	return set
}

//==============================================================================

// SegmentMatcher defines a single piece of pattern to be matched against.
type SegmentMatcher struct {
	*regexp.Regexp
	original string
	param    bool
}

// Segment returns a Matchable for a specific part of a pattern eg. :name, age,
// {id:[\\d+]}.
func Segment(segment string) Matchable {
	id, rx, b := YankSpecial(segment)
	mrk := regexp.MustCompile(rx)

	sm := SegmentMatcher{
		Regexp:   mrk,
		original: id,
		param:    b,
	}

	return &sm
}

// IsParam returns true/false if the segment is also a paramter.
func (s *SegmentMatcher) IsParam() bool {
	return s.param
}

// Segment returns the original string that makes up this segment matcher.
func (s *SegmentMatcher) Segment() string {
	return s.original
}

// Validate validates the value against the matcher.
func (s *SegmentMatcher) Validate(m string) bool {
	return s.MatchString(m)
}

//==============================================================================
