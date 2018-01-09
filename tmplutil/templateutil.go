package tmplutil

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

var (
	defaultFuncs = template.FuncMap{
		"greaterThanEqualF": func(b, a float64) bool {
			return b >= a
		},
		"lessThanEqualF": func(b, a float64) bool {
			return b <= a
		},
		"greaterThanEqual": func(b, a int) bool {
			return b >= a
		},
		"capitalize": func(b string) string {
			if len(b) == 0 {
				return b
			}

			return strings.ToUpper(b[:1]) + b[1:]
		},
		"notempty": func(b string) bool {
			return strings.TrimSpace(b) != ""
		},
		"empty": func(b string) bool {
			return strings.TrimSpace(b) == ""
		},
		"title": func(b string) string {
			return strings.ToTitle(b)
		},
		"trim": func(b, suff string) string {
			return strings.Trim(b, suff)
		},
		"trimSuffix": func(b, suff string) string {
			return strings.TrimSuffix(b, suff)
		},
		"trimPrefix": func(b, pre string) string {
			return strings.TrimPrefix(b, pre)
		},
		"hasSuffix": func(b, suff string) bool {
			return strings.HasSuffix(b, suff)
		},
		"hasPrefix": func(b, pre string) bool {
			return strings.HasPrefix(b, pre)
		},
		"replaceOnce": func(b, target, sub string) string {
			return strings.Replace(b, target, sub, 1)
		},
		"replaceAll": func(b, target, sub string) string {
			return strings.Replace(b, target, sub, -1)
		},
		"lower": func(b string) string {
			return strings.ToLower(b)
		},
		"upper": func(b string) string {
			return strings.ToUpper(b)
		},
		"joinPath": func(b ...string) string {
			return path.Join(b...)
		},
		"basePathName": func(b string) string {
			return path.Base(b)
		},
		"join": func(vals []string, jn string) string {
			return strings.Join(vals, jn)
		},
		"joinInterface": func(vals []interface{}, jn string) string {
			var items []string
			for _, val := range vals {
				items = append(items, fmt.Sprintf("%+s", val))
			}
			return strings.Join(items, jn)
		},
		"joinSlice": func(vals []string, jn string) string {
			return strings.Join(vals, jn)
		},
		"joinVariadic": func(jn string, vals ...string) string {
			return strings.Join(vals, jn)
		},
		"splitAfter": func(b string, sp string, n int) []string {
			return strings.SplitAfterN(b, sp, n)
		},
		"split": func(b string, sp string) []string {
			return strings.Split(b, sp)
		},
		"indent": func(b string) string {
			return strings.Join(strings.Split(b, "\n"), "\n\t")
		},
		"lessThanEqual": func(b, a int) bool {
			return b <= a
		},
		"bytesToString": func(b []byte) string {
			return string(b)
		},
		"greaterThanF": func(b, a float64) bool {
			return b > a
		},
		"lessThanF": func(b, a float64) bool {
			return b < a
		},
		"greaterThan": func(b, a int) bool {
			return b > a
		},
		"lessThan": func(b, a int) bool {
			return b < a
		},
		"trimspace": func(b string) string {
			return strings.TrimSpace(b)
		},
		"equal": func(b, a interface{}) bool {
			return b == a
		},
		"not": func(b bool) bool {
			return !b
		},
		"notequal": func(b, a interface{}) bool {
			return b != a
		},
		"quote": quote,
		"prefixInt": func(prefix string, b int) string {
			return fmt.Sprintf("%s%d", prefix, b)
		},
		"stub": func(count int) string {
			var vals []string

			for i := count; i > 0; i-- {
				vals = append(vals, "_")
			}

			return strings.Join(vals, ",")
		},
		"subs": func(word string, b int) string {
			return word[:b]
		},
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"divide": func(a, b int) int {
			return a / b
		},
		"lenNotEqual": func(b interface{}, target int) bool {
			return lenOff(b) != target
		},
		"lenEqual": func(b interface{}, target int) bool {
			return lenOff(b) == target
		},
		"percentage": func(a, b float64) float64 {
			return (a / b) * 100
		},
		"lenOf":         lenOff,
		"nthOf":         nthOf,
		"doCut":         cutList,
		"doTimesPrefix": doTimePrefix,
		"doTimesSuffix": doTimeSuffix,
		"doPrefix":      doPrefix,
		"doSuffix":      doSuffix,
		"doCutSplit":    doCutSplit,
		"doPrefixCut":   cutListPrefix,
		"doSuffixCut":   cutListSuffix,
		"intsToString":  doIntToString,
	}
)

