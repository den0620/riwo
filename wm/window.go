package wm

import (
	"strconv"
	"syscall/js"
)

// Type `Window` manages single abstract window’s properties.
type Window struct {
	ID      int
	Element js.Value // Connected DOM element.
	ContextEntries map[string]js.Value // Tho "Move", "Resize", "Delete" and "Hide" are basic ones
}


// NewWindow creates a new Window, sets up its DOM element, and returns a pointer to it.
func WindowCreate(x, y, width, height, content string) *Window {
	windowCount++
	id := windowCount

	document := js.Global().Get("document")
	body := document.Get("body")

	// Create the DOM element for the window.
	winElem := document.Call("createElement", "div")
	style := "overflow: hidden; position: absolute; z-index: 10; left: "+x+
		"; top: "+y+"; width: "+width+"; height: "+height+
		"; background-color: #f0f0f0; border: solid #55AAAA; padding: 0;"
	winElem.Set("style", style)
	winElem.Set("innerHTML", content)
	winElem.Set("id", strconv.Itoa(id))

	body.Call("appendChild", winElem)

	// Logging
	if verbose {
		Print("Generated window's title is \""
		+window.Get("title").String()+
		"\"; Window's ID (wid) is \""
		+window.Get("wid").String()+"\"")
	}

	return &Window{
		ID:      id,
		Element: winElem,
	}
}

// Move sets pos and dimensions for the window.
func (w *Window) WindowEditPos(newX, newY, newWidth, newHeight string) {
	w.Element.Get("style").Set("left", newX)
	w.Element.Get("style").Set("top", newY)
	w.Element.Get("style").Set("width", newWidth)
	w.Element.Get("style").Set("height", newHeight)
}

// Remove deletes the window from the DOM.
func (w *Window) WindowRemove() {
	w.Element.Call("remove")
}

