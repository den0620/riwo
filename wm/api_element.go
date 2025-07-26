package wm

import (
	"fmt"
	"syscall/js"
)

// RiwoElement
// Main container type for WindowManager
type RiwoElement struct {
	jsValue js.Value
}

// Create
// Base function for <div> container initialization
// in Riwo environment
// (initializes public constructor for RiwoElement)
func Create() *RiwoElement {
	return CreateKnown("div")
}

// CreateKnown
// Initializes new container known tag.
// In example: CreateKnown("div") abcolutely equals Create()
// (initializes public constructor for RiwoElement)
func CreateKnown(tag string) *RiwoElement {
	return &RiwoElement{
		jsValue: js.Global().Get("document").Call("createElement", tag),
	}
}

// CreateFrom
// Makes new RiwoElement from by given DOM data
// (initializes private constructor for RiwoElement)
func CreateFrom(from *js.Value) *RiwoElement {
	return &RiwoElement{
		jsValue: *from,
	}
}

// ByID
func ByID(id string) *RiwoElement {
	return &RiwoElement{
		jsValue: js.Global().Get("document").Call("getElementById", id),
	}
}

func (e *RiwoElement) Id(id string) *RiwoElement {
	e.jsValue.Set("id", id)
	return e
}

// Class
// Gets and applies values from classes range
func (e *RiwoElement) Class(classes ...string) *RiwoElement {
	classList := e.jsValue.Get("classList")
	for _, cls := range classes {
		classList.Call("add", cls)
	}
	return e
}

// Style
// Applies target style ny name/value for current RiwoElement
func (e *RiwoElement) Style(style string, value interface{}) *RiwoElement {
	e.jsValue.Get("style").Set(style, value)
	return e
}

// Text
// Sets text content for current element
func (e *RiwoElement) Text(content string) *RiwoElement {
	e.jsValue.Set("textContent", content)
	return e
}

// Inner
// Sets Inner HTML value for current RiwoElement
func (e *RiwoElement) Inner(content string) *RiwoElement {
	e.jsValue.Set("innerHTML", content)
	return e
}

// Callback
// Applies event to current RiwoElement
func (e *RiwoElement) Callback(event string, handler func(this js.Value, args []js.Value) any) *RiwoElement {
	e.jsValue.Call("addEventListener", event, js.FuncOf(handler))
	return e
}

// Append
// Includes (appends) next expected elements range to current RiwoElement
func (e *RiwoElement) Append(children ...*RiwoElement) *RiwoElement {
	for _, child := range children {
		e.jsValue.Call("appendChild", child.jsValue)
	}
	return e
}

// Attr
// Applies attribute to current element
func (e *RiwoElement) Attr(name, value string) *RiwoElement {
	e.jsValue.Call("setAttribute", name, value)
	return e
}

// Set
// Updates element's attrubute by name
func (e *RiwoElement) Set(name string, value interface{}) *RiwoElement {
	e.jsValue.Set(name, fmt.Sprintf("%s", value))
	return e
}

// Mount
// Appends current container to parent RiwoElement
func (e *RiwoElement) Mount(parent *RiwoElement) *RiwoElement {
	parent.Append(e)
	return e
}

// DOM
// returns JavaScript data structure for current RiwoElement
func (e *RiwoElement) DOM() js.Value {
	return e.jsValue
}

// From
// returns property value string
func (e *RiwoElement) From(property string) js.Value {
	return e.jsValue.Get(property)
}

// Call
// calls target key for current element
func (e *RiwoElement) Call(property string) js.Value {
	return e.DOM().Call(property)
}

// Delete
// Erases all data
func (e *RiwoElement) Delete() {
	e.jsValue = js.Null()
}