//======================================================================================

// From returns a new template.Template instance which has the provided functions added.
func From(tag string, data string) (*template.Template, error) {
	return template.New(tag).Funcs(defaultFuncs).Parse(data)
}

// MustFrom returns a new template.Template instance with all functions added.
// It panics if instance failed to be created with returned error.
func MustFrom(tag string, data string) *template.Template {
	tmpl, err := From(tag, data)
	if err != nil {
		panic(err)
	}
	return tmpl
}

//======================================================================================

// Group holds a series of template content with associated name.
// Group returns new templates based on provided sets.
type Group struct {
	ml        sync.RWMutex
	templates map[string]string
}

// New returns a new instance of Group.
func New() *Group {
	return &Group{
		templates: make(map[string]string),
	}
}

// Add adds giving name and template value into group.
func (g *Group) Add(name string, data string) *Group {
	g.ml.Lock()
	defer g.ml.Unlock()
	g.templates[name] = data
	return g
}

// FromWith acts like New but calls template.Lookup on the retunred template
// to return the first template object forthe first name in the slice provided.
func (g *Group) FromWith(fn template.FuncMap, names ...string) (*template.Template, error) {
	if len(names) == 0 {
		return nil, errors.New("Require name slice")
	}

	tml, err := g.NewWith(fn, names...)
	if err != nil {
		return nil, err
	}

	if initialTml := tml.Lookup(names[0]); initialTml != nil {
		return initialTml, nil
	}

	return nil, fmt.Errorf("template object for %+q not found", names[0])
}

// NewWith returns a template from all templates within group if no names
// are provided. Either arrange names of support template to main template
// else use template.Lookup method to retrieve template for respective name.
func (g *Group) NewWith(fn template.FuncMap, names ...string) (*template.Template, error) {
	g.ml.RLock()
	defer g.ml.RUnlock()

	var tmpl *template.Template

	if len(names) != 0 {
		for _, name := range names {
			data, ok := g.templates[name]
			if !ok {
				continue
			}

			if tmpl != nil {
				ml, err := tmpl.New(name).Funcs(defaultFuncs).Funcs(fn).Parse(data)
				if err != nil {
					return nil, err
				}

				tmpl = ml
				continue
			}

			tml, err := template.New(name).Funcs(defaultFuncs).Funcs(fn).Parse(data)
			if err != nil {
				return nil, err
			}

			tmpl = tml
		}

		return tmpl, nil
	}

	for name, data := range g.templates {
		if tmpl != nil {
			ml, err := tmpl.New(name).Funcs(defaultFuncs).Funcs(fn).Parse(data)
			if err != nil {
				return nil, err
			}

			tmpl = ml
			continue
		}

		tml, err := template.New(name).Funcs(defaultFuncs).Funcs(fn).Parse(data)
		if err != nil {
			return nil, err
		}

		tmpl = tml
	}

	return tmpl, nil
}

// From acts like New but calls template.Lookup on the retunred template
// to return the first template object forthe first name in the slice provided.
func (g *Group) From(names ...string) (*template.Template, error) {
	if len(names) == 0 {
		return nil, errors.New("Require name slice")
	}

	tml, err := g.New(names...)
	if err != nil {
		return nil, err
	}

	if initialTml := tml.Lookup(names[0]); initialTml != nil {
		return initialTml, nil
	}

	return nil, fmt.Errorf("template object for %+q not found", names[0])
}

