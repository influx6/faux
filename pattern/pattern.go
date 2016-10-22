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
	HasHash() bool
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
	pattern = addSlash(pattern)

	pm := SegmentList(pattern)

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
func (m *matchProvider) Validate(path string) (Params, string, bool) {
	path = addSlash(path)
	stripped := stripAndClean(path)
	hashedSrc := stripAndCleanButHash(path)

	cleaned := cleanPath(stripped)
	src := splitPattern(cleaned)

	srclen := len(src)
	total := len(m.matchers)

	if !m.endless && total != srclen {
		return nil, "", false
	}

	if m.endless && total > srclen {
		return nil, "", false
	}

	var state bool

	param := make(Params)

	var lastIndex int
	var doneHash bool

	for index, v := range m.matchers {
		lastIndex = index

		if v.HasHash() {
			doneHash = true
		}

		if v.Validate(src[index]) {

			if v.IsParam() {
				param[v.Segment()] = src[index]
			}

			state = true
			continue
		} else {
			state = false
			break
		}
	}

	if lastIndex+1 == srclen {
		return param, "", state
	}

	remPath := strings.Join(src[lastIndex+1:], "/")
	if doneHash || !strings.Contains(hashedSrc, "#") {
		return param, addSlash(remPath), state
	}

	var rems []string

	fragment := SegmentList(hashedSrc)[lastIndex+1:]
	for _, item := range fragment {
		if item.HasHash() {
			hashed := "#" + item.Segment()
			rems = append(rems, hashed)
			continue
		}

		rems = append(rems, item.Segment())
	}

	return param, addSlash(strings.Join(rems, "/")), state
}

//==============================================================================

// SegmentList returns list of SegmentMatcher which implements the Matchable
// interface, with each made of each segment of the pattern.
func SegmentList(pattern string) Matchers {
	pattern = stripAndCleanButHash(pattern)

	var set Matchers

	if hashIndex := strings.Index(pattern, "#"); hashIndex != -1 {
		if hashIndex == 0 {
			pattern = strings.Join([]string{"/", pattern}, "")
		} else {
			last := pattern[hashIndex-1 : hashIndex]
			if string(last[0]) != "/" {
				splits := strings.Split(pattern, "#")
				pattern = strings.Join([]string{splits[0], "/#", splits[1]}, "")
			}
		}
	}

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
	hashed   bool
}

// Segment returns a Matchable for a specific part of a pattern eg. :name, age,
// {id:[\\d+]}.
func Segment(segment string) Matchable {
	if segment == "*" {
		segment = "/*"
	}

	hashed := strings.HasPrefix(segment, "#")
	if hashed {
		segment = segment[1:]
	}

	id, rx, b := YankSpecial(segment)
	mrk := regexp.MustCompile(rx)

	sm := SegmentMatcher{
		Regexp:   mrk,
		original: id,
		param:    b,
		hashed:   hashed,
	}

	return &sm
}

// HasHashed returns true/false if this segment hash the hash.
func (s *SegmentMatcher) HasHash() bool {
	return s.hashed
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
