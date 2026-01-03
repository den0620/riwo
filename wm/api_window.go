/*
Window struct for WM
and its functions + listeners
*/

package wm

import (
	"fmt"
	"strconv"
	"syscall/js"
)

// Define a named type for the context entry.
type ContextEntry struct {
	Name     string
	Callback func()
}

// Type `RiwoWindow` manages single abstract window’s properties.
type RiwoWindow struct {
	ID          int            // For the most part unites DOM object and Go object
	Title       string         // Yes, Riwo, Window title.
	Content     *RiwoObject    // Connected DOM element.
	MenuEntries []ContextEntry // Tho "Move", "Resize", "Delete" and "Hide" are basic ones
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

	windowContent := Create().
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
		Set("id", id).  // <-- assing shared ID
		Inner(content). // <-- spookie-dookie inner HTML
		Mount(bodyContent)

	// Logging
	JSLog("Generated window's ID (wid) is \"" +
		strconv.Itoa(id) + "\"")

	window := &RiwoWindow{
		ID:      id,
		Content: windowContent,
		Title:   fmt.Sprintf(" (wid=%d)", id),
		// No custom ContextEntries
	}

	CurrentWindow = window
	ActiveWindow = *windowContent

	AllWindows[strconv.Itoa(window.ID)] = window // <-- why string????? // i dont remember but probably because of js

	// Bring to front when clicked
	windowContent.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
		if !IsResizingInit {
			CurrentWindow = window
			ActiveWindow = *windowContent
		}

		// Right-click (RMB) on the window to select it for resizing, second right-click activates resizing
		if IsResizingMode && !IsResizingInit && args[0].Get("button").Int() == 2 {
			// First RMB hold - Select the window for resizing
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			JSLog("First right-click: Window selected for resizing.")

			windowContent.Style("z-index", strconv.Itoa(HighestZIndex))
			HighestZIndex++
			IsResizingInit = true

			bodyContent.Style("cursor", "url(assets/cursor-selection.svg) 12 12, auto")
		}

		// Mouse down event for selecting and dragging the window (click brings it to front)
		if !IsResizingInit {
			HighestZIndex++
			windowContent.Style("z-index", strconv.Itoa(HighestZIndex))
			JSLog("Window brought to front.")

			if IsMovingMode && args[0].Get("button").Int() == 2 {
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")
				//JustSelected = true
				StartX = args[0].Get("clientX").Float() - windowContent.From("offsetLeft").Float()
				StartY = args[0].Get("clientY").Float() - windowContent.From("offsetTop").Float()
				IsDragging = true

				// Create ghost window
				rect := windowContent.Call("getBoundingClientRect")
				width := rect.Get("width").Float()
				height := rect.Get("height").Float()

				// Ensure ghost window is above everything during drag
				GhostWindow = *Create().
					Style("left", Ftoa(windowContent.From("offsetLeft").Float())+"px").
					Style("top", Ftoa(windowContent.From("offsetTop").Float())+"px").
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
				hiddenWindowOption := CreateMenuObject(fmt.Sprintf("%s (#%d)", window.Title, window.ID))

				if windowContent.From("title").String() != "" {
					hiddenWindowOption = CreateMenuObject(windowContent.From("title").String())
				}

				hiddenWindowOption.DOM().Set("id", "menuopt"+strconv.Itoa(window.ID))

				// ??? option activation
				hiddenWindowOption.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
					if args[0].Get("button").Int() == 2 {
						args[0].Call("preventDefault")
						args[0].Call("stopPropagation")
						JustSelected = true

						RemoveMenuOption(hiddenWindowOption.DOM())
						windowContent.Style("display", "block")

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

				windowContent.Style("display", "none")
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
// Deletes the window from DOM and Go.
func removeWindow(w *RiwoWindow) {
	w.ID = -1                              // Remove reference for apps
	w.Content.Call("remove")               // Remove html part
	delete(AllWindows, strconv.Itoa(w.ID)) // Remove from list
}
