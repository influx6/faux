package pattern

import (
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

//==============================================================================

var endless = regexp.MustCompile(`/\*$`)

//IsEndless returns true/false if the pattern as a /*
func IsEndless(s string) bool {
	return endless.MatchString(s)
}

//==============================================================================

// moreSlashes this to check for more than one forward slahes
var moreSlashes = regexp.MustCompile(`\/+`)

// CleanSlashes cleans all double forward slashes into one
func CleanSlashes(p string) string {
	if strings.Contains(p, "\\") {
		return moreSlashes.ReplaceAllString(filepath.ToSlash(p), "/")
	}
	return moreSlashes.ReplaceAllString(p, "/")
}

//==============================================================================

//RemoveCurly removes '{' and '}' from any string
func RemoveCurly(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{")
}

//==============================================================================

//RemoveBracket removes '[' and ']' from any string
func RemoveBracket(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
}

//==============================================================================

//SplitPattern splits a pattern with the '/'
func SplitPattern(c string) []string {
	return strings.Split(c, "/")
}

// addSlash appends a / infront of the giving string if not there.
func addSlash(ps string) string {
	if strings.HasPrefix(ps, "/") {
		return ps
	}

	return "/" + ps
}

//==============================================================================

//TrimEndSlashe removes the '/' at the end of string.
func TrimEndSlashe(c string) string {
	return strings.TrimSuffix(cleanPath(c), "/")
}

//==============================================================================

//TrimSlashes removes the '/' at the beginning and end of string.
func TrimSlashes(c string) string {
	return strings.TrimSuffix(strings.TrimPrefix(cleanPath(c), "/"), "/")
}

//==============================================================================

//SplitPatternAndRemovePrefix splits a pattern with the '/'
func SplitPatternAndRemovePrefix(c string) []string {
	return strings.Split(strings.TrimPrefix(cleanPath(c), "/"), "/")
}

//==============================================================================

var morespecial = regexp.MustCompile(`{\w+:[\w\W]+}`)

// HasKeyParam returns true/false if the special pattern {:[..]} exists in the string
func HasKeyParam(p string) bool {
	return morespecial.MatchString(p)
}

// CheckPriority is used to return the priority of a pattern.
// 0 for highest(when no parameters).
// 1 for restricted parameters({id:[]}).
// 2 for no paramters.
// The first parameter catched is used for rating.
// The ratings go from highest to lowest .i.e (0-2).
func CheckPriority(patt string) int {
	sets := splitPattern(patt)

	for _, so := range sets {
		if morespecial.MatchString(so) {
			return 1
		}

		if special.MatchString(so) {
			return 2
		}

		continue
	}

	return 0
}

//cleanPattern cleans any /* * pattern found
func cleanPattern(patt string) string {
	cleaned := endless.ReplaceAllString(patt, "")
	return morespecial.ReplaceAllString(cleaned, "/")
}

//==============================================================================

// CleanPath provides a public path cleaner
func CleanPath(p string) string {
	return cleanPath(p)
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

//==============================================================================

var special = regexp.MustCompile(`{\w+:[\w\W]+}|:[\w]+`)

// HasParam returns true/false if the special pattern {:[..]} exists in the string
func HasParam(p string) bool {
	return special.MatchString(p)
}

//==============================================================================

var picker = regexp.MustCompile(`^:[\w\W]+$`)

// HasPick matches string of type :id,:name
func HasPick(p string) bool {
	return picker.MatchString(p)
}

//==============================================================================

var specs = regexp.MustCompile(`\n\t\s+`)
var paramd = regexp.MustCompile(`^{[\w\W]+}$`)
var anyvalue = `[\w\W]+`

//YankSpecial provides a means of extracting parts of form `{id:[\d+]}`
func YankSpecial(val string) (string, string, bool) {
	if HasPick(val) {
		cls := strings.TrimPrefix(val, ":")
		return cls, anyvalue, true
	}

	if !paramd.MatchString(val) {
		cls := specs.ReplaceAllString(val, "")
		return cls, cls, false
	}

	part := strings.Split(removeCurly(val), ":")
	return part[0], removeBracket(part[1]), true
}

//==============================================================================

func removeCurly(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{")
}

func removeBracket(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
}

func splitPattern(c string) []string {
	parts := strings.Split(c, "/")

	// Re-add the first slash to respect root supremacy.
	if len(parts) > 0 && parts[0] == "" {
		parts[0] = "/"
	}

	return parts
}

//==============================================================================

// stripAndClean strips the slahes from the path.
func stripAndClean(c string) string {
	return CleanSlashes(strings.Replace(strings.TrimSuffix(strings.TrimSuffix(c, "/*"), "/"), "#", "/", -1))
}

// stripAndCleanButHash strips the slahes from the path.
func stripAndCleanButHash(c string) string {
	return CleanSlashes(strings.TrimSuffix(strings.TrimSuffix(c, "/*"), "/"))
}
