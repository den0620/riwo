/*
Window struct for WM
and its functions + listeners
*/

package wm

import (
	//"fmt"
	"strconv"
	"syscall/js"
)

// Define a named type for the context entry.
type ContextEntry struct {
	Name     string
	Callback func()
}

// Type `RiwoWindow` manages single abstract window's properties.
type RiwoWindow struct {
	ID             int             // DOM id shared with Frame
	Title          string          // Riwo window title string
	Frame          *RiwoObject     // Outer draggable/resizable frame (fills placement)
	Content        *RiwoObject     // Guest / launcher surface mounted inside Frame
	MenuEntries    []ContextEntry // WM-local menu rows plus guest rows bridged separately
}

// windowPlacement is a data container, visible only in `wm` package
// Functions which manage windows, don't need information about window location
// Once call of this structure is a call from big arguments of `CreateWindow` function.
// It would be better if Riwo will try to hide `wm` implementation details from other packages.
type windowPlacement struct {
	x      string
	y      string
	width  string
	height string
}

// createWindow
// Creates a new Window, sets up its DOM element, and returns a pointer.
func createWindow(p *windowPlacement, content string) *RiwoWindow {
	WindowCount++
	id := WindowCount

	document := js.Global().Get("document")
	body := document.Get("body")

	bodyContent := CreateFrom(&body)

	frame := Create().
		Style("overflow", "hidden").
		Style("position", "absolute").
		Style("width", p.width).
		Style("height", p.height).
		Style("top", p.y).
		Style("left", p.x).
		Style("z-index", strconv.Itoa(HighestZIndex)).
		Style("background-color", "#f0f0f0").
		Style("border", "solid #55AAAA").
		Style("padding", "0").
		Set("id", id).
		Mount(bodyContent)

	pane := Create().
		Style("height", "100%").
		Style("width", "100%").
		Style("overflow", "hidden").
		Style("background-color", "#f0f0f0").
		Style("position", "relative").
		Style("boxSizing", "border-box")

	frame.Append(pane)
	if content != "" {
		pane.Inner(content)
	}

	// Logging
	JSLog("Generated window's ID (wid) is \"" +
		strconv.Itoa(id) + "\"")

	window := &RiwoWindow{
		ID:      id,
		Frame:   frame,
		Content: pane,
		Title:   " (wid=" + strconv.Itoa(id) + ")",
	}

	CurrentWindow = window
	ActiveWindow = *frame

	AllWindows[strconv.Itoa(window.ID)] = window // <-- why string????? // i dont remember but probably because of js

	// Bring to front when clicked
	frame.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
		if !IsResizingInit {
			CurrentWindow = window
			ActiveWindow = *frame
		}

		// Right-click (RMB) on the window to select it for resizing, second right-click activates resizing
		if IsResizingMode && !IsResizingInit && args[0].Get("button").Int() == 2 {
			// First RMB hold - Select the window for resizing
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			JSLog("First right-click: Window selected for resizing.")

			frame.Style("z-index", strconv.Itoa(HighestZIndex))
			HighestZIndex++
			IsResizingInit = true

			bodyContent.Style("cursor", "url(assets/cursor-selection.svg) 12 12, auto")
		}

		// Mouse down event for selecting and dragging the window (click brings it to front)
		if !IsResizingInit {
			HighestZIndex++
			frame.Style("z-index", strconv.Itoa(HighestZIndex))
			JSLog("Window brought to front.")

			if IsMovingMode && args[0].Get("button").Int() == 2 {
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")
				//JustSelected = true
				StartX = args[0].Get("clientX").Float() - frame.From("offsetLeft").Float()
				StartY = args[0].Get("clientY").Float() - frame.From("offsetTop").Float()
				IsDragging = true

				// Create ghost window
				rect := frame.Call("getBoundingClientRect")
				width := rect.Get("width").Float()
				height := rect.Get("height").Float()

				// Ensure ghost window is above everything during drag
				GhostWindow = *Create().
					Style("left", Ftoa(frame.From("offsetLeft").Float())+"px").
					Style("top", Ftoa(frame.From("offsetTop").Float())+"px").
					Style("position", "absolute").Style("z-index", strconv.Itoa(HighestZIndex+1)).
					Style("width", Ftoa(width)+"px").
					Style("height", Ftoa(height)+"px").
					Style("border", "solid 2px #FF0000").
					Style("cursor", "url(assets/cursor-drag.svg) 12 12, auto").
					Mount(bodyContent) // |<-- Append it to bodyContent

				JustSelected = true
				JSLog("Dragging initiated with ghost window.")
			}

			if IsHiding && args[0].Get("button").Int() == 2 {
				// Hide window
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")
				JustSelected = true
				IsHiding = false

				// prepare menu item
				hiddenWindowOption := CreateMenuObject(window.Title + " (#" + strconv.Itoa(window.ID) + ")")

				if frame.From("title").String() != "" {
					hiddenWindowOption = CreateMenuObject(frame.From("title").String())
				}

				hiddenWindowOption.DOM().Set("id", "menuopt"+strconv.Itoa(window.ID))

				// ??? option activation
				hiddenWindowOption.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
					if args[0].Get("button").Int() == 2 {
						args[0].Call("preventDefault")
						args[0].Call("stopPropagation")
						JustSelected = true

						RemoveMenuOption(hiddenWindowOption.DOM())
						frame.Style("display", "block")

						ContextMenu.Style("display", "none")

						// Delete by value
						for index, value := range ContextMenuHides {
							if value.Get("id").String() == hiddenWindowOption.DOM().Get("id").String() {
								ContextMenuHides = append(ContextMenuHides[:index], ContextMenuHides[index+1:]...)
							}
						}
						JustSelected = false
						JSLog("Unhide activated.")
					}
					return nil
				})
				ContextMenuHides = append(ContextMenuHides, hiddenWindowOption.DOM())

				frame.Style("display", "none")
				bodyContent.Style("cursor", "url(assets/cursor.svg), auto")

				JustSelected = false
				JSLog("WID " + strconv.Itoa(window.ID) + " hidden")

			}
		}
		// Right-click (RMB) deletes the window in delete mode
		if IsDeleteMode && args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			removeWindow(window)
			IsDeleteMode = false

			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")

			JustSelected = false
			JSLog("Window deleted.")
		}
		return nil
	})

	return window
}

// removeWindow
func removeWindow(w *RiwoWindow) {
	oldStr := strconv.Itoa(w.ID)
	js.Global().Get("__riwoKernel").Call("disposeGuestForWindow", w.ID)

	if w.Frame != nil && w.Frame.DOM().Truthy() {
		w.Frame.Call("remove")
	} else if w.Content != nil && w.Content.DOM().Truthy() {
		w.Content.Call("remove")
	}
	delete(AllWindows, oldStr)
	w.ID = -1
}
