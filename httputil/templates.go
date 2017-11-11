package httputil

import (
	htemplate "html/template"
	"text/template"
)

// textContextFunctions returns a map of tempalte funcs for usage
// with text/template.Template.
func textContextFunctions(c *Context) template.FuncMap {
	return template.FuncMap{
		"flash":         c.Flash,
		"flashMessages": c.FlashMessages,
		"clearFlashMessages": func() string {
			c.ClearFlashMessages()
			return ""
		},
		"clearFlash": func(name string) string {
			c.ClearFlash(name)
			return ""
		},
		"setFlash": func(name, message string) string {
			c.SetFlash(name, message)
			return ""
		},
	}
}

// htmlContextFunctions returns a map of tempalte funcs for usage
// with text/template.Template.
func htmlContextFunctions(c *Context) htemplate.FuncMap {
	return htemplate.FuncMap{
		"flash":         c.Flash,
		"flashMessages": c.FlashMessages,
		"clearFlashMessages": func() string {
			c.ClearFlashMessages()
			return ""
		},
		"clearFlash": func(name string) string {
			c.ClearFlash(name)
			return ""
		},
		"setFlash": func(name, message string) string {
			c.SetFlash(name, message)
			return ""
		},
	}
}
