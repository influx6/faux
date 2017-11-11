package httputil

import (
	htemplate "html/template"
	"text/template"
)

// textContextFunctions returns a map of tempalte funcs for usage
// with text/template.Template.
func textContextFunctions(c *Context) template.FuncMap {
	return template.FuncMap{
		"flash":              c.Flash,
		"setFlash":           c.SetFlash,
		"flashMessages":      c.FlashMessages,
		"clearFlash":         c.ClearFlash,
		"clearFlashMessages": c.ClearFlashMessages,
	}
}

// htmlContextFunctions returns a map of tempalte funcs for usage
// with text/template.Template.
func htmlContextFunctions(c *Context) htemplate.FuncMap {
	return htemplate.FuncMap{
		"flash":              c.Flash,
		"setFlash":           c.SetFlash,
		"flashMessages":      c.FlashMessages,
		"clearFlash":         c.ClearFlash,
		"clearFlashMessages": c.ClearFlashMessages,
	}
}
