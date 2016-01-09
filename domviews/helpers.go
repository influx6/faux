package domviews

import (
	"bytes"
	"html/template"
)

// TemplateRenderable defines a basic example of a Renderable
type TemplateRenderable struct {
	tmpl  *template.Template
	cache *bytes.Buffer
}

// NewTemplateRenderable returns a new Renderable
func NewTemplateRenderable(content string) (*TemplateRenderable, error) {
	tl, err := template.New("").Parse(content)
	if err != nil {
		return nil, err
	}

	tr := TemplateRenderable{
		tmpl:  tl,
		cache: bytes.NewBuffer([]byte{}),
	}

	return &tr, nil
}

// Execute effects the inner template with the supplied data
func (t *TemplateRenderable) Execute(v interface{}) error {
	t.cache.Reset()
	err := t.tmpl.Execute(t.cache, v)
	return err
}

// Render renders out the internal cache
func (t *TemplateRenderable) Render(_ ...string) string {
	return string(t.cache.Bytes())
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (t *TemplateRenderable) RenderHTML(_ ...string) template.HTML {
	return template.HTML(t.Render())
}

// String calls the render function
func (t *TemplateRenderable) String() string {
	return t.Render()
}
