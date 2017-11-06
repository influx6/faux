package tmplutil

import (
	"errors"
	"fmt"
	"sync"
	"text/template"
)

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
				ml, err := tmpl.New(name).Parse(data)
				if err != nil {
					return nil, err
				}

				tmpl = ml
				continue
			}

			tml, err := template.New(name).Parse(data)
			if err != nil {
				return nil, err
			}

			tmpl = tml
		}

		return tmpl, nil
	}

	for name, data := range g.templates {
		if tmpl != nil {
			ml, err := tmpl.New(name).Parse(data)
			if err != nil {
				return nil, err
			}

			tmpl = ml
			continue
		}

		tml, err := template.New(name).Parse(data)
		if err != nil {
			return nil, err
		}

		tmpl = tml
	}

	return tmpl, nil
}
