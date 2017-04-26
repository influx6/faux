package dom

import (
	"regexp"
	"strings"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/faux/utils"
	"honnef.co/go/js/dom"
)

// Window returns the root window for the current dom.
func Window() dom.Window {
	return dom.GetWindow()
}

// Document returns the current document attached to the window
func Document() dom.Document {
	return Window().Document()
}

// RootElement returns the root parent element that for the provided element.
func RootElement(elem dom.Node) dom.Node {
	var parent dom.Node

	parent = elem

	for parent.ParentNode() != nil {
		parent = parent.ParentNode()
	}

	return parent
}

// HasShadowRoot returns true/false if the element type matches
// the [object shadowRoot].
func HasShadowRoot(elem dom.Node) bool {
	rs := RootElement(elem)
	return rs.Underlying().Call("toString").String() == "[object ShadowRoot]"
}

// ShadowRootDocument returns the DocumentFragment interface for the provided
// shadowRoot.
func ShadowRootDocument(elem dom.Node) dom.DocumentFragment {
	rs := RootElement(elem)
	return dom.WrapDocumentFragment(rs.Underlying())
}

// GetShadowRoot retrieves the shadowRoot connected to the pass dom.Node, else
// returns false as the second argument if the node has no shadowRoot.
func GetShadowRoot(elem dom.Node) (dom.DocumentFragment, bool) {
	if elem == nil {
		return nil, false
	}

	var root *js.Object

	if root = elem.Underlying().Get("shadowRoot"); root == nil {
		if root = elem.Underlying().Get("root"); root == nil {
			return nil, false
		}
	}

	return dom.WrapDocumentFragment(root), true
}

//==============================================================================

// topScrollAttr defines the apppropriate property to retrieve the top scroll
// value.
var topScrollAttr string

// leftScrollAttr defines the apppropriate property to retrieve the left scroll
// value.
var leftScrollAttr string

var useDocForOffset bool

// initScrollProperties initializes the necessary scroll property names.
func initScrollProperties() {
	if js.Global.Get("pageYOffset") != nil {
		topScrollAttr = "pageYOffset"
	} else {
		topScrollAttr = "scrollTop"
		useDocForOffset = true
	}

	if js.Global.Get("pageXOffset") != nil {
		leftScrollAttr = "pageXOffset"
	} else {
		leftScrollAttr = "scrollLeft"
		useDocForOffset = true
	}
}

//==============================================================================

// PageBox returns the offset of the current page bounding box.
func PageBox() (float64, float64) {
	var cursor *js.Object

	if useDocForOffset {
		cursor = Document().Underlying()
	} else {
		cursor = js.Global
	}

	top := cursor.Get(topScrollAttr)
	left := cursor.Get(leftScrollAttr)
	return utils.ParseFloat(top.String()), utils.ParseFloat(left.String())
}

// ClientBox returns the offset of the current page client box.
func ClientBox() (float64, float64) {
	top := Document().Underlying().Get("clientTop")
	left := Document().Underlying().Get("clientLeft")

	if top == nil || left == nil {
		return 0, 0
	}

	return utils.ParseFloat(top.String()), utils.ParseFloat(left.String())
}

// rootName defines a regexp for matching the string to either be body/html.
var rootName = regexp.MustCompile("^(?:body|html)$")

// Position returns the current position of the dom.Element.
func Position(elem dom.Element) (float64, float64) {
	parent := OffsetParent(elem)

	var parentTop, parentLeft float64
	var marginTop, marginLeft float64
	var pBorderTop, pBorderLeft float64
	var pBorderTopObject *js.Object
	var pBorderLeftObject *js.Object

	nodeNameObject, err := GetProp(parent, "nodeName")
	if err == nil && !rootName.MatchString(strings.ToLower(nodeNameObject.String())) {
		parentElem := dom.WrapElement(parent)
		parentTop, parentLeft = Offset(parentElem)
	}

	if parent.Get("style") != nil {

		pBorderTopObject, err = GetProp(parent, "style.borderTopWidth")
		if err == nil {
			pBorderTop = utils.ParseFloat(pBorderTopObject.String())
		}

		pBorderLeftObject, err = GetProp(parent, "style.borderLeftWidth")
		if err == nil {
			pBorderLeft = utils.ParseFloat(pBorderLeftObject.String())
		}

		parentTop += pBorderTop
		parentLeft += pBorderLeft
	}

	css, _ := GetComputedStyle(elem, "")

	marginTopObject, err := GetComputedStyleValueWith(css, "margin-top")
	if err == nil {
		marginTop = utils.ParseFloat(marginTopObject.String())
	}

	marginLeftObject, err := GetComputedStyleValueWith(css, "margin-left")
	if err == nil {
		marginLeft = utils.ParseFloat(marginLeftObject.String())
	}

	elemTop, elemLeft := Offset(elem)

	elemTop -= marginTop
	elemLeft -= marginLeft

	return elemTop - parentTop, elemLeft - parentLeft
}

// Offset returns the top,left offset of a dom.Element.
func Offset(elem dom.Element) (float64, float64) {
	boxTop, _, _, boxLeft := BoundingBox(elem)
	clientTop, clientLeft := ClientBox()
	pageTop, pageLeft := PageBox()

	top := (boxTop + pageTop) - clientTop
	left := (boxLeft + pageLeft) - clientLeft

	return top, left
}

// BoundingBox returns the top,right,down,left corners of a dom.Element.
func BoundingBox(elem dom.Element) (float64, float64, float64, float64) {
	rect := elem.GetBoundingClientRect()
	return rect.Top, rect.Right, rect.Bottom, rect.Left
}

//==============================================================================

// OffsetParent returns the offset parent element for a specific element.
func OffsetParent(elem dom.Element) *js.Object {
	und := elem.Underlying()

	osp, err := GetProp(und, "offsetParent")
	if err != nil {
		osp = Document().Underlying()
	}

	for osp != nil && !MatchProp(osp, "nodeType", "html") && MatchProp(osp, "style.position", "static") {
		val, err := GetProp(osp, "offsetParent")
		if err != nil {
			break
		}

		osp = val
	}

	return osp
}

//==============================================================================

// init initalizes properties and functions necessary for package wide varaibles.
func init() {
	if detect.IsBrowser() {
		initScrollProperties()
	}
}

//==============================================================================
