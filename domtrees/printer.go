package domtrees

import (
	"fmt"
	"strings"

	"github.com/go-humble/detect"
)

// This contains printers for the tree dom definition structures

// AttrPrinter defines a printer interface for writing out a Attribute objects into a string form
type AttrPrinter interface {
	Print([]*Attribute) string
}

// AttrWriter provides a concrete struct that meets the AttrPrinter interface
type AttrWriter struct{}

// SimpleAttrWriter provides a basic attribute writer
var SimpleAttrWriter = &AttrWriter{}

const attrformt = ` %s="%s"`

// Print returns a stringed repesentation of the attribute object
func (m *AttrWriter) Print(a []*Attribute) string {
	if len(a) <= 0 {
		return ""
	}

	attrs := []string{}

	for _, ar := range a {
		attrs = append(attrs, fmt.Sprintf(attrformt, ar.Name, ar.Value))
	}

	return strings.Join(attrs, " ")
}

// StylePrinter defines a printer interface for writing out a style objects into a string form
type StylePrinter interface {
	Print([]*Style) string
}

// StyleWriter provides a concrete struct that meets the AttrPrinter interface
type StyleWriter struct{}

// SimpleStyleWriter provides a basic style writer
var SimpleStyleWriter = &StyleWriter{}

const styleformt = " %s:%s;"

// Print returns a stringed repesentation of the style object
func (m *StyleWriter) Print(s []*Style) string {
	if len(s) <= 0 {
		return ""
	}

	css := []string{}

	for _, cs := range s {
		css = append(css, fmt.Sprintf(styleformt, cs.Name, cs.Value))
	}

	return strings.Join(css, " ")
}

// TextPrinter defines a printer interface for writing out a text type markup into a string form
type TextPrinter interface {
	Print(Markup) string
}

// TextWriter writes out the text element/node for the vdom into a string
type TextWriter struct{}

// SimpleTextWriter provides a basic text writer
var SimpleTextWriter = &TextWriter{}

// Print returns the string representation of the text object
func (m *TextWriter) Print(t Markup) string {
	return t.TextContent()
}

// ElementWriter writes out the element out as a string matching the html tag rules
type ElementWriter struct {
	attrWriter   AttrPrinter
	styleWriter  StylePrinter
	text         TextPrinter
	allowRemoved bool
}

// SimpleElementWriter provides a default writer using the basic attribute and style writers
var SimpleElementWriter = NewElementWriter(SimpleAttrWriter, SimpleStyleWriter, SimpleTextWriter)

// NewElementWriter returns a new writer for Element objects
func NewElementWriter(aw AttrPrinter, sw StylePrinter, tw TextPrinter) *ElementWriter {
	return &ElementWriter{
		attrWriter:  aw,
		styleWriter: sw,
		text:        tw,
	}
}

/*<<<---------------code within this region is usually for testing purposes------------*/

// DisallowRemoved is used to switch off the check to allow rendering of elements set as removed
func (m *ElementWriter) DisallowRemoved() {
	m.allowRemoved = false
}

// AllowRemoved is used to switch off the check to allow rendering of elements set as removed
func (m *ElementWriter) AllowRemoved() {
	m.allowRemoved = true
}

/* ----------------code within this region is usually for testing purposes----------->>>*/

// Print returns the string representation of the element
func (m *ElementWriter) Print(e *Element) string {
	// if we are on the server && is this element marked as removed, if so we skip and return an empty string
	if detect.IsServer() {
		if e.Removed() && !m.allowRemoved {
			return ""
		}
	}

	//if we are dealing with a text type just return the content
	if e.Name() == "text" {
		return m.text.Print(e)
	}

	//collect uid and hash of the element so we can write them along
	hash := &Attribute{"hash", e.Hash()}
	uid := &Attribute{"uid", e.UID()}

	//management attributes
	mido := []*Attribute{hash, uid}

	// if e.Removed() {
	// 	mido = append(mido, &Attribute{"haikuRemoved", ""})
	// }

	//write out the hash and uid as attributes
	hashes := m.attrWriter.Print(mido)

	//write out the elements attributes using the AttrWriter
	attrs := m.attrWriter.Print(e.Attributes())

	//write out the elements inline-styles using the StyleWriter
	style := m.styleWriter.Print(e.Styles())

	var closer string
	var beginbrack string

	if e.AutoClosed() {
		closer = "/>"
	} else {
		beginbrack = ">"
		closer = fmt.Sprintf("</%s>", e.Name())
	}

	var children = []string{}

	for _, ch := range e.Children() {
		// if ch.Name() == "text" {
		// 	children = append(children, m.text.Print(ch))
		// 	continue
		// }
		if ech, ok := ch.(*Element); ok {
			if ech == e {
				continue
			}
			children = append(children, m.Print(ech))
		}
	}

	//lets create the elements markup now
	return strings.Join([]string{
		fmt.Sprintf("<%s", e.Name()),
		hashes,
		attrs,
		fmt.Sprintf(` style="%s"`, style),
		beginbrack,
		e.textContent,
		strings.Join(children, ""),
		closer,
	}, "")
}

// MarkupWriter defines a printer interface for writing out a markup object into a string form
type MarkupWriter interface {
	Write(Markup) (string, error)
}

// MarkupWriter provides the concrete struct that meets the MarkupPrinter interface
type markupWriter struct {
	*ElementWriter
}

// SimpleMarkupWriter provides a basic markup writer for handling the different markup elements
var SimpleMarkupWriter = NewMarkupWriter(SimpleElementWriter)

// NewMarkupWriter returns a new markup instance
func NewMarkupWriter(em *ElementWriter) MarkupWriter {
	return &markupWriter{em}
}

// Write returns a stringed repesentation of the markup object
func (m *markupWriter) Write(ma Markup) (string, error) {
	// if tmr, ok := ma.(*Text); ok {
	// 	return m.ElementWriter.text.Print(tmr), nil
	// }

	if emr, ok := ma.(*Element); ok {
		return m.ElementWriter.Print(emr), nil
	}

	return "", ErrNotMarkup
}
