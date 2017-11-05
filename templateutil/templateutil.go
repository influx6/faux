package templateutil

import (
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
func (g *Group) Add(name string, data []byte) {
	g.ml.Lock()
	defer g.ml.Unlock()
	g.templates[name] = string(data)
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
