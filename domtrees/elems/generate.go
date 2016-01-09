// +build ignore
// Credit to Richard Musiol (https://github.com/neelance/dom)
// His code was crafted to fit haiku's use

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/influx6/flux"
)

var elemNameMap = map[string]string{
	"a":          "Anchor",
	"article":    "Article",
	"aside":      "Aside",
	"area":       "Area",
	"abbr":       "Abbreviation",
	"b":          "Bold",
	"base":       "Base",
	"bdi":        "BidirectionalIsolation",
	"bdo":        "BidirectionalOverride",
	"blockquote": "BlockQuote",
	"br":         "Break",
	"cite":       "Citation",
	"col":        "Column",
	"colgroup":   "ColumnGroup",
	"datalist":   "DataList",
	"dialog":     "Dialog",
	"details":    "Details",
	"dd":         "Description",
	"del":        "DeletedText",
	"dfn":        "Definition",
	"dl":         "DescriptionList",
	"dt":         "DefinitionTerm",
	"em":         "Emphasis",
	"embed":      "Embed",
	"footer":     "Footer",
	"figure":     "Figure",
	"figcaption": "FigureCaption",
	"fieldset":   "FieldSet",
	"h1":         "Header1",
	"h2":         "Header2",
	"h3":         "Header3",
	"h4":         "Header4",
	"h5":         "Header5",
	"h6":         "Header6",
	"hgroup":     "HeadingsGroup",
	"header":     "Header",
	"hr":         "HorizontalRule",
	"i":          "Italic",
	"iframe":     "InlineFrame",
	"img":        "Image",
	"ins":        "InsertedText",
	"kbd":        "KeyboardInput",
	"keygen":     "KeyGen",
	"li":         "ListItem",
	"meta":       "Meta",
	"menuitem":   "MenuItem",
	"nav":        "Navigation",
	"noframes":   "NoFrames",
	"noscript":   "NoScript",
	"ol":         "OrderedList",
	"option":     "Option",
	"optgroup":   "OptionsGroup",
	"p":          "Paragraph",
	"param":      "Parameter",
	"pre":        "Preformatted",
	"q":          "Quote",
	"rp":         "RubyParenthesis",
	"rt":         "RubyText",
	"s":          "Strikethrough",
	"samp":       "Sample",
	"source":     "Source",
	"section":    "Section",
	"sub":        "Subscript",
	"sup":        "Superscript",
	"tbody":      "TableBody",
	"textarea":   "TextArea",
	"td":         "TableData",
	"tfoot":      "TableFoot",
	"th":         "TableHeader",
	"thead":      "TableHead",
	"tr":         "TableRow",
	"u":          "Underline",
	"ul":         "UnorderedList",
	"var":        "Variable",
	"track":      "Track",
	"wbr":        "WordBreakOpportunity",
}

//list of self closing tags
var autoclosers = map[string]bool{
	"area":    true,
	"base":    true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"input":   true,
	"keygen":  true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
	"br":      true,
}

func main() {
	doc, err := goquery.NewDocument("https://developer.mozilla.org/en-US/docs/Web/HTML/Element")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("elems.gen.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, `// Package elems contains definition for different html element types and some custom node types

//go:generate go run generate.go

// Documentation source: "HTML element reference" by Mozilla Contributors, https://developer.mozilla.org/en-US/docs/Web/HTML/Element, licensed under CC-BY-SA 2.5.

// Documentation for custom types lies within the  "github.com/influx6/faux/domtrees" package

package elems

import (
	"github.com/influx6/faux/domtrees"
)

// Text provides the concrete implementation for using the domtrees.Text struct
func Text(txt string) *domtrees.Element {
	return domtrees.NewText(txt)
}
`)

	doc.Find(".quick-links a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if !strings.HasPrefix(link, "/en-US/docs/Web/HTML/Element/") {
			return
		}

		if s.Parent().Find(".icon-trash, .icon-thumbs-down-alt, .icon-warning-sign").Length() > 0 {
			return
		}

		desc, _ := s.Attr("title")

		text := s.Text()
		if text == "Heading elements" {
			writeElem(file, "h1", desc, link)
			writeElem(file, "h2", desc, link)
			writeElem(file, "h3", desc, link)
			writeElem(file, "h4", desc, link)
			writeElem(file, "h5", desc, link)
			writeElem(file, "h6", desc, link)
			return
		}

		name := text[1 : len(text)-1]
		if name == "html" || name == "head" || name == "body" {
			return
		}

		writeElem(file, name, desc, link)
	})
}

func writeElem(w io.Writer, name, desc, link string) {
	var autocloser = autoclosers[name]
	funName := elemNameMap[name]

	if funName == "" {
		funName = flux.Capitalize(name)
	}

	fmt.Fprintf(w, `
// %s provides the following for html elements ->
// %s
// https://developer.mozilla.org%s
func %s(markup ...domtrees.Appliable) *domtrees.Element {
	e := domtrees.NewElement("%s",%t)
	for _, m := range markup {
		m.Apply(e)
	}
	return e
}
`, funName, desc, link, funName, name, autocloser)
}
