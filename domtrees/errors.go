package domtrees

import "errors"

// Markup based errors relating to the type of markup

//ErrNotText is returned when the markup type is not a text markup
var ErrNotText = errors.New("Markup is not a *Text type")

// ErrNotElem is returned when the markup type does not match the *Element type
var ErrNotElem = errors.New("Markup is not a *Element type")

// ErrNotMarkup is returned when the value/pointer type does not match the Markup interface type
var ErrNotMarkup = errors.New("Value does not match Markup interface types")

// ErrNotAttr relating to the attribute types
var ErrNotAttr = errors.New("Value type is not n Attribute type")

// ErrNotFound relating to the attribute types
var ErrNotFound = errors.New("Item not found")

// Errors relating to the style types
var ErrNotStyle = errors.New("Value type is not a Style type")
