package wm

import (
	"fmt"
	"syscall/js"
)

// RiwoObject
// Main container type for WindowManager
type RiwoObject struct {
	jsValue js.Value
}

// Create
// Base function for <div> container initialization
// in Riwo environment
// (initializes public constructor for RiwoElement)
func Create() *RiwoObject {
	return CreateKnown("div")
}

// CreateKnown
// Initializes new container known tag.
// In example: CreateKnown("div") abcolutely equals Create()
// (initializes public constructor for RiwoElement)
func CreateKnown(tag string) *RiwoObject {
	return &RiwoObject{
		jsValue: js.Global().Get("document").Call("createElement", tag),
	}
}

// CreateFrom
// Makes new RiwoElement from by given DOM data
// (initializes private constructor for RiwoElement)
func CreateFrom(from *js.Value) *RiwoObject {
	return &RiwoObject{
		jsValue: *from,
	}
}

// ByID
func ByID(id string) *RiwoObject {
	return &RiwoObject{
		jsValue: js.Global().Get("document").Call("getElementById", id),
	}
}

func (e *RiwoObject) Id(id string) *RiwoObject {
	e.jsValue.Set("id", id)
	return e
}

// Class
// Gets and applies values from classes range
func (e *RiwoObject) Class(classes ...string) *RiwoObject {
	classList := e.jsValue.Get("classList")
	for _, cls := range classes {
		classList.Call("add", cls)
	}
	return e
}

// Style
// Applies target style ny name/value for current RiwoElement
func (e *RiwoObject) Style(style string, value interface{}) *RiwoObject {
	e.jsValue.Get("style").Set(style, value)
	return e
}

// Text
// Sets text content for current element
func (e *RiwoObject) Text(content string) *RiwoObject {
	e.jsValue.Set("textContent", content)
	return e
}

// Inner
// Sets Inner HTML value for current RiwoElement
func (e *RiwoObject) Inner(content string) *RiwoObject {
	e.jsValue.Set("innerHTML", content)
	return e
}

// Listen
// Applies event to current RiwoElement
func (e *RiwoObject) Listen(event string, handler func(this js.Value, args []js.Value) any) *RiwoObject {
	e.jsValue.Call("addEventListener", event, js.FuncOf(handler))
	return e
}

// Append
// Includes (appends) next expected elements range to current RiwoElement
func (e *RiwoObject) Append(children ...*RiwoObject) *RiwoObject {
	for _, child := range children {
		e.jsValue.Call("appendChild", child.jsValue)
	}
	return e
}
func (e *RiwoObject) AppendByDom(children ...js.Value) *RiwoObject {
	for _, child := range children {
		e.jsValue.Call("appendChild", child)
	}

	return e
}

// Attr
// Applies attribute to current element
func (e *RiwoObject) Attr(name, value string) *RiwoObject {
	e.jsValue.Call("setAttribute", name, value)
	return e
}

// Set
// Updates element's attrubute by name
func (e *RiwoObject) Set(name string, value interface{}) *RiwoObject {
	e.jsValue.Set(name, fmt.Sprintf("%s", value))
	return e
}

// Mount
// Appends current container to parent RiwoElement
func (e *RiwoObject) Mount(parent *RiwoObject) *RiwoObject {
	parent.Append(e)
	return e
}

// DOM
// returns JavaScript data structure for current RiwoElement
func (e *RiwoObject) DOM() js.Value {
	return e.jsValue
}

// From
// returns property value string
func (e *RiwoObject) From(property string) js.Value {
	return e.jsValue.Get(property)
}

// Call
// calls target key for current element
func (e *RiwoObject) Call(property string) js.Value {
	return e.DOM().Call(property)
}

// Delete
// Erases all data
func (e *RiwoObject) Delete() {
	e.jsValue = js.Null()
}