// New returns a template from all templates within group if no names
// are provided. Either arrange names of support template to main template
// else use template.Lookup method to retrieve template for respective name.
func (g *Group) New(names ...string) (*template.Template, error) {
	g.ml.RLock()
	defer g.ml.RUnlock()

	var tmpl *template.Template

	if len(names) != 0 {
		for _, name := range names {
			data, ok := g.templates[name]
			if !ok {
				continue
			}

			if tmpl != nil {
				ml, err := tmpl.New(name).Funcs(defaultFuncs).Parse(data)
				if err != nil {
					return nil, err
				}

				tmpl = ml
				continue
			}

			tml, err := template.New(name).Funcs(defaultFuncs).Parse(data)
			if err != nil {
				return nil, err
			}

			tmpl = tml
		}

		return tmpl, nil
	}

	for name, data := range g.templates {
		if tmpl != nil {
			ml, err := tmpl.New(name).Funcs(defaultFuncs).Parse(data)
			if err != nil {
				return nil, err
			}

			tmpl = ml
			continue
		}

		tml, err := template.New(name).Funcs(defaultFuncs).Parse(data)
		if err != nil {
			return nil, err
		}

		tmpl = tml
	}

	return tmpl, nil
}

func doIntToString(vals []interface{}) []string {
	var items []string
	for _, val := range vals {
		items = append(items, fmt.Sprintf("%+s", val))
	}
	return items
}

func cutListSuffix(sets []string, cutsuffix string) []string {
	var do []string

	for _, set := range sets {
		do = append(do, strings.TrimSuffix(set, cutsuffix))
	}

	return do
}

func cutListPrefix(sets []string, cutprefix string) []string {
	var do []string

	for _, set := range sets {
		do = append(do, strings.TrimPrefix(set, cutprefix))
	}

	return do
}

func doCutSplit(sets []string, sp string, index int) []string {
	var do []string

	for _, set := range sets {
		parts := strings.Split(set, sp)
		if index >= len(parts) {
			continue
		}

		do = append(do, parts[index])
	}

	return do
}

func cutList(sets []string, cut string) []string {
	var do []string

	for _, set := range sets {
		do = append(do, strings.Trim(set, cut))
	}

	return do
}

func doSuffix(elems []string, suffix string) []string {
	var do []string

	for _, item := range elems {
		do = append(do, fmt.Sprintf("%s%s", item, suffix))
	}

	return do
}

func doPrefix(elems []string, prefix string) []string {
	var do []string

	for _, item := range elems {
		do = append(do, fmt.Sprintf("%s%s", prefix, item))
	}

	return do
}

func doTimeSuffix(times int, suffix string) []string {
	var do []string

	for i := 0; i < times; i++ {
		do = append(do, fmt.Sprintf("%d%s", i, suffix))
	}

	return do
}

func doTimePrefix(times int, prefix string) []string {
	var do []string

	for i := 0; i < times; i++ {
		do = append(do, fmt.Sprintf("%s%d", prefix, i))
	}

	return do
}

func quote(b interface{}) string {
	switch bo := b.(type) {
	case string:
		return strconv.Quote(bo)
	case int:
		return strconv.Quote(strconv.Itoa(bo))
	case bool:
		return strconv.Quote(strconv.FormatBool(bo))
	case int64:
		return strconv.Quote(strconv.Itoa(int(bo)))
	case float32:
		mo := strconv.FormatFloat(float64(bo), 'f', 4, 32)
		return strconv.Quote(mo)
	case float64:
		mo := strconv.FormatFloat(bo, 'f', 4, 32)
		return strconv.Quote(mo)
	case byte:
		return strconv.QuoteRune(rune(bo))
	case rune:
		return strconv.QuoteRune(bo)
	default:
		return "Unconvertible Type"
	}
}

func nthOf(b interface{}, index int) (val interface{}) {
	switch bo := b.(type) {
	case []interface{}:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	case []string:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	case string:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	case []int:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	case []bool:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	case []int64:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	case []float32:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	case []float64:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	case []byte:
		if index >= len(bo) {
			return
		}
		val = bo[index]
	}

	return
}

func lenOff(b interface{}) int {
	switch bo := b.(type) {
	case []interface{}:
		return len(bo)
	case []string:
		return len(bo)
	case string:
		return len(bo)
	case []int:
		return len(bo)
	case []bool:
		return len(bo)
	case []int64:
		return len(bo)
	case []float32:
		return len(bo)
	case []float64:
		return len(bo)
	case []byte:
		return len(bo)
	default:
		return 0
	}
}
