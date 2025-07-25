package wm

import (
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
func Create() *RiwoElement {
	return CreateKnown("div")
}

// CreateKnown
// Initializes new container known tag.
// In example: CreateKnown("div") abcolutely equals Create()
func CreateKnown(tag string) *RiwoElement {
	return &RiwoElement{
		jsValue: js.Global().Get("document").Call("createElement", tag),
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

// Mount
// Appends current container to parent RiwoElement
func (e *RiwoElement) Mount(parent *RiwoElement) *RiwoElement {
	parent.Append(e)
	return e
}

// Get
// returns JavaScript data structure for current RiwoElement
func (e *RiwoElement) Get() js.Value {
	return e.jsValue
}
